package service

import (
	"context"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/dto"
)

func TestRewriteLoopbackUpstreamImageURL(t *testing.T) {
	got := RewriteLoopbackUpstreamImageURL("https://api.example.com", "http://127.0.0.1:3001/files/a.png")
	if !strings.HasPrefix(got, "https://api.example.com:3001/files/a.png") {
		t.Fatalf("rewrite = %q", got)
	}
	if RewriteLoopbackUpstreamImageURL("", "https://cdn.example.com/a.png") != "https://cdn.example.com/a.png" {
		t.Fatal("empty channel base should keep original url")
	}
}

func TestRehostSyncImageResponseBodySkipsInternalPrefixModel(t *testing.T) {
	body := []byte(`{"created":1,"data":[{"url":"https://upstream.example/a.png"}]}`)
	out, err := RehostSyncImageResponseBody(context.Background(), 1, "go2api-gpt-image-2-1k", "https://api.example.com", body, false)
	if err != nil {
		t.Fatalf("RehostSyncImageResponseBody: %v", err)
	}
	if string(out) != string(body) {
		t.Fatalf("internal prefixed model should passthrough body unchanged")
	}
}

func TestRehostImageDataURLsRequiresR2ForUpstreamURL(t *testing.T) {
	images := []dto.ImageData{{Url: "https://upstream.example/a.png"}}
	for _, model := range []string{"geek2-gpt-image-2-4k", "flux-pro-2", "cy-img1-gpt-image-2", "manju-gemini-banana-pro-1/2k"} {
		_, err := RehostImageDataURLs(context.Background(), 1, "task_test", "https://api.example.com", model, images)
		if err == nil || !strings.Contains(err.Error(), "R2 not configured") {
			t.Fatalf("model %s: expected R2 not configured error, got %v", model, err)
		}
	}
}

func TestRehostImageDataURLsRequiresR2ForGulieLoopbackURL(t *testing.T) {
	images := []dto.ImageData{{Url: "http://127.0.0.1:3001/files/a.png"}}
	_, err := RehostImageDataURLs(context.Background(), 1, "task_test", "http://gulie.204.group:25555", "cy-img1-gpt-image-2", images)
	if err == nil || !strings.Contains(err.Error(), "R2 not configured") {
		t.Fatalf("expected R2 not configured error, got %v", err)
	}
}

func TestRehostImageDataForClientB64RequiresR2(t *testing.T) {
	images := []dto.ImageData{{B64Json: "aGVsbG8="}}
	_, err := RehostImageDataForClient(context.Background(), 1, "task_test", "https://api.example.com", "cy-img1-gpt-image-2", images, true)
	if err == nil || !strings.Contains(err.Error(), "R2 not configured") {
		t.Fatalf("expected R2 not configured error, got %v", err)
	}
}

func TestDecodeImageDataItemDataURI(t *testing.T) {
	data, mime, err := DecodeImageDataItem(dto.ImageData{Url: "data:image/png;base64,aGVsbG8="})
	if err != nil {
		t.Fatalf("DecodeImageDataItem: %v", err)
	}
	if mime != "image/png" {
		t.Fatalf("mime = %q", mime)
	}
	if string(data) != "hello" {
		t.Fatalf("data = %q", data)
	}
}

func TestRehostTaskImageResultURLsRejectsUpstreamURLWithoutPolicy(t *testing.T) {
	images := []dto.ImageData{{Url: "https://upstream.example/a.png"}}
	_, err := RehostTaskImageResultURLs(context.Background(), 1, "task_test", "https://api.example.com", "go2api-gpt-image-2-1k", images)
	if err == nil || !strings.Contains(err.Error(), "upstream returned url without b64_json") {
		t.Fatalf("expected policy rejection, got %v", err)
	}
}

func TestRehostTaskImageResultURLsRequiresR2ForB64(t *testing.T) {
	images := []dto.ImageData{{B64Json: "aGVsbG8="}}
	_, err := RehostTaskImageResultURLs(context.Background(), 1, "task_test", "https://api.example.com", "cy-img1-gpt-image-2", images)
	if err == nil || !strings.Contains(err.Error(), "R2 not configured") {
		t.Fatalf("expected R2 not configured error, got %v", err)
	}
}

func TestRehostTaskImageResultURLsRequiresR2ForAcceptedURL(t *testing.T) {
	images := []dto.ImageData{{Url: "https://upstream.example/a.png"}}
	_, err := RehostTaskImageResultURLs(context.Background(), 1, "task_test", "https://api.example.com", "manju-gemini-banana-pro-1/2k", images)
	if err == nil || !strings.Contains(err.Error(), "R2 not configured") {
		t.Fatalf("expected R2 not configured error, got %v", err)
	}
}
