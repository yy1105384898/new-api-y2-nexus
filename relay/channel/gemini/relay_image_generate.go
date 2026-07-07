package gemini

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

// ConvertImageRequestToGeminiGenerateContent 将 OpenAI Image 请求转为 Gemini generateContent body。
// 用于渠道 58/71：上游 /v1beta/models/gemini-3-pro-image-preview:generateContent。
func ConvertImageRequestToGeminiGenerateContent(request dto.ImageRequest) (*dto.GeminiChatRequest, error) {
	prompt := strings.TrimSpace(request.Prompt)
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}

	geminiRequest := &dto.GeminiChatRequest{
		Contents: []dto.GeminiChatContent{{
			Role:  "user",
			Parts: []dto.GeminiPart{{Text: prompt}},
		}},
		GenerationConfig: dto.GeminiChatGenerationConfig{
			ResponseModalities: []string{"TEXT", "IMAGE"},
		},
		SafetySettings: buildGeminiSafetySettings(),
	}

	imageConfig := map[string]any{}
	if aspect := resolveGeminiImageAspectRatio(request.Size); aspect != "" {
		imageConfig["aspectRatio"] = aspect
	}
	if imageSize := resolveGeminiImageSize(request.Quality); imageSize != "" {
		imageConfig["imageSize"] = imageSize
	}
	if len(imageConfig) > 0 {
		imageConfigBytes, err := common.Marshal(imageConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal image_config: %w", err)
		}
		geminiRequest.GenerationConfig.ImageConfig = imageConfigBytes
	}
	return geminiRequest, nil
}

func resolveGeminiImageAspectRatio(size string) string {
	value := strings.TrimSpace(size)
	if value == "" || strings.EqualFold(value, "auto") {
		return ""
	}
	if strings.Contains(value, ":") {
		return value
	}
	return ""
}

func resolveGeminiImageSize(quality string) string {
	switch strings.ToLower(strings.TrimSpace(quality)) {
	case "high", "hd", "4k":
		return "4K"
	case "medium", "2k":
		return "2K"
	case "low", "standard", "1k", "auto":
		return "1K"
	default:
		return ""
	}
}

func imageDataFromGeminiGenerateContent(response *dto.GeminiChatResponse) ([]dto.ImageData, error) {
	if response == nil {
		return nil, errors.New("empty response from Gemini API")
	}
	images := make([]dto.ImageData, 0, 2)
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData == nil || !strings.HasPrefix(part.InlineData.MimeType, "image") {
				continue
			}
			images = append(images, dto.ImageData{B64Json: part.InlineData.Data})
		}
	}
	if len(images) == 0 {
		return nil, errors.New("generateContent response has no image data")
	}
	return images, nil
}

// GeminiGenerateContentImageHandler 将 generateContent 出图响应转为 OpenAI ImageResponse。
func GeminiGenerateContentImageHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}
	service.CloseResponseBodyGracefully(resp)

	var geminiResponse dto.GeminiChatResponse
	if err := common.Unmarshal(responseBody, &geminiResponse); err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}
	if len(geminiResponse.Candidates) == 0 {
		if geminiResponse.PromptFeedback != nil && geminiResponse.PromptFeedback.BlockReason != nil {
			return nil, types.NewOpenAIError(
				fmt.Errorf("request blocked by Gemini API: %s", *geminiResponse.PromptFeedback.BlockReason),
				types.ErrorCodePromptBlocked,
				http.StatusBadRequest,
			)
		}
		return nil, types.NewOpenAIError(errors.New("empty response from Gemini API"), types.ErrorCodeEmptyResponse, http.StatusInternalServerError)
	}

	images, err := imageDataFromGeminiGenerateContent(&geminiResponse)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponse, http.StatusBadGateway)
	}

	rehosted, err := service.RehostImageDataForClient(c.Request.Context(), c.GetInt("id"), "", info.ChannelBaseUrl, info.OriginModelName, images, info.ImageClientWantsURL)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponse, http.StatusBadGateway)
	}

	imageResp := dto.ImageResponse{
		Created: common.GetTimestamp(),
		Data:    rehosted,
	}
	out, err := common.Marshal(imageResp)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}
	service.IOCopyBytesGracefully(c, resp, out)

	usage := buildUsageFromGeminiMetadata(geminiResponse.UsageMetadata, info.GetEstimatePromptTokens())
	return &usage, nil
}
