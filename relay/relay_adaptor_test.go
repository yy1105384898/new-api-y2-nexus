package relay

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
)

func TestGetTaskAdaptorSupportsXaiVideo(t *testing.T) {
	adaptor := GetTaskAdaptor(constant.TaskPlatform("48"))
	if adaptor == nil {
		t.Fatal("xAI channel should use the OpenAI video router")
	}
}
