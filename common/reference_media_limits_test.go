package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReferenceMediaLimitsSeedDataSync(t *testing.T) {
	seedPath := filepath.Join("..", "scripts", "seed_data", "model_ui_params_video.json")
	raw, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("read seed data: %v", err)
	}

	var payload struct {
		Profiles []struct {
			ID              string                 `json:"id"`
			ReferenceLimits map[string]interface{} `json:"referenceLimits"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal seed data: %v", err)
	}

	for _, profile := range payload.Profiles {
		limits := profile.ReferenceLimits
		if len(limits) == 0 {
			continue
		}
		if value, ok := limits["imageMaxBytes"]; ok {
			assertLimitEquals(t, profile.ID, "imageMaxBytes", value, ReferenceImageMaxBytes)
		}
		if value, ok := limits["videoMaxBytes"]; ok {
			assertLimitEquals(t, profile.ID, "videoMaxBytes", value, ReferenceVideoMaxBytes)
		}
		if value, ok := limits["audioMaxBytes"]; ok {
			assertLimitEquals(t, profile.ID, "audioMaxBytes", value, ReferenceAudioMaxBytes)
		}
	}
}

func assertLimitEquals(t *testing.T, profileID, key string, got any, want int64) {
	t.Helper()
	asFloat, ok := got.(float64)
	if !ok {
		t.Fatalf("profile %s %s expected number, got %T", profileID, key, got)
	}
	if int64(asFloat) != want {
		t.Fatalf("profile %s %s = %d, want %d", profileID, key, int64(asFloat), want)
	}
}

func TestReferenceImageTooLargeDetail(t *testing.T) {
	if got := ReferenceImageTooLargeDetail(); got != "image too large, max 30MB" {
		t.Fatalf("ReferenceImageTooLargeDetail() = %q", got)
	}
}

func TestFormatReferenceByteLimitMessageZH(t *testing.T) {
	got := FormatReferenceByteLimitMessageZH("image", ReferenceImageMaxBytes)
	if got != "参考图超过 30MB，请压缩后再上传" {
		t.Fatalf("FormatReferenceByteLimitMessageZH() = %q", got)
	}
}
