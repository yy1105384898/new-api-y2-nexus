package openai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

var adobe2APIImageModelPrefixes = []string{
	"nano-banana",
	"gpt-image",
	"adobe-nano-banana",
	"adobe-gpt-image",
	"adobe2api-nano-banana",
	"adobe2api-gpt-image",
	"firefly-nano-banana",
	"firefly-gpt-image",
}

var adobe2APIVideoModelPrefixes = []string{
	"sora2",
	"veo31",
	"kling3",
	"kling-o3",
	"seedance2",
	"adobe-sora2",
	"adobe-veo31",
	"adobe-kling3",
	"adobe-kling-o3",
	"adobe-seedance2",
	"adobe2api-sora2",
	"adobe2api-veo31",
	"adobe2api-kling3",
	"adobe2api-kling-o3",
	"adobe2api-seedance2",
	"firefly-sora",
	"firefly-veo",
	"firefly-kling",
	"firefly-seedance",
}

func IsAdobe2APIImageOriginModel(model string) bool {
	return hasAdobe2APIPrefix(model, adobe2APIImageModelPrefixes)
}

func IsAdobe2APIVideoChatOriginModel(model string) bool {
	return hasAdobe2APIPrefix(model, adobe2APIVideoModelPrefixes)
}

func IsAdobe2APIImageRelay(info *relaycommon.RelayInfo) bool {
	if info == nil {
		return false
	}
	if !isAdobe2APIChannel(info) {
		return false
	}
	if IsAdobe2APIImageOriginModel(info.OriginModelName) {
		return true
	}
	if info.ChannelMeta != nil && IsAdobe2APIImageOriginModel(info.UpstreamModelName) {
		return true
	}
	return false
}

func IsAdobe2APIVideoChatRelay(info *relaycommon.RelayInfo) bool {
	if info == nil {
		return false
	}
	if !isAdobe2APIChannel(info) {
		return false
	}
	if IsAdobe2APIVideoChatOriginModel(info.OriginModelName) {
		return true
	}
	if info.ChannelMeta != nil && IsAdobe2APIVideoChatOriginModel(info.UpstreamModelName) {
		return true
	}
	return false
}

func isAdobe2APIChannel(info *relaycommon.RelayInfo) bool {
	if info == nil {
		return false
	}
	baseURL := ""
	if info.ChannelMeta != nil {
		if info.ChannelMeta.ChannelId == 75 {
			return true
		}
		baseURL = info.ChannelMeta.ChannelBaseUrl
	}
	baseURL = strings.TrimSpace(strings.ToLower(baseURL))
	if baseURL == "" {
		return false
	}
	if strings.Contains(baseURL, "adobe2api") || strings.Contains(baseURL, "gongju") {
		return true
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	return parsed.Host == "45.67.221.45:6001"
}

func hasAdobe2APIPrefix(model string, prefixes []string) bool {
	name := strings.ToLower(strings.TrimSpace(model))
	name = strings.TrimPrefix(name, "adobe2api/")
	name = strings.TrimPrefix(name, "adobe/")
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func resolveAdobe2APIUpstreamModel(info *relaycommon.RelayInfo, fallback string) string {
	if info != nil && info.ChannelMeta != nil && strings.TrimSpace(info.UpstreamModelName) != "" {
		return strings.TrimSpace(info.UpstreamModelName)
	}
	name := strings.TrimSpace(fallback)
	for _, prefix := range []string{"adobe2api/", "adobe/", "adobe2api-", "adobe-"} {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			return strings.TrimSpace(name[len(prefix):])
		}
	}
	if info != nil && strings.TrimSpace(info.OriginModelName) != "" {
		return strings.TrimSpace(info.OriginModelName)
	}
	return name
}

func ConvertAdobe2APIImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	modelName := resolveAdobe2APIUpstreamModel(info, request.Model)
	if modelName == "" {
		return nil, fmt.Errorf("model is required")
	}
	body := map[string]any{
		"model":  modelName,
		"prompt": request.Prompt,
	}
	if request.N != nil && *request.N > 0 {
		body["n"] = *request.N
	}
	imageSize := adobe2APIImageSize(request)
	if imageSize != "" {
		body["image_size"] = imageSize
		body["output_resolution"] = imageSize
	}
	if aspectRatio := adobe2APIAspectRatio(request); aspectRatio != "" {
		body["aspect_ratio"] = aspectRatio
	}
	refs, err := adobe2APIReferenceImages(c, request)
	if err != nil {
		return nil, err
	}
	if len(refs) > 0 {
		body["reference_images"] = refs
	}
	return body, nil
}

func adobe2APIImageSize(request dto.ImageRequest) string {
	for _, key := range []string{"image_size", "output_resolution"} {
		if raw, ok := request.Extra[key]; ok {
			if value := rawJSONString(raw); value != "" {
				return value
			}
		}
	}
	size := strings.ToUpper(strings.TrimSpace(request.Size))
	if size == "1K" || size == "2K" || size == "4K" {
		return size
	}
	switch strings.ToLower(strings.TrimSpace(request.Quality)) {
	case "high", "hd", "4k":
		return "4K"
	case "medium", "2k":
		return "2K"
	case "low", "standard", "1k":
		return "1K"
	default:
		return ""
	}
}

