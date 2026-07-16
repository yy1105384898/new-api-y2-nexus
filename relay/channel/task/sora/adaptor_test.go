package sora

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

func TestBuildRequestURL_GrokImagineVideo(t *testing.T) {
	adaptor := &TaskAdaptor{baseURL: "http://sub2api:8091"}
	url, err := adaptor.BuildRequestURL(&relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelId:         192,
			UpstreamModelName: "grok-imagine-video",
		},
		TaskRelayInfo: &relaycommon.TaskRelayInfo{Action: constant.TaskActionGenerate},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "http://sub2api:8091/v1/videos/generations" {
		t.Fatalf("got %q", url)
	}
}

func TestBuildGrok2APIVideoJSON(t *testing.T) {
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)
	for key, value := range map[string]string{
		"prompt":          "animate this image",
		"seconds":         "8",
		"size":            "1792x1024",
		"resolution_name": "720p",
		"input_reference": "https://example.com/reference.png",
		"preset":          "normal",
	} {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	form, err := multipart.NewReader(&payload, writer.Boundary()).ReadForm(1024 * 1024)
	if err != nil {
		t.Fatal(err)
	}
	defer form.RemoveAll()

	body, err := buildGrok2APIVideoJSON(form, "grok-imagine-video")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"duration":"8"`, `"aspect_ratio":"16:9"`, `"resolution":"720p"`,
		`"image":{"url":"https://example.com/reference.png"}`,
	} {
		if !bytes.Contains(data, []byte(want)) {
			t.Fatalf("missing %s in %s", want, data)
		}
	}
	for _, forbidden := range []string{`"seconds"`, `"size"`, `"resolution_name"`, `"preset"`} {
		if bytes.Contains(data, []byte(forbidden)) {
			t.Fatalf("unexpected %s in %s", forbidden, data)
		}
	}
}

func TestBuildJSONGrokVideoBody(t *testing.T) {
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)
	if err := writer.WriteField("prompt", "animate this image"); err != nil {
		t.Fatal(err)
	}
	file, err := writer.CreateFormFile("input_reference", "reference.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = file.Write([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	form, err := multipart.NewReader(&payload, writer.Boundary()).ReadForm(1024 * 1024)
	if err != nil {
		t.Fatal(err)
	}
	defer form.RemoveAll()

	body, err := buildJSONGrokVideoBody(form, "grok-video-1.5")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"model":"grok-video-1.5"`)) ||
		!bytes.Contains(data, []byte(`"image_urls":["data:image/png;base64,`)) {
		t.Fatalf("unexpected JSON body: %s", data)
	}
}

