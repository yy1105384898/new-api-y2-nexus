package gemini

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
)

func TestConvertImageRequestToGeminiGenerateContent4K(t *testing.T) {
	out, err := ConvertImageRequestToGeminiGenerateContent(dto.ImageRequest{
		Prompt:  "a cat",
		Size:    "1:1",
		Quality: "high",
	})
	if err != nil {
		t.Fatalf("ConvertImageRequestToGeminiGenerateContent: %v", err)
	}
	if len(out.Contents) != 1 || out.Contents[0].Parts[0].Text != "a cat" {
		t.Fatalf("unexpected contents: %+v", out.Contents)
	}
	if len(out.GenerationConfig.ResponseModalities) != 2 {
		t.Fatalf("modalities = %v", out.GenerationConfig.ResponseModalities)
	}
	var imageConfig map[string]any
	if err := common.Unmarshal(out.GenerationConfig.ImageConfig, &imageConfig); err != nil {
		t.Fatalf("unmarshal image_config: %v", err)
	}
	if imageConfig["aspectRatio"] != "1:1" {
		t.Fatalf("aspectRatio = %v", imageConfig["aspectRatio"])
	}
	if imageConfig["imageSize"] != "4K" {
		t.Fatalf("imageSize = %v", imageConfig["imageSize"])
	}
}

func TestGetRequestURLImagineImageUsesGenerateContent(t *testing.T) {
	adaptor := Adaptor{}
	info := &relaycommon.RelayInfo{
		ChannelBaseUrl:    "https://api.0lll0.cn",
		UpstreamModelName: "gemini-3-pro-image-preview",
		RelayMode:         relayconstant.RelayModeImagesGenerations,
	}
	url, err := adaptor.GetRequestURL(info)
	if err != nil {
		t.Fatalf("GetRequestURL: %v", err)
	}
	want := "https://api.0lll0.cn/v1beta/models/gemini-3-pro-image-preview:generateContent"
	if url != want {
		t.Fatalf("url = %q, want %q", url, want)
	}
}

func TestImageDataFromGeminiGenerateContent(t *testing.T) {
	images, err := imageDataFromGeminiGenerateContent(&dto.GeminiChatResponse{
		Candidates: []dto.GeminiChatCandidate{{
			Content: dto.GeminiChatContent{
				Parts: []dto.GeminiPart{{
					InlineData: &dto.GeminiInlineData{
						MimeType: "image/png",
						Data:     "abc123",
					},
				}},
			},
		}},
	})
	if err != nil {
		t.Fatalf("imageDataFromGeminiGenerateContent: %v", err)
	}
	if len(images) != 1 || images[0].B64Json != "abc123" {
		t.Fatalf("images = %+v", images)
	}
}

func TestResolveGeminiImageSize(t *testing.T) {
	if got := resolveGeminiImageSize("high"); got != "4K" {
		t.Fatalf("high -> %q", got)
	}
	if got := resolveGeminiImageSize("standard"); got != "1K" {
		t.Fatalf("standard -> %q", got)
	}
	if got := resolveGeminiImageSize(""); got != "" {
		t.Fatalf("empty -> %q", got)
	}
}

func TestConvertImageRequestToGeminiGenerateContentRequiresPrompt(t *testing.T) {
	_, err := ConvertImageRequestToGeminiGenerateContent(dto.ImageRequest{Prompt: "   "})
	if err == nil || !strings.Contains(err.Error(), "prompt is required") {
		t.Fatalf("err = %v", err)
	}
}
