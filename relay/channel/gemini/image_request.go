package gemini

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

var geminiMarkdownImageRE = regexp.MustCompile(`!\[[^\]]*\]\((data:image/[^;]+;base64,[^)]+)\)`)

func convertFlashImageRequest(c *gin.Context, request dto.ImageRequest) (*dto.GeminiChatRequest, error) {
	prompt := strings.TrimSpace(request.Prompt)
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	parts := []dto.GeminiPart{{Text: prompt}}
	refParts, err := buildGeminiImageReferenceParts(c, request)
	if err != nil {
		return nil, err
	}
	parts = append(parts, refParts...)

	maskParts, err := buildGeminiMaskReferenceParts(c, request)
	if err != nil {
		return nil, err
	}
	parts = append(parts, maskParts...)

	aspectRatio, imageSize := mapImageRequestDimensions(request)
	imageConfig := map[string]interface{}{}
	if aspectRatio != "" && aspectRatio != "auto" {
		imageConfig["aspectRatio"] = aspectRatio
	}
	if imageSize != "" {
		imageConfig["imageSize"] = imageSize
	}

	generationConfig := dto.GeminiChatGenerationConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	if len(imageConfig) > 0 {
		imageConfigBytes, err := json.Marshal(imageConfig)
		if err != nil {
			return nil, err
		}
		generationConfig.ImageConfig = imageConfigBytes
	}

	return &dto.GeminiChatRequest{
		Contents: []dto.GeminiChatContent{
			{
				Role:  "user",
				Parts: parts,
			},
		},
		GenerationConfig: generationConfig,
	}, nil
}

func mapImageRequestDimensions(request dto.ImageRequest) (aspectRatio string, imageSize string) {
	aspectRatio = "1:1"
	size := strings.TrimSpace(request.Size)
	if size != "" {
		if strings.Contains(size, ":") {
			aspectRatio = size
		} else {
			switch size {
			case "256x256", "512x512", "1024x1024":
				aspectRatio = "1:1"
			case "1536x1024":
				aspectRatio = "3:2"
			case "1024x1536":
				aspectRatio = "2:3"
			case "1024x1792":
				aspectRatio = "9:16"
			case "1792x1024":
				aspectRatio = "16:9"
			default:
				aspectRatio = size
			}
		}
	}

	if request.Quality != "" {
		switch strings.ToLower(strings.TrimSpace(request.Quality)) {
		case "hd", "high", "4k", "2k":
			imageSize = "2K"
			if strings.EqualFold(request.Quality, "4k") || strings.EqualFold(request.Quality, "high") {
				imageSize = "4K"
			}
		case "medium":
			imageSize = "2K"
		case "standard", "low", "auto", "1k":
			imageSize = "1K"
		}
	}
	return aspectRatio, imageSize
}

func buildGeminiImageReferenceParts(c *gin.Context, request dto.ImageRequest) ([]dto.GeminiPart, error) {
	refs, err := collectImageReferenceStrings(request.Image, request.Images)
	if err != nil {
		return nil, err
	}
	return fileSourcesToGeminiParts(c, refs)
}

func buildGeminiMaskReferenceParts(c *gin.Context, request dto.ImageRequest) ([]dto.GeminiPart, error) {
	if len(request.Mask) == 0 {
		return nil, nil
	}
	refs, err := collectImageReferenceStrings(request.Mask, nil)
	if err != nil {
		return nil, err
	}
	parts, err := fileSourcesToGeminiParts(c, refs)
	if err != nil {
		return nil, err
	}
	if len(parts) > 0 {
		parts = append([]dto.GeminiPart{{Text: "以下图片为蒙版，请按蒙版区域进行编辑。"}}, parts...)
	}
	return parts, nil
}

func collectImageReferenceStrings(imageField json.RawMessage, imagesField json.RawMessage) ([]string, error) {
	out := make([]string, 0, 4)
	if len(imageField) > 0 {
		refs, err := parseImageReferenceField(imageField)
		if err != nil {
			return nil, err
		}
		out = append(out, refs...)
	}
	if len(imagesField) > 0 {
		refs, err := parseImageReferenceField(imagesField)
		if err != nil {
			return nil, err
		}
		out = append(out, refs...)
	}
	return out, nil
}

func parseImageReferenceField(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			return nil, nil
		}
		return []string{single}, nil
	}
	var many []string
	if err := json.Unmarshal(raw, &many); err == nil {
		out := make([]string, 0, len(many))
		for _, item := range many {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
			}
		}
		return out, nil
	}
	return nil, fmt.Errorf("image reference must be a string or string array")
}

func fileSourcesToGeminiParts(c *gin.Context, refs []string) ([]dto.GeminiPart, error) {
	parts := make([]dto.GeminiPart, 0, len(refs))
	for _, ref := range refs {
		source := types.NewFileSourceFromData(ref, "")
		base64Data, mimeType, err := service.GetBase64Data(c, source, "formatting image for Gemini image generation")
		if err != nil {
			return nil, fmt.Errorf("load image reference failed: %w", err)
		}
		if _, ok := geminiSupportedMimeTypes[strings.ToLower(mimeType)]; !ok {
			return nil, fmt.Errorf("mime type is not supported by Gemini: %s", mimeType)
		}
		parts = append(parts, dto.GeminiPart{
			InlineData: &dto.GeminiInlineData{
				MimeType: mimeType,
				Data:     base64Data,
			},
		})
	}
	return parts, nil
}

func extractGeminiChatImageData(response *dto.GeminiChatResponse) []dto.ImageData {
	if response == nil {
		return nil
	}
	images := make([]dto.ImageData, 0, 1)
	seen := make(map[string]struct{})
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && strings.HasPrefix(part.InlineData.MimeType, "image") {
				data := strings.TrimSpace(part.InlineData.Data)
				if data == "" {
					continue
				}
				if _, ok := seen[data]; ok {
					continue
				}
				seen[data] = struct{}{}
				images = append(images, dto.ImageData{B64Json: data})
			}
			if part.Text == "" {
				continue
			}
			for _, match := range geminiMarkdownImageRE.FindAllStringSubmatch(part.Text, -1) {
				if len(match) < 2 {
					continue
				}
				format, base64String, err := service.DecodeBase64FileData(match[1])
				if err != nil || base64String == "" {
					continue
				}
				if _, ok := seen[base64String]; ok {
					continue
				}
				seen[base64String] = struct{}{}
				_ = format
				images = append(images, dto.ImageData{B64Json: base64String})
			}
		}
	}
	return images
}
