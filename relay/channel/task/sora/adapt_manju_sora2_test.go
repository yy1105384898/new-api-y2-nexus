package sora

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

func TestIsManjuSora2Relay(t *testing.T) {
	if !IsManjuSora2Relay("manju-openai-sora2", "sora2") {
		t.Fatal("expected manju sora2 relay")
	}
	if IsManjuSora2Relay("sora-2", "sora-2") {
		t.Fatal("expected standard sora not manju")
	}
}

func TestConvertManjuSora2ChatBody(t *testing.T) {
	out, err := ConvertManjuSora2ChatBody(map[string]interface{}{
		"prompt":  "cat on beach",
		"seconds": "8",
		"size":    "1280x720",
	}, "sora2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["model"] != "sora2" {
		t.Fatalf("expected upstream model sora2, got %v", out["model"])
	}
	if out["stream"] != false {
		t.Fatalf("expected stream false, got %v", out["stream"])
	}
	if out["sora2_ratio"] != "16:9" {
		t.Fatalf("expected sora2_ratio 16:9, got %v", out["sora2_ratio"])
	}
	if out["sora2_duration"] != "8" {
		t.Fatalf("expected sora2_duration 8, got %v", out["sora2_duration"])
	}
	msgs, ok := out["messages"].([]map[string]interface{})
	if !ok || len(msgs) == 0 || msgs[0]["content"] != "cat on beach" {
		t.Fatalf("expected messages with prompt, got %v", out["messages"])
	}
}

func TestParseTaskResult_ManjuSora2Succeeded(t *testing.T) {
	adaptor := &TaskAdaptor{}
	body := []byte(`{
		"id":"sora2-e949ceaef92c",
		"platform":"sora2",
		"status":"succeeded",
		"progress":100,
		"properties":{"duration":"8","aspect_ratio":"16:9","output_resolution":"720p"},
		"raw_data":{"video_url":"https://dlff.manjuapi.com/files/demo.mp4","video_urls":["https://dlff.manjuapi.com/files/demo.mp4"]},
		"video":{"url":"https://dlff.manjuapi.com/files/demo.mp4"}
	}`)
	result, err := adaptor.ParseTaskResult(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != model.TaskStatusSuccess {
		t.Fatalf("expected SUCCESS, got %s", result.Status)
	}
	want := "https://dlff.manjuapi.com/files/demo.mp4"
	if result.Url != want {
		t.Fatalf("expected url %q, got %q", want, result.Url)
	}
	if result.CompletionTokens != 8 {
		t.Fatalf("expected 8 seconds, got %d", result.CompletionTokens)
	}
}

func TestParseTaskResult_ManjuSora2Running(t *testing.T) {
	adaptor := &TaskAdaptor{}
	body := []byte(`{"id":"sora2-abc","platform":"sora2","status":"running","progress":13}`)
	result, err := adaptor.ParseTaskResult(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != model.TaskStatusInProgress {
		t.Fatalf("expected IN_PROGRESS, got %s", result.Status)
	}
	if result.Progress != "13%" {
		t.Fatalf("expected 13%%, got %q", result.Progress)
	}
}

func TestParseResponseTask_ManjuNestedData(t *testing.T) {
	body := []byte(`{"id":"sora2-abc","platform":"sora2","status":"running","progress":13,"data":{"id":1,"data":{"video_url":""}}}`)
	res, err := parseResponseTask(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "sora2-abc" {
		t.Fatalf("expected id sora2-abc, got %q", res.ID)
	}
	if res.Status != "running" {
		t.Fatalf("expected running, got %q", res.Status)
	}
}

func TestBuildOpenAIVideoCreateResponse(t *testing.T) {
	info := &relaycommon.RelayInfo{
		PublicTaskID:    "task_public",
		OriginModelName: "manju-openai-sora2",
	}
	out := buildOpenAIVideoCreateResponse(info, responseTask{Status: "running", Progress: 13, Seconds: "8"}, nil)
	if out["id"] != "task_public" {
		t.Fatalf("expected public task id")
	}
	if out["status"] != "in_progress" {
		t.Fatalf("expected in_progress, got %v", out["status"])
	}
}

func TestExtractManjuSoraVideoURL(t *testing.T) {
	body := []byte(`{"metadata":{"url":"https://example.com/a.mp4"}}`)
	if got := extractManjuSoraVideoURL(body); got != "https://example.com/a.mp4" {
		t.Fatalf("expected metadata url, got %q", got)
	}
}

func TestParseTaskResult_ManjuSora2FailedWithMessage(t *testing.T) {
	adaptor := &TaskAdaptor{}
	body := []byte(`{
		"id":"sora2-failed001",
		"platform":"sora2",
		"status":"failed",
		"message":"某张上传的参考图未通过平台内容审核（常见于含可识别真人肖像或敏感内容）；重试无效，请更换涉及的参考图后重试"
	}`)
	result, err := adaptor.ParseTaskResult(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != model.TaskStatusFailure {
		t.Fatalf("expected FAILURE, got %s", result.Status)
	}
	want := "某张上传的参考图未通过平台内容审核（常见于含可识别真人肖像或敏感内容）；重试无效，请更换涉及的参考图后重试"
	if result.Reason != want {
		t.Fatalf("expected reason %q, got %q", want, result.Reason)
	}
}

func TestBuildManjuSoraOpenAIErrorResponse(t *testing.T) {
	body := []byte(`{"id":"sora2-failed001","platform":"sora2","status":"failed","message":"审核失败"}`)
	out, ok := BuildManjuSoraOpenAIErrorResponse(body)
	if !ok {
		t.Fatal("expected error conversion")
	}
	if !strings.Contains(string(out), "审核失败") {
		t.Fatalf("expected message in output, got %s", string(out))
	}
}
