package omniv2v

import (
	"strings"

	oaivideo "github.com/QuantumNous/new-api/relay/channel/task/oaivideo/shared"
)

const (
	flatKeyReferenceVideos    = "reference_videos"
	flatKeyReferenceImageURLs = "reference_image_urls"
)

func buildUpstreamBody(body map[string]interface{}, upstreamModel string, duration int) map[string]interface{} {
	out := map[string]interface{}{
		"model":  strings.TrimSpace(upstreamModel),
		"prompt": strings.TrimSpace(oaivideo.AsString(body["prompt"])),
	}
	if aspectRatio := strings.TrimSpace(oaivideo.AsString(body["aspect_ratio"])); aspectRatio != "" {
		out["aspect_ratio"] = aspectRatio
	}
	if duration > 0 {
		out["seconds"] = duration
	}

	refVideos := collectReferenceVideos(body)
	switch len(refVideos) {
	case 1:
		out["video_url"] = refVideos[0]
	case 2:
		out["videos"] = []interface{}{refVideos[0], refVideos[1]}
	}

	if refImages := collectReferenceImages(body); len(refImages) > 0 {
		out["images"] = stringSliceToInterface(refImages)
	}
	return out
}

func collectReferenceVideos(body map[string]interface{}) []string {
	if body == nil {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 2)
	add := func(url string) {
		url = strings.TrimSpace(url)
		if url == "" {
			return
		}
		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}
		out = append(out, url)
	}
	for _, url := range oaivideo.CollectStringList(body[flatKeyReferenceVideos]) {
		add(url)
	}
	for _, url := range oaivideo.CollectStringList(body["videos"]) {
		add(url)
	}
	if url := strings.TrimSpace(oaivideo.AsString(body["video_url"])); url != "" {
		add(url)
	}
	if url := strings.TrimSpace(oaivideo.AsString(body["video"])); url != "" {
		add(url)
	}
	if len(out) > 2 {
		out = out[:2]
	}
	return out
}

func collectReferenceImages(body map[string]interface{}) []string {
	if body == nil {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 2)
	add := func(url string) {
		url = strings.TrimSpace(url)
		if url == "" {
			return
		}
		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}
		out = append(out, url)
	}
	for _, url := range oaivideo.CollectStringList(body[flatKeyReferenceImageURLs]) {
		add(url)
	}
	for _, url := range oaivideo.CollectStringList(body["images"]) {
		add(url)
	}
	for _, url := range oaivideo.CollectStringList(body["image_urls"]) {
		add(url)
	}
	if len(out) > 2 {
		out = out[:2]
	}
	return out
}

func stringSliceToInterface(urls []string) []interface{} {
	out := make([]interface{}, len(urls))
	for i, url := range urls {
		out[i] = url
	}
	return out
}
