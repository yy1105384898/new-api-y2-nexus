package controller

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

func shouldRunSyncImageViaQueue(c *gin.Context) bool {
	return c != nil && c.Request != nil && common.GetEnvOrDefaultBool("IMAGE_SYNC_VIA_QUEUE", true)
}

// requestedSyncImageResponseFormat freezes the public response contract before
// the request enters the durable queue. The public default is url; b64_json is
// retained only when the client requests it explicitly.
func requestedSyncImageResponseFormat(c *gin.Context) string {
	format := syncImageResponseFormat(c)
	if strings.EqualFold(format, "b64_json") {
		return "b64_json"
	}
	if format == "" && !common.GetEnvOrDefaultBool("IMAGE_SYNC_QUEUE_DEFAULT_RESPONSE_IS_URL", true) {
		return "b64_json"
	}
	return "url"
}

func syncImageResponseFormat(c *gin.Context) string {
	if c != nil && c.Request != nil && c.Request.MultipartForm == nil && strings.Contains(strings.ToLower(c.Request.Header.Get("Content-Type")), "multipart/form-data") {
		if form, err := common.ParseMultipartFormReusable(c); err == nil {
			c.Request.MultipartForm = form
			c.Request.PostForm = form.Value
		}
	}
	if c.Request.MultipartForm != nil {
		values := c.Request.MultipartForm.Value["response_format"]
		if len(values) > 0 {
			return strings.TrimSpace(values[0])
		}
		return ""
	}
	if c.Request.PostForm != nil {
		if value := strings.TrimSpace(c.Request.PostForm.Get("response_format")); value != "" {
			return value
		}
	}
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		return ""
	}
	body, err := storage.Bytes()
	if err != nil || len(body) == 0 {
		return ""
	}
	var request struct {
		ResponseFormat string `json:"response_format"`
	}
	if err := common.Unmarshal(body, &request); err != nil {
		return ""
	}
	return strings.TrimSpace(request.ResponseFormat)
}

func relaySyncImageViaQueue(c *gin.Context) {
	responseFormat := requestedSyncImageResponseFormat(c)
	recorder := httptest.NewRecorder()
	submit, _ := gin.CreateTestContext(recorder)
	submit.Request = c.Request
	submit.Params = c.Params
	submit.Keys = make(map[string]any, len(c.Keys))
	for key, value := range c.Keys {
		submit.Keys[key] = value
	}
	submit.Set("image_sync_wait", true)

	RelayImageTaskSubmit(submit)
	if recorder.Code != http.StatusOK {
		c.Data(recorder.Code, recorder.Header().Get("Content-Type"), recorder.Body.Bytes())
		return
	}

	var job dto.OpenAIImageJob
	if err := common.Unmarshal(recorder.Body.Bytes(), &job); err != nil || job.ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": "failed to create internal image task",
			"type":    "server_error",
		}})
		return
	}
	c.Header("X-Cangyuan-Task-Id", job.ID)
	waitForQueuedSyncImage(c, job.ID, responseFormat)
}

func waitForQueuedSyncImage(c *gin.Context, taskID, responseFormat string) {
	timeout := time.Duration(common.GetEnvOrDefault("IMAGE_SYNC_QUEUE_WAIT_SECONDS", 300)) * time.Second
	interval := time.Duration(common.GetEnvOrDefault("IMAGE_SYNC_QUEUE_POLL_INTERVAL_MS", 250)) * time.Millisecond
	if interval < 50*time.Millisecond {
		interval = 50 * time.Millisecond
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	userID := c.GetInt("id")
	for {
		task, exists, err := model.GetByTaskIdForFetch(userID, taskID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"message": "failed to read image task status",
				"type":    "server_error",
			}})
			return
		}
		if exists && task != nil {
			switch task.Status {
			case model.TaskStatusSuccess:
				writeQueuedSyncImageResponse(c, task, responseFormat)
				return
			case model.TaskStatusFailure:
				message := strings.TrimSpace(task.FailReason)
				if message == "" {
					message = "image generation failed"
				}
				c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
					"message": message,
					"type":    "upstream_error",
				}})
				return
			}
		}

		select {
		case <-c.Request.Context().Done():
			return
		case <-timer.C:
			c.Header("Retry-After", "5")
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": gin.H{
				"message": "image task is still running; retry the request later",
				"type":    "timeout_error",
			}})
			return
		case <-ticker.C:
		}
	}
}

