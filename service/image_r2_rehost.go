package service

// 通用 R2 图片转存执行（上传/下载）；转存策略见 relay/imagevendor。

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay/imagevendor"
)

// RewriteLoopbackUpstreamImageURL 将上游 loopback 图片地址（如 Gulie 127.0.0.1:3001）
// 映射为渠道主机名 + 原端口，便于服务端下载。
func RewriteLoopbackUpstreamImageURL(channelBaseURL, imageURL string) string {
	channelBaseURL = strings.TrimSpace(channelBaseURL)
	if channelBaseURL == "" {
		return imageURL
	}
	img, err := url.Parse(imageURL)
	if err != nil {
		return imageURL
	}
	host := strings.ToLower(img.Hostname())
	if host != "127.0.0.1" && host != "localhost" {
		return imageURL
	}
	base, err := url.Parse(channelBaseURL)
	if err != nil || base.Hostname() == "" {
		return imageURL
	}
	port := img.Port()
	if port == "" {
		if img.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	out := &url.URL{
		Scheme:   base.Scheme,
		Host:     net.JoinHostPort(base.Hostname(), port),
		Path:     img.Path,
		RawQuery: img.RawQuery,
	}
	return out.String()
}

func imageDataNeedsURLRehost(images []dto.ImageData) bool {
	for _, item := range images {
		if strings.TrimSpace(item.Url) != "" && strings.TrimSpace(item.B64Json) == "" {
			return true
		}
	}
	return false
}

func imageDataHasB64(images []dto.ImageData) bool {
	for _, item := range images {
		if strings.TrimSpace(item.B64Json) != "" {
			return true
		}
	}
	return false
}

func decodeImageB64Payload(b64 string) ([]byte, string, error) {
	data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
	if err != nil {
		return nil, "", err
	}
	mimeType := http.DetectContentType(data)
	if !strings.HasPrefix(mimeType, "image/") {
		mimeType = "image/png"
	}
	return data, mimeType, nil
}

func syncImageNeedsRehost(originModel string, images []dto.ImageData, clientWantsURL bool) bool {
	policy := imagevendor.ResolveRehostPolicy(originModel)
	if !policy.AcceptUpstreamURL || len(images) == 0 {
		return false
	}
	if imageDataNeedsURLRehost(images) {
		return true
	}
	return clientWantsURL && imageDataHasB64(images)
}

// RehostImageDataForClient 将上游 b64 或（4K/FLUX）url 转存 R2；clientWantsURL 时清除 b64_json 仅留公网 url。
func RehostImageDataForClient(ctx context.Context, userID int, storeID, channelBaseURL, originModel string, images []dto.ImageData, clientWantsURL bool) ([]dto.ImageData, error) {
	if !syncImageNeedsRehost(originModel, images, clientWantsURL) {
		return images, nil
	}
	if getR2Config() == nil {
		return nil, fmt.Errorf("R2 not configured")
	}
	storeID = strings.TrimSpace(storeID)
	if storeID == "" {
		storeID = model.GenerateTaskID()
	}
	out := make([]dto.ImageData, len(images))
	copy(out, images)
	for index := range out {
		item := &out[index]
		if b64 := strings.TrimSpace(item.B64Json); b64 != "" {
			data, mimeType, err := decodeImageB64Payload(b64)
			if err != nil {
				return nil, err
			}
			uploaded, err := UploadGeneratedImageBytes(ctx, userID, storeID, index, data, mimeType)
			if err != nil {
				return nil, fmt.Errorf("rehost upstream image b64: %w", err)
			}
			item.Url = uploaded.PublicURL
			if clientWantsURL {
				item.B64Json = ""
			}
			continue
		}
		if strings.TrimSpace(item.Url) == "" {
			continue
		}
		if !imagevendor.ResolveRehostPolicy(originModel).AcceptUpstreamURL {
			continue
		}
		downloadURL := RewriteLoopbackUpstreamImageURL(channelBaseURL, item.Url)
		uploaded, err := UploadGeneratedImageFromURL(ctx, userID, storeID, index, downloadURL)
		if err != nil {
			return nil, fmt.Errorf("rehost upstream image url: %w", err)
		}
		item.Url = uploaded.PublicURL
	}
	return out, nil
}

// RehostImageDataURLs 将需转存的模型上游 url 落 R2；未命中或无 url 时原样返回。
func RehostImageDataURLs(ctx context.Context, userID int, storeID, channelBaseURL, originModel string, images []dto.ImageData) ([]dto.ImageData, error) {
	return RehostImageDataForClient(ctx, userID, storeID, channelBaseURL, originModel, images, false)
}

// RehostSyncImageResponseBody 同步生图 JSON 响应：b64 或 url 转存 R2 后返回公网 URL。
func RehostSyncImageResponseBody(ctx context.Context, userID int, originModel, channelBaseURL string, responseBody []byte, clientWantsURL bool) ([]byte, error) {
	if len(responseBody) == 0 {
		return responseBody, nil
	}
	var payload struct {
		Data []dto.ImageData `json:"data"`
	}
	if err := common.Unmarshal(responseBody, &payload); err != nil || len(payload.Data) == 0 {
		return responseBody, nil
	}
	if !syncImageNeedsRehost(originModel, payload.Data, clientWantsURL) {
		return responseBody, nil
	}
	rehosted, err := RehostImageDataForClient(ctx, userID, model.GenerateTaskID(), channelBaseURL, originModel, payload.Data, clientWantsURL)
	if err != nil {
		return nil, err
	}
	var raw map[string]json.RawMessage
	if err := common.Unmarshal(responseBody, &raw); err != nil {
		return responseBody, nil
	}
	dataJSON, err := common.Marshal(rehosted)
	if err != nil {
		return nil, err
	}
	raw["data"] = dataJSON
	return common.Marshal(raw)
}

// DecodeImageDataItem 解析 ImageData：b64_json、data URI（url 字段）或上游 http(s) url。
// 有字节数据时返回 (data, mimeType, nil)；仅有 URL 时返回 (nil, url, nil)。
func DecodeImageDataItem(item dto.ImageData) ([]byte, string, error) {
	if item.B64Json != "" {
		data, err := base64.StdEncoding.DecodeString(item.B64Json)
		if err != nil {
			return nil, "", err
		}
		return data, detectImageBytesMimeType(data), nil
	}
	if item.Url == "" {
		return nil, "", fmt.Errorf("image item has no url or b64_json")
	}
	if strings.HasPrefix(item.Url, "data:") {
		data, mime, err := decodeDataURI(item.Url)
		if err != nil {
			return nil, "", err
		}
		return data, mime, nil
	}
	return nil, item.Url, nil
}

func decodeDataURI(uri string) ([]byte, string, error) {
	comma := strings.Index(uri, ",")
	if comma < 0 {
		return nil, "", fmt.Errorf("invalid data uri")
	}
	meta := uri[5:comma]
	payload := uri[comma+1:]
	mimeType := "image/png"
	if semi := strings.Index(meta, ";"); semi > 0 {
		mimeType = meta[:semi]
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", err
	}
	if !strings.HasPrefix(mimeType, "image/") {
		mimeType = detectImageBytesMimeType(data)
	}
	return data, mimeType, nil
}

func detectImageBytesMimeType(data []byte) string {
	if len(data) == 0 {
		return "image/png"
	}
	mimeType := http.DetectContentType(data)
	if strings.HasPrefix(mimeType, "image/") {
		return mimeType
	}
	return "image/png"
}

// RehostTaskImageResultURLs 异步 task 专用：转存后返回公网 URL 列表。
// 支持 b64_json、data URI（url 字段）、上游 http(s) url（受 policy 约束）。
func RehostTaskImageResultURLs(ctx context.Context, userID int, storeID, channelBaseURL, originModel string, images []dto.ImageData) ([]string, error) {
	acceptUpstreamURL := imagevendor.ImageAsyncAcceptsUpstreamURL(originModel)
	resultURLs := make([]string, 0, len(images))
	storeID = strings.TrimSpace(storeID)
	for index, item := range images {
		data, mimeOrURL, err := DecodeImageDataItem(item)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			if getR2Config() == nil {
				return nil, fmt.Errorf("R2 not configured")
			}
			mimeType := mimeOrURL
			if !strings.HasPrefix(mimeType, "image/") {
				mimeType = "image/png"
			}
			uploaded, err := UploadGeneratedImageBytes(ctx, userID, storeID, index, data, mimeType)
			if err != nil {
				return nil, err
			}
			resultURLs = append(resultURLs, uploaded.PublicURL)
			continue
		}
		if mimeOrURL != "" {
			if acceptUpstreamURL {
				if getR2Config() == nil {
					return nil, fmt.Errorf("R2 not configured")
				}
				downloadURL := RewriteLoopbackUpstreamImageURL(channelBaseURL, mimeOrURL)
				uploaded, err := UploadGeneratedImageFromURL(ctx, userID, storeID, index, downloadURL)
				if err != nil {
					return nil, fmt.Errorf("rehost upstream image url: %w", err)
				}
				resultURLs = append(resultURLs, uploaded.PublicURL)
				continue
			}
			return nil, fmt.Errorf("upstream returned url without b64_json; use response_format=b64_json")
		}
	}
	if len(resultURLs) == 0 {
		return nil, fmt.Errorf("no image results from upstream")
	}
	return resultURLs, nil
}
