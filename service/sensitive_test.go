package service

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
)

func TestPromptSensitiveRejection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	prevEnabled := setting.CheckSensitiveEnabled
	prevPrompt := setting.CheckSensitiveOnPromptEnabled
	prevWords := append([]string(nil), setting.SensitiveWords...)
	t.Cleanup(func() {
		setting.CheckSensitiveEnabled = prevEnabled
		setting.CheckSensitiveOnPromptEnabled = prevPrompt
		setting.SensitiveWords = prevWords
	})

	setting.CheckSensitiveEnabled = true
	setting.CheckSensitiveOnPromptEnabled = true
	setting.SensitiveWords = []string{"雷管炸弹", "换脸"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	rejected, apiErr := PromptSensitiveRejection(c, "生成圆柱体雷管炸弹")
	if !rejected || apiErr == nil {
		t.Fatalf("expected sensitive rejection")
	}
	if got := apiErr.Error(); got != ContentPolicyMessageEN {
		t.Fatalf("PromptSensitiveRejection() = %q, want %q", got, ContentPolicyMessageEN)
	}

	rejected, apiErr = PromptSensitiveRejection(c, "生成一只鸟异兽")
	if rejected || apiErr != nil {
		t.Fatalf("expected no rejection for safe prompt")
	}
}

func TestTaskErrorIfSensitivePrompt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	prevEnabled := setting.CheckSensitiveEnabled
	prevPrompt := setting.CheckSensitiveOnPromptEnabled
	prevWords := append([]string(nil), setting.SensitiveWords...)
	t.Cleanup(func() {
		setting.CheckSensitiveEnabled = prevEnabled
		setting.CheckSensitiveOnPromptEnabled = prevPrompt
		setting.SensitiveWords = prevWords
	})

	setting.CheckSensitiveEnabled = true
	setting.CheckSensitiveOnPromptEnabled = true
	setting.SensitiveWords = []string{"锁定人脸"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/images/generations", nil)
	c.Request.Header.Set("X-Cangyuan-Client", "infinite-canvas")

	taskErr := TaskErrorIfSensitivePrompt(c, "需要锁定参考我的人脸信息")
	if taskErr == nil {
		t.Fatal("expected task error")
	}
	if !taskErr.LocalError {
		t.Fatal("expected local task error")
	}
	if taskErr.Message != ContentPolicyMessageZH {
		t.Fatalf("TaskErrorIfSensitivePrompt() = %q, want %q", taskErr.Message, ContentPolicyMessageZH)
	}
}
