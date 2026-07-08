package image

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
)

func TestDowngradeImageSizeForChannel72(t *testing.T) {
	tests := []struct {
		size string
		want string
		ok   bool
	}{
		// OpenAI 官方 4K → 2K
		{size: "3840x2160", want: "2048x1152", ok: true},
		{size: "2160x3840", want: "1152x2048", ok: true},
		{size: "16:9-4k", want: "2048x1152", ok: true},
		{size: "4k", want: "2048x1152", ok: true},
		{size: "1:1-4k", want: "2048x2048", ok: true},
		// OpenAI 2K 上限内：不动
		{size: "2560x1440", want: "", ok: false},
		{size: "2048x1152", want: "", ok: false},
		{size: "1920x1080", want: "", ok: false},
		{size: "1024x1024", want: "", ok: false},
		{size: "16:9-2k", want: "", ok: false},
		{size: "1:1", want: "", ok: false},
		// 超 OpenAI 2K 像素上限：等比缩到 3686400 内
		{size: "4096x4096", want: "2048x2048", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			got, ok := downgradeImageSizeForChannel72(tt.size)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestScaleImagePixelsToOpenAI2KMax(t *testing.T) {
	// 3840x2160 走 direct map，这里测非标超大分辨率
	got, ok := scaleImagePixelsToOpenAI2KMax(3000, 2000)
	if !ok {
		t.Fatal("3000x2000 should scale")
	}
	w, h, parsed := parseImagePixelSize(got)
	if !parsed {
		t.Fatalf("parse %q", got)
	}
	if w*h > openAIImage2KMaxPixels {
		t.Fatalf("scaled area %d exceeds OpenAI 2K max %d", w*h, openAIImage2KMaxPixels)
	}
	if w%16 != 0 || h%16 != 0 {
		t.Fatalf("scaled size %dx%d not 16px aligned", w, h)
	}

	got, ok = scaleImagePixelsToOpenAI2KMax(2560, 1440)
	if ok || got != "" {
		t.Fatalf("2560x1440 is OpenAI 2K max, got %q ok=%v", got, ok)
	}
}

func TestApplyChannelImageSizeDowngrade(t *testing.T) {
	request := &dto.ImageRequest{Size: "3840x2160"}
	if !ApplyChannelImageSizeDowngrade(channelImageSize4KDowngradeID, request) {
		t.Fatal("expected downgrade")
	}
	if request.Size != "2048x1152" {
		t.Fatalf("size = %q", request.Size)
	}
}
