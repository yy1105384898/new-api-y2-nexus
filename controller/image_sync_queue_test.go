package controller

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func imageSyncQueueTestContext(t *testing.T, body string, contentType string) *gin.Context {
	t.Helper()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", contentType)
	storage, err := common.CreateBodyStorage([]byte(body))
	require.NoError(t, err)
	c.Set(common.KeyBodyStorage, storage)
	t.Cleanup(func() { storage.Close() })
	return c
}

func TestSyncImageRequestsB64JSON(t *testing.T) {
	c := imageSyncQueueTestContext(t, `{"model":"gpt-image-2","response_format":"b64_json"}`, "application/json")
	require.Equal(t, "b64_json", requestedSyncImageResponseFormat(c))
	require.True(t, shouldRunSyncImageViaQueue(c))
}

func TestSyncImageDefaultResponseUsesQueueAndReturnsURL(t *testing.T) {
	c := imageSyncQueueTestContext(t, `{"model":"gpt-image-2"}`, "application/json")
	t.Setenv("IMAGE_SYNC_QUEUE_DEFAULT_RESPONSE_IS_URL", "true")
	require.Equal(t, "url", requestedSyncImageResponseFormat(c))
	require.True(t, shouldRunSyncImageViaQueue(c))
}

func TestSyncImageURLResponseCanUseQueue(t *testing.T) {
	c := imageSyncQueueTestContext(t, `{"model":"gpt-image-2","response_format":"url"}`, "application/json")
	require.Equal(t, "url", requestedSyncImageResponseFormat(c))
	require.True(t, shouldRunSyncImageViaQueue(c))
}

func TestSyncImageMultipartB64ResponseUsesQueueAndKeepsExplicitFormat(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("model", "gpt-image-2"))
	require.NoError(t, writer.WriteField("prompt", "test"))
	require.NoError(t, writer.WriteField("response_format", "b64_json"))
	require.NoError(t, writer.Close())

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body.Bytes()))
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	storage, err := common.CreateBodyStorage(body.Bytes())
	require.NoError(t, err)
	c.Set(common.KeyBodyStorage, storage)
	t.Cleanup(func() {
		if c.Request.MultipartForm != nil {
			_ = c.Request.MultipartForm.RemoveAll()
		}
		_ = storage.Close()
	})

	require.Equal(t, "b64_json", requestedSyncImageResponseFormat(c))
	require.True(t, shouldRunSyncImageViaQueue(c))
}

func TestQueuedSyncB64ResponseStreamsStoredResult(t *testing.T) {
	previous := downloadTaskImageResult
	downloadTaskImageResult = func(_ context.Context, _, _ string) (*os.File, error) {
		file, err := os.CreateTemp("", "queued-sync-b64-test-*")
		if err != nil {
			return nil, err
		}
		if _, err = file.Write([]byte("image-bytes")); err != nil {
			return nil, err
		}
		_, err = file.Seek(0, 0)
		return file, err
	}
	t.Cleanup(func() { downloadTaskImageResult = previous })

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	task := &model.Task{
		Properties: model.Properties{OriginModelName: "gpt-image-2"},
		PrivateData: model.TaskPrivateData{
			ImageResultURLs: []string{"https://tmp.cangyuansuanli.cn/gen-images/example.png"},
		},
	}

	writeQueuedSyncImageResponse(c, task, "b64_json")
	require.Equal(t, http.StatusOK, recorder.Code)
	var response dto.ImageResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.Len(t, response.Data, 1)
	require.Equal(t, "aW1hZ2UtYnl0ZXM=", response.Data[0].B64Json)
	require.Empty(t, response.Data[0].Url)
}
