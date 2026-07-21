package omniv2v

import "testing"

func TestBuildUpstreamBody_ReferenceVideosToUpstream(t *testing.T) {
	t.Run("single reference video", func(t *testing.T) {
		out := buildUpstreamBody(map[string]interface{}{
			"prompt":           "restyle",
			"aspect_ratio":     "16:9",
			"reference_videos": []interface{}{"https://cdn.example.com/a.mp4"},
		}, "omni-fast-v2v", 0)
		if out["video_url"] != "https://cdn.example.com/a.mp4" {
			t.Fatalf("video_url = %#v", out["video_url"])
		}
		if _, exists := out["videos"]; exists {
			t.Fatalf("videos should not be set for single input: %#v", out)
		}
	})

	t.Run("dual reference videos", func(t *testing.T) {
		out := buildUpstreamBody(map[string]interface{}{
			"prompt": "blend",
			"reference_videos": []interface{}{
				"https://cdn.example.com/a.mp4",
				"https://cdn.example.com/b.mp4",
			},
			"reference_image_urls": []interface{}{"https://cdn.example.com/ref.jpg"},
		}, "omni-fast-v2v", 0)
		videos, ok := out["videos"].([]interface{})
		if !ok || len(videos) != 2 {
			t.Fatalf("videos = %#v", out["videos"])
		}
		images, ok := out["images"].([]interface{})
		if !ok || len(images) != 1 {
			t.Fatalf("images = %#v", out["images"])
		}
		if out["model"] != "omni-fast-v2v" {
			t.Fatalf("model = %#v", out["model"])
		}
	})
}

func TestIsRelay(t *testing.T) {
	if !IsRelay("cy-sd1-omni-v2v", "omni-fast-v2v") {
		t.Fatal("expected omni v2v relay match")
	}
	if IsRelay("omni-fast", "omni-fast") {
		t.Fatal("omni i2v should not match omni v2v relay")
	}
}
