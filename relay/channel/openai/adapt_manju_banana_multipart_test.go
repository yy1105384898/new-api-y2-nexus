package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/gin-gonic/gin"
)

func TestConvertManjuBananaMultipartEditsMarshalsUpstreamModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("model", "manju-gemini-banana-pro-4k"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("prompt", "make background blue"); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("image", "ref.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("fakepng")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body.Bytes()))
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	imageReq, err := helper.GetAndValidOpenAIImageRequest(c, relayconstant.RelayModeImagesEdits)
	if err != nil {
		t.Fatalf("parse multipart: %v", err)
	}

	c.Set("model_mapping", `{"manju-gemini-banana-pro-4k":"gemini-3.0-pro-image 4K"}`)
	common.SetContextKey(c, constant.ContextKeyOriginalModel, "manju-gemini-banana-pro-4k")

	info := &relaycommon.RelayInfo{
		RelayMode:       relayconstant.RelayModeImagesEdits,
		OriginModelName: "manju-gemini-banana-pro-4k",
		Request:         imageReq,
	}
	info.InitChannelMeta(c)

	request, err := common.DeepCopy(imageReq)
	if err != nil {
		t.Fatalf("deep copy: %v", err)
	}
	if err := helper.ModelMappedHelper(c, info, request); err != nil {
		t.Fatalf("model mapped: %v", err)
	}

	converted, err := ConvertManjuBananaImageRequest(c, info, *request)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	upstreamBody, ok := converted.(map[string]any)
	if !ok {
		t.Fatalf("expected map body, got %T", converted)
	}
	if upstreamBody["model"] != "gemini-3.0-pro-image 4K" {
		t.Fatalf("body model = %q", upstreamBody["model"])
	}
	if upstreamBody["stream"] != false {
		t.Fatalf("stream = %v", upstreamBody["stream"])
	}
	if _, ok := upstreamBody["extra_body"]; ok {
		t.Fatalf("extra_body should not be present: %v", upstreamBody["extra_body"])
	}

	jsonData, err := common.Marshal(converted)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	modelRaw, ok := payload["model"]
	if !ok || string(modelRaw) == `""` || string(modelRaw) == "null" {
		t.Fatalf("marshaled body missing model: %s", string(jsonData))
	}
	if got := string(modelRaw); got != `"gemini-3.0-pro-image 4K"` {
		t.Fatalf("marshaled model = %s", got)
	}

	imageRaw, ok := payload["image"]
	if !ok {
		t.Fatalf("missing image field: %s", string(jsonData))
	}
	var image string
	if err := json.Unmarshal(imageRaw, &image); err != nil {
		t.Fatalf("image: %v", err)
	}
	if len(image) < len("data:image/") || image[:len("data:image/")] != "data:image/" {
		t.Fatalf("expected image data URI, got %q", image)
	}
	if _, ok := payload["messages"]; ok {
		t.Fatalf("Image API payload must not contain chat messages: %s", string(jsonData))
	}

	// body must stay replayable after conversion path touched multipart
	if _, err := io.ReadAll(c.Request.Body); err != nil {
		t.Fatalf("read body: %v", err)
	}
}