func adobe2APIAspectRatio(request dto.ImageRequest) string {
	if raw, ok := request.Extra["aspect_ratio"]; ok {
		if value := rawJSONString(raw); value != "" {
			return value
		}
	}
	value := strings.TrimSpace(request.Size)
	if strings.Contains(value, ":") {
		return value
	}
	switch strings.ToLower(value) {
	case "1024x1024", "2048x2048", "4096x4096":
		return "1:1"
	case "1536x1024":
		return "3:2"
	case "1024x1536":
		return "2:3"
	case "1792x1024", "1920x1080":
		return "16:9"
	case "1024x1792", "1080x1920":
		return "9:16"
	default:
		return ""
	}
}

func adobe2APIReferenceImages(c *gin.Context, request dto.ImageRequest) ([]string, error) {
	refs := make([]string, 0, 6)
	for _, key := range []string{"reference_images", "image_refs"} {
		if raw, ok := request.Extra[key]; ok {
			refs = append(refs, rawJSONStringList(raw)...)
		}
	}
	extracted, err := collectChatImageReferenceURLs(c, request)
	if err != nil {
		return nil, err
	}
	refs = append(refs, extracted...)
	return uniqueNonEmptyStrings(refs), nil
}

func ConvertAdobe2APIOpenAIChatRequest(c *gin.Context, request *dto.GeneralOpenAIRequest, info *relaycommon.RelayInfo) (any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	body := map[string]any{}
	if c != nil {
		if storage, err := common.GetBodyStorage(c); err == nil {
			if raw, readErr := storage.Bytes(); readErr == nil && len(raw) > 0 {
				_ = common.Unmarshal(raw, &body)
			}
		}
	}
	body["model"] = resolveAdobe2APIUpstreamModel(info, request.Model)
	body["messages"] = adobe2APIChatMessages(body, request)

	for _, key := range []string{"duration", "aspect_ratio", "resolution", "generate_audio", "reference_mode", "video_reference_mode"} {
		if value, ok := body[key]; ok {
			body[key] = value
		}
	}
	if refMode := strings.TrimSpace(rawMapString(body, "reference_mode")); refMode != "" {
		if strings.TrimSpace(rawMapString(body, "video_reference_mode")) == "" {
			body["video_reference_mode"] = refMode
		}
	}
	delete(body, "image_urls")
	delete(body, "image_url")
	delete(body, "stream")
	return body, nil
}

func adobe2APIChatMessages(body map[string]any, request *dto.GeneralOpenAIRequest) []dto.Message {
	messages := request.Messages
	imageURLs := make([]string, 0, 4)
	imageURLs = append(imageURLs, rawMapStringList(body, "image_urls")...)
	if single := rawMapString(body, "image_url"); single != "" {
		imageURLs = append(imageURLs, single)
	}
	imageURLs = uniqueNonEmptyStrings(imageURLs)
	if len(imageURLs) == 0 || len(messages) == 0 {
		return messages
	}
	last := len(messages) - 1
	content := make([]dto.MediaContent, 0, len(imageURLs)+1)
	if text := strings.TrimSpace(messages[last].StringContent()); text != "" {
		content = append(content, dto.MediaContent{Type: dto.ContentTypeText, Text: text})
	}
	for _, imageURL := range imageURLs {
		content = append(content, dto.MediaContent{
			Type: dto.ContentTypeImageURL,
			ImageUrl: &dto.MessageImageUrl{
				Url: imageURL,
			},
		})
	}
	messages[last].SetMediaContent(content)
	return messages
}

func rawJSONString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := common.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var value any
	if err := common.Unmarshal(raw, &value); err == nil {
		switch v := value.(type) {
		case float64:
			if v == float64(int64(v)) {
				return fmt.Sprintf("%d", int64(v))
			}
			return strings.TrimSpace(fmt.Sprintf("%v", v))
		case bool:
			return fmt.Sprintf("%t", v)
		}
	}
	return ""
}

func rawJSONStringList(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var single string
	if err := common.Unmarshal(raw, &single); err == nil {
		if single = strings.TrimSpace(single); single != "" {
			return []string{single}
		}
		return nil
	}
	var list []string
	if err := common.Unmarshal(raw, &list); err == nil {
		return uniqueNonEmptyStrings(list)
	}
	return nil
}

func rawMapString(body map[string]any, key string) string {
	if body == nil {
		return ""
	}
	switch v := body[key].(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return strings.TrimSpace(v.String())
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}

func rawMapStringList(body map[string]any, key string) []string {
	if body == nil {
		return nil
	}
	switch v := body[key].(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return uniqueNonEmptyStrings(out)
	case []string:
		return uniqueNonEmptyStrings(v)
	case string:
		parts := strings.FieldsFunc(v, func(r rune) bool {
			return r == ',' || r == '\n' || r == '\r'
		})
		return uniqueNonEmptyStrings(parts)
	default:
		return nil
	}
}

func uniqueNonEmptyStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
