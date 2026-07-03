package openai

import (
	"encoding/json"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
)

func shouldRouteImageRequestViaChat(model string) bool {
	name := strings.ToLower(strings.TrimSpace(model))
	return strings.Contains(name, "gemini-banana") ||
		strings.Contains(name, "nano-banana") ||
		strings.Contains(name, "flash-lite-image")
}

func convertImageRequestToChatCompletion(info *relaycommon.RelayInfo, request dto.ImageRequest) (*dto.GeneralOpenAIRequest, error) {
	if info != nil {
		info.RequestURLPath = "/v1/chat/completions"
		info.FinalRequestRelayFormat = types.RelayFormatOpenAI
	}

	stream := false
	body := &dto.GeneralOpenAIRequest{
		Model:    request.Model,
		Stream:   &stream,
		Messages: []dto.Message{{Role: "user", Content: request.Prompt}},
	}

	if extraBody := buildGeminiImageExtraBody(request); len(extraBody) > 0 {
		raw, err := json.Marshal(extraBody)
		if err != nil {
			return nil, err
		}
		body.ExtraBody = raw
	}

	return body, nil
}

func buildGeminiImageExtraBody(request dto.ImageRequest) map[string]any {
	imageConfig := map[string]any{}
	if aspectRatio := imageRequestAspectRatio(request.Size); aspectRatio != "" {
		imageConfig["aspect_ratio"] = aspectRatio
	}
	if imageSize := imageRequestImageSize(request.Quality); imageSize != "" {
		imageConfig["image_size"] = imageSize
	}
	if len(imageConfig) == 0 {
		return nil
	}
	return map[string]any{
		"google": map[string]any{
			"image_config": imageConfig,
		},
	}
}

func imageRequestAspectRatio(size string) string {
	size = strings.TrimSpace(size)
	if size == "" || size == "auto" {
		return ""
	}
	if strings.Contains(size, ":") {
		return size
	}
	switch size {
	case "1024x1024", "512x512", "256x256":
		return "1:1"
	case "1536x1024":
		return "3:2"
	case "1024x1536":
		return "2:3"
	case "1792x1024":
		return "16:9"
	case "1024x1792":
		return "9:16"
	}
	return ""
}

func imageRequestImageSize(quality string) string {
	switch strings.ToLower(strings.TrimSpace(quality)) {
	case "low":
		return "1K"
	case "medium":
		return "2K"
	case "high", "hd":
		return "4K"
	default:
		return ""
	}
}
