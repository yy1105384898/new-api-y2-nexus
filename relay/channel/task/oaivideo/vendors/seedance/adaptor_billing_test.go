package seedance

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

func multipartSeedanceContext(t *testing.T, duration string) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("model", "cy-sd1-seedance-2.0-4k")
	_ = writer.WriteField("prompt", "test")
	if duration != "" {
		_ = writer.WriteField("duration", duration)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1/videos", &body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return c
}

func TestValidateRequestRejectsMissingOairegboxDuration(t *testing.T) {
	c := multipartSeedanceContext(t, "")
	info := &relaycommon.RelayInfo{OriginModelName: "cy-sd1-seedance-2.0-4k"}
	if taskErr := (&TaskAdaptor{}).ValidateRequestAndSetAction(c, info); taskErr == nil {
		t.Fatal("expected missing duration error")
	}
}

func TestValidateRequestRejectsOutOfRangeOairegboxDuration(t *testing.T) {
	for _, duration := range []string{"1", "16"} {
		c := multipartSeedanceContext(t, duration)
		info := &relaycommon.RelayInfo{OriginModelName: "cy-sd1-seedance-2.0-4k"}
		if taskErr := (&TaskAdaptor{}).ValidateRequestAndSetAction(c, info); taskErr == nil {
			t.Fatalf("expected duration %s to be rejected", duration)
		}
	}
}

func TestValidateRequestAcceptsOairegboxDurationRange(t *testing.T) {
	c := multipartSeedanceContext(t, "15")
	info := &relaycommon.RelayInfo{OriginModelName: "cy-sd1-seedance-2.0-4k"}
	if taskErr := (&TaskAdaptor{}).ValidateRequestAndSetAction(c, info); taskErr != nil {
		t.Fatalf("unexpected error: %v", taskErr)
	}
}

func TestValidateRequestRejectsLeonardoMini8sDurationOverEight(t *testing.T) {
	for _, duration := range []string{"9", "15"} {
		c := multipartSeedanceContext(t, duration)
		info := &relaycommon.RelayInfo{OriginModelName: leonardoSeedanceMini8sModel}
		if taskErr := (&TaskAdaptor{}).ValidateRequestAndSetAction(c, info); taskErr == nil {
			t.Fatalf("expected duration %s to be rejected", duration)
		}
	}
}

func TestValidateRequestAcceptsLeonardoMini8sDurationAtMostEight(t *testing.T) {
	for _, duration := range []string{"", "4", "8"} {
		c := multipartSeedanceContext(t, duration)
		info := &relaycommon.RelayInfo{OriginModelName: leonardoSeedanceMini8sModel}
		if taskErr := (&TaskAdaptor{}).ValidateRequestAndSetAction(c, info); taskErr != nil {
			t.Fatalf("duration %s should be accepted: %v", duration, taskErr)
		}
	}
}