func queuedSyncImageResultURLs(task *model.Task) []string {
	if task == nil {
		return nil
	}
	urls := task.PrivateData.ImageResultURLs
	if len(urls) == 0 && strings.TrimSpace(task.PrivateData.ResultURL) != "" {
		urls = []string{task.PrivateData.ResultURL}
	}
	return urls
}

func writeQueuedSyncImageResponse(c *gin.Context, task *model.Task, responseFormat string) {
	urls := queuedSyncImageResultURLs(task)
	if strings.EqualFold(responseFormat, "b64_json") {
		writeQueuedSyncImageB64Response(c, task, urls)
		return
	}
	data := make([]dto.ImageData, 0, len(urls))
	for _, resultURL := range urls {
		if strings.TrimSpace(resultURL) != "" {
			data = append(data, dto.ImageData{Url: resultURL})
		}
	}
	if len(data) == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"message": "image task completed without a result",
			"type":    "upstream_error",
		}})
		return
	}
	c.JSON(http.StatusOK, dto.ImageResponse{Created: time.Now().Unix(), Data: data})
}

var (
	b64DeliveryLimiterOnce  sync.Once
	b64DeliveryLimiter      chan struct{}
	downloadTaskImageResult = service.DownloadTaskImageResult
)

func acquireB64DeliverySlot(ctx *gin.Context) (func(), error) {
	b64DeliveryLimiterOnce.Do(func() {
		limit := common.GetEnvOrDefault("IMAGE_B64_DELIVERY_MAX_CONCURRENT", 8)
		if limit < 1 {
			limit = 1
		}
		b64DeliveryLimiter = make(chan struct{}, limit)
	})
	select {
	case b64DeliveryLimiter <- struct{}{}:
		return func() { <-b64DeliveryLimiter }, nil
	case <-ctx.Request.Context().Done():
		return nil, ctx.Request.Context().Err()
	}
}

func writeQueuedSyncImageB64Response(c *gin.Context, task *model.Task, urls []string) {
	if len(urls) == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"message": "image task completed without a result",
			"type":    "upstream_error",
		}})
		return
	}
	release, err := acquireB64DeliverySlot(c)
	if err != nil {
		return
	}
	defer release()

	files := make([]*os.File, 0, len(urls))
	defer func() {
		for _, file := range files {
			name := file.Name()
			_ = file.Close()
			_ = os.Remove(name)
		}
	}()
	for _, resultURL := range urls {
		file, downloadErr := downloadTaskImageResult(c.Request.Context(), task.Properties.OriginModelName, resultURL)
		if downloadErr != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
				"message": "failed to prepare queued image response",
				"type":    "upstream_error",
			}})
			return
		}
		files = append(files, file)
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Status(http.StatusOK)
	_, _ = io.WriteString(c.Writer, `{"created":`+strconv.FormatInt(time.Now().Unix(), 10)+`,"data":[`)
	for index, file := range files {
		if index > 0 {
			_, _ = io.WriteString(c.Writer, ",")
		}
		_, _ = io.WriteString(c.Writer, `{"b64_json":"`)
		encoder := base64.NewEncoder(base64.StdEncoding, c.Writer)
		if _, err = io.Copy(encoder, file); err != nil {
			_ = encoder.Close()
			_ = c.Error(fmt.Errorf("stream queued image result: %w", err))
			return
		}
		if err = encoder.Close(); err != nil {
			_ = c.Error(fmt.Errorf("finish queued image base64 response: %w", err))
			return
		}
		_, _ = io.WriteString(c.Writer, `"}`)
	}
	_, _ = io.WriteString(c.Writer, `]}`)
}
