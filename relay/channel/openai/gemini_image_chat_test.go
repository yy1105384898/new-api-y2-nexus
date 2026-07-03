package openai

import (
	"encoding/json"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
)

func TestConvertGeminiBananaImageRequestToChatCompletion(t *testing.T) {
	info := &relaycommon.RelayInfo{
		RequestURLPath: "/v1/images/generations",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "gemini-banana-2.0",
		},
	}
	req := dto.ImageRequest{
		Model:   "gemini-banana-2.0",
		Prompt:  "draw a cat",
		Size:    "1792x1024",
		Quality: "high",
	}

	converted, err := convertImageRequestToChatCompletion(info, req)
	if err != nil {
		t.Fatalf("convertImageRequestToChatCompletion() error = %v", err)
	}
	if info.RequestURLPath != "/v1/chat/completions" {
		t.Fatalf("RequestURLPath = %q", info.RequestURLPath)
	}
	if info.FinalRequestRelayFormat != types.RelayFormatOpenAI {
		t.Fatalf("FinalRequestRelayFormat = %q", info.FinalRequestRelayFormat)
	}
	if converted.Model != "gemini-banana-2.0" {
		t.Fatalf("model = %q", converted.Model)
	}
	if len(converted.Messages) != 1 || converted.Messages[0].Content != "draw a cat" {
		t.Fatalf("messages = %#v", converted.Messages)
	}

	var extra map[string]map[string]map[string]string
	if err := json.Unmarshal(converted.ExtraBody, &extra); err != nil {
		t.Fatalf("extra_body unmarshal: %v", err)
	}
	config := extra["google"]["image_config"]
	if config["aspect_ratio"] != "16:9" || config["image_size"] != "4K" {
		t.Fatalf("image_config = %#v", config)
	}
}

func TestConvertChatImageResponseBody(t *testing.T) {
	body := []byte(`{
		"created": 1783040000,
		"choices": [{
			"message": {
				"content": "![image](data:image/png;base64,QUJDRA==)"
			}
		}],
		"usage": {"prompt_tokens": 1, "completion_tokens": 2, "total_tokens": 3}
	}`)

	converted, ok, err := convertChatImageResponseBody(body)
	if err != nil {
		t.Fatalf("convertChatImageResponseBody() error = %v", err)
	}
	if !ok {
		t.Fatal("convertChatImageResponseBody() ok = false")
	}

	var payload struct {
		Created int64           `json:"created"`
		Data    []dto.ImageData `json:"data"`
		Usage   dto.Usage       `json:"usage"`
	}
	if err := json.Unmarshal(converted, &payload); err != nil {
		t.Fatalf("converted unmarshal: %v", err)
	}
	if payload.Created != 1783040000 {
		t.Fatalf("created = %d", payload.Created)
	}
	if len(payload.Data) != 1 || payload.Data[0].B64Json != "QUJDRA==" {
		t.Fatalf("data = %#v", payload.Data)
	}
	if payload.Usage.TotalTokens != 3 {
		t.Fatalf("usage = %#v", payload.Usage)
	}
}