func TestParseTaskResult_GZFormat(t *testing.T) {
	adaptor := &TaskAdaptor{}

	t.Run("running", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"running","videoUrl":null,"error":null}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusInProgress {
			t.Fatalf("expected IN_PROGRESS, got %s", result.Status)
		}
	})

	t.Run("succeeded with videoUrl", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"succeeded","videoUrl":"https://example.com/a.mp4","error":null}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusSuccess {
			t.Fatalf("expected SUCCESS, got %s", result.Status)
		}
		if result.Url != "https://example.com/a.mp4" {
			t.Fatalf("expected video url, got %q", result.Url)
		}
	})

	t.Run("failed with string error", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"failed","videoUrl":null,"error":"content policy violation"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusFailure {
			t.Fatalf("expected FAILURE, got %s", result.Status)
		}
		if result.Reason != "content policy violation" {
			t.Fatalf("expected error reason, got %q", result.Reason)
		}
	})

	t.Run("grok moderation without status", func(t *testing.T) {
		body := []byte(`{"code":"Client specified an invalid argument","error":"Generated video rejected by content moderation.","id":"task_upstream","task_id":"task_upstream","model":"grok-image-video"}`)
		result, err := adaptor.ParseTaskResult(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusFailure {
			t.Fatalf("expected FAILURE, got %s", result.Status)
		}
		if result.Reason != "Generated video rejected by content moderation." {
			t.Fatalf("expected moderation reason, got %q", result.Reason)
		}
	})

	t.Run("omni prefers data url over relative video_url", func(t *testing.T) {
		body := []byte(`{
			"id":"task_upstream",
			"status":"completed",
			"video_url":"/v1/videos/vid-4444bf370600/content",
			"data":[{"url":"https://download-2.oaibox.xyz/v1/videos/task_abc/content"}]
		}`)
		result, err := adaptor.ParseTaskResult(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusSuccess {
			t.Fatalf("expected SUCCESS, got %s", result.Status)
		}
		want := "https://download-2.oaibox.xyz/v1/videos/task_abc/content"
		if result.Url != want {
			t.Fatalf("expected %q, got %q", want, result.Url)
		}
	})

	t.Run("seedance accepts metadata url", func(t *testing.T) {
		body := []byte(`{"id":"task_upstream","status":"completed","metadata":{"url":"https://tmp.example.com/video.mp4"}}`)
		result, err := adaptor.ParseTaskResult(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusSuccess {
			t.Fatalf("expected SUCCESS, got %s", result.Status)
		}
		if result.Url != "https://tmp.example.com/video.mp4" {
			t.Fatalf("expected metadata url, got %q", result.Url)
		}
	})

	t.Run("unknown without error keeps polling", func(t *testing.T) {
		body := []byte(`{"created_at":1783042146,"id":"task_upstream","model":"cy-sd1-seedance-2.0-fast-480p","object":"video","progress":0,"status":"unknown","task_id":"task_upstream"}`)
		result, err := adaptor.ParseTaskResult(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusInProgress {
			t.Fatalf("expected IN_PROGRESS, got %s", result.Status)
		}
		if result.Progress == "100%" {
			t.Fatalf("unknown status without error should not be terminal")
		}
	})
}

func TestParseTaskResult_OpenAIFormat(t *testing.T) {
	adaptor := &TaskAdaptor{}

	result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"completed","usage":{"seconds":8}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != model.TaskStatusSuccess {
		t.Fatalf("expected SUCCESS, got %s", result.Status)
	}
	if result.CompletionTokens != 8 {
		t.Fatalf("expected 8 seconds, got %d", result.CompletionTokens)
	}
}

func TestAdjustBillingOnComplete_OAIREGBoxFallbackToRequestedSeconds(t *testing.T) {
	adaptor := &TaskAdaptor{}
	task := &model.Task{
		Quota: 400000,
		Properties: model.Properties{
			OriginModelName: "cy-sd1-seedance-2.0-mini-480p",
		},
		Data: []byte(`{"created_at":1783265344,"status":"completed","video_url":"https://example.com/a.mp4","data":[{"url":"https://example.com/a.mp4"}]}`),
		PrivateData: model.TaskPrivateData{
			BillingContext: &model.TaskBillingContext{
				ModelPrice:      0.2,
				GroupRatio:      1,
				OtherRatios:     map[string]float64{"seconds": 4, "size": 1},
				OriginModelName: "cy-sd1-seedance-2.0-mini-480p",
			},
		},
	}

	got := adaptor.AdjustBillingOnComplete(task, &relaycommon.TaskInfo{Status: model.TaskStatusSuccess})
	want := int(0.2 * float64(common.QuotaPerUnit) * 4)
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func TestAdjustBillingOnComplete_PrefersUpstreamUsageSeconds(t *testing.T) {
	adaptor := &TaskAdaptor{}
	task := &model.Task{
		Quota: 200000,
		Data:  []byte(`{"status":"completed","usage":{"seconds":8}}`),
		PrivateData: model.TaskPrivateData{
			BillingContext: &model.TaskBillingContext{
				ModelPrice:  0.2,
				GroupRatio:  1,
				OtherRatios: map[string]float64{"seconds": 4},
			},
		},
	}

	got := adaptor.AdjustBillingOnComplete(task, &relaycommon.TaskInfo{Status: model.TaskStatusSuccess})
	want := int(0.2 * float64(common.QuotaPerUnit) * 8)
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func TestUsageSecondsFromTaskData(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int
	}{
		{
			name: "usage object",
			data: `{"status":"completed","usage":{"seconds":6}}`,
			want: 6,
		},
		{
			name: "top level seconds string",
			data: `{"status":"completed","seconds":"5"}`,
			want: 5,
		},
		{
			name: "oairegbox without seconds",
			data: `{"status":"completed","video_url":"https://example.com/a.mp4"}`,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usageSecondsFromTaskData([]byte(tt.data)); got != tt.want {
				t.Fatalf("usageSecondsFromTaskData() = %d, want %d", got, tt.want)
			}
		})
	}
}
