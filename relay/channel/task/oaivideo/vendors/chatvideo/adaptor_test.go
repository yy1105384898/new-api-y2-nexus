package chatvideo

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

func TestBuildRequestBodyConvertsCanonicalVideoRequestToChat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewBufferString(`{}`))
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Model:         "sora-2",
		Prompt:        "make a video",
		Images:        []string{"https://example.com/ref.png"},
		Duration:      8,
		AspectRatio:   "16:9",
		Resolution:    "720p",
		GenerateAudio: common.GetPointer(true),
	})

	reader, err := (&TaskAdaptor{}).BuildRequestBody(c, &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{UpstreamModelName: "sora-2"},
	})
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	body, _ := io.ReadAll(reader)
	var got map[string]any
	if err := common.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if got["model"] != "sora-2" || got["stream"] != true {
		t.Fatalf("unexpected request: %#v", got)
	}
	if got["duration"] != float64(8) || got["aspect_ratio"] != "16:9" || got["resolution"] != "720p" || got["generate_audio"] != true {
		t.Fatalf("video parameters were not preserved: %#v", got)
	}
}

func TestDoResponseConvertsSSEVideoURLToCompletedTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: io.NopCloser(bytes.NewBufferString(
			"data: {\"choices\":[{\"delta\":{\"content\":\"[video](https://example.com/out.mp4)\"}}]}\n\n" +
				"data: [DONE]\n\n",
		)),
	}

	upstreamID, taskData, taskErr := (&TaskAdaptor{}).DoResponse(c, resp, &relaycommon.RelayInfo{
		TaskRelayInfo: &relaycommon.TaskRelayInfo{PublicTaskID: "task_public"},
		ChannelMeta:   &relaycommon.ChannelMeta{UpstreamModelName: "sora-2"},
	})
	if taskErr != nil {
		t.Fatalf("response error: %v", taskErr)
	}
	if upstreamID != "task_public" || !bytes.Contains(taskData, []byte(`"status":"completed"`)) || !bytes.Contains(taskData, []byte("https://example.com/out.mp4")) {
		t.Fatalf("unexpected task: id=%q data=%s", upstreamID, taskData)
	}
}

func TestIsRelayUsesInternalVideoRoutePrefix(t *testing.T) {
	if !IsRelay("cy-vid2-sora-2") || !IsRelay("cy-sd1-grok-video") || IsRelay("cy-sd1-grok-video-cli") || IsRelay("cy-gv1-grok-video") {
		t.Fatal("chat video route matching is incorrect")
	}
}
