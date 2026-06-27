package gemini

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/gin-gonic/gin"
)

func TestConvertFlashImageRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/images/generations", nil)

	req := dto.ImageRequest{
		Model:  "gemini-banana-2.0",
		Prompt: "a red apple",
		Size:   "1536x1024",
	}
	got, err := convertFlashImageRequest(c, req)
	if err != nil {
		t.Fatalf("convertFlashImageRequest returned error: %v", err)
	}
	chatReq, ok := got.(*dto.GeminiChatRequest)
	if !ok {
		t.Fatalf("expected *dto.GeminiChatRequest, got %T", got)
	}
	if len(chatReq.Contents) != 1 || len(chatReq.Contents[0].Parts) == 0 {
		t.Fatalf("expected user content with prompt")
	}
	if chatReq.Contents[0].Parts[0].Text != "a red apple" {
		t.Fatalf("prompt = %q", chatReq.Contents[0].Parts[0].Text)
	}
	if len(chatReq.GenerationConfig.ResponseModalities) != 2 {
		t.Fatalf("expected responseModalities, got %#v", chatReq.GenerationConfig.ResponseModalities)
	}
	var imageConfig map[string]interface{}
	if err := json.Unmarshal(chatReq.GenerationConfig.ImageConfig, &imageConfig); err != nil {
		t.Fatalf("unmarshal imageConfig: %v", err)
	}
	if imageConfig["aspectRatio"] != "3:2" {
		t.Fatalf("aspectRatio = %#v", imageConfig["aspectRatio"])
	}
}

func TestExtractGeminiChatImageData(t *testing.T) {
	response := &dto.GeminiChatResponse{
		Candidates: []dto.GeminiChatCandidate{
			{
				Content: dto.GeminiChatContent{
					Parts: []dto.GeminiPart{
						{InlineData: &dto.GeminiInlineData{MimeType: "image/png", Data: "abc123"}},
					},
				},
			},
		},
	}
	images := extractGeminiChatImageData(response)
	if len(images) != 1 || images[0].B64Json != "abc123" {
		t.Fatalf("unexpected images: %#v", images)
	}
}
