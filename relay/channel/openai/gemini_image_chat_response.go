package openai

import (
	"encoding/json"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
)

type chatImageCompatResponse struct {
	Created int64 `json:"created"`
	Choices []struct {
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage json.RawMessage `json:"usage,omitempty"`
}

func convertChatImageResponseBody(body []byte) ([]byte, bool, error) {
	var chatResp chatImageCompatResponse
	if err := common.Unmarshal(body, &chatResp); err != nil {
		return nil, false, nil
	}
	if len(chatResp.Choices) == 0 {
		return nil, false, nil
	}

	images := make([]dto.ImageData, 0)
	for _, choice := range chatResp.Choices {
		items := extractImageDataFromChatContent(choice.Message.Content)
		images = append(images, items...)
	}
	if len(images) == 0 {
		return nil, false, nil
	}

	payload := map[string]any{
		"created": chatResp.Created,
		"data":    images,
	}
	if len(chatResp.Usage) > 0 && string(chatResp.Usage) != "null" {
		var usage any
		if err := common.Unmarshal(chatResp.Usage, &usage); err == nil && usage != nil {
			payload["usage"] = usage
		}
	}

	converted, err := common.Marshal(payload)
	if err != nil {
		return nil, true, err
	}
	return converted, true, nil
}

func extractImageDataFromChatContent(raw json.RawMessage) []dto.ImageData {
	raw = json.RawMessage(strings.TrimSpace(string(raw)))
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var text string
	if err := common.Unmarshal(raw, &text); err == nil {
		return imageDataFromText(text)
	}

	var parts []struct {
		Type     string          `json:"type"`
		Text     string          `json:"text,omitempty"`
		ImageURL json.RawMessage `json:"image_url,omitempty"`
	}
	if err := common.Unmarshal(raw, &parts); err != nil {
		return nil
	}
	var images []dto.ImageData
	for _, part := range parts {
		images = append(images, imageDataFromText(part.Text)...)
		if len(part.ImageURL) == 0 {
			continue
		}
		if url := imageURLFromRaw(part.ImageURL); url != "" {
			images = append(images, imageDataFromURL(url))
		}
	}
	return images
}

func imageDataFromText(text string) []dto.ImageData {
	uris := extractDataImageURIs(text)
	out := make([]dto.ImageData, 0, len(uris))
	for _, uri := range uris {
		out = append(out, imageDataFromURL(uri))
	}
	return out
}

func imageURLFromRaw(raw json.RawMessage) string {
	var text string
	if err := common.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var obj struct {
		URL string `json:"url"`
	}
	if err := common.Unmarshal(raw, &obj); err == nil {
		return strings.TrimSpace(obj.URL)
	}
	return ""
}

func imageDataFromURL(url string) dto.ImageData {
	item := dto.ImageData{Url: url}
	if b64, ok := splitDataImageURI(url); ok {
		item.B64Json = b64
	}
	return item
}

func extractDataImageURIs(text string) []string {
	var out []string
	rest := text
	for {
		idx := strings.Index(rest, "data:image/")
		if idx < 0 {
			break
		}
		rest = rest[idx:]
		end := len(rest)
		for i, r := range rest {
			if i == 0 {
				continue
			}
			if r == ')' || r == '"' || r == '\'' || r == '<' || r == '>' || r == '\n' || r == '\r' || r == '\t' || r == ' ' {
				end = i
				break
			}
		}
		uri := strings.TrimSpace(rest[:end])
		if _, ok := splitDataImageURI(uri); ok {
			out = append(out, uri)
		}
		if end >= len(rest) {
			break
		}
		rest = rest[end:]
	}
	return out
}

func splitDataImageURI(uri string) (string, bool) {
	const marker = ";base64,"
	if !strings.HasPrefix(uri, "data:image/") {
		return "", false
	}
	idx := strings.Index(uri, marker)
	if idx < 0 {
		return "", false
	}
	b64 := strings.TrimSpace(uri[idx+len(marker):])
	if b64 == "" {
		return "", false
	}
	if strings.ContainsAny(b64, " \n\r\t") {
		return "", false
	}
	if strings.Contains(b64, ",") {
		return "", false
	}
	return b64, true
}
