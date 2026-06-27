package model_setting

import "testing"

func TestIsGeminiFlashImageModel(t *testing.T) {
	tests := []struct {
		model string
		want  bool
	}{
		{model: "gemini-3.1-flash-image-preview", want: true},
		{model: "nano-banana-pro", want: true},
		{model: "imagen-4.0-generate-001", want: false},
		{model: "gemini-2.0-flash", want: false},
	}
	for _, tt := range tests {
		if got := IsGeminiFlashImageModel(tt.model); got != tt.want {
			t.Fatalf("IsGeminiFlashImageModel(%q) = %v, want %v", tt.model, got, tt.want)
		}
	}
}
