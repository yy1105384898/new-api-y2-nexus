package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExtractGeminiPathModel(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "/v1beta/models/gemini-banana-2.0:generateContent", want: "gemini-banana-2.0"},
		{path: "/v1/models/gemini-2.0-flash:streamGenerateContent", want: "gemini-2.0-flash"},
		{path: "/v1/chat/completions", want: ""},
	}
	for _, tt := range tests {
		if got := ExtractGeminiPathModel(tt.path); got != tt.want {
			t.Fatalf("ExtractGeminiPathModel(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestRewriteGeminiRequestPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1beta/models/gemini-banana-2.0:generateContent", nil)

	if err := rewriteGeminiRequestPath(c, "byte-gemini-banana-2.0"); err != nil {
		t.Fatalf("rewriteGeminiRequestPath returned error: %v", err)
	}
	want := "/v1beta/models/byte-gemini-banana-2.0:generateContent"
	if c.Request.URL.Path != want {
		t.Fatalf("path = %q, want %q", c.Request.URL.Path, want)
	}
}
