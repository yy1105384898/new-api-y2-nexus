package chatvideo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/task/oaivideo/vendors/defaultvideo"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

// TaskAdaptor converts NewAPI's public video-task contract to legacy video
// providers that expose generation through chat/completions. The frontend must
// never select or parse this upstream protocol directly.
type TaskAdaptor struct {
	defaultvideo.TaskAdaptor
}

func (a *TaskAdaptor) GetModelList() []string { return nil }

func (a *TaskAdaptor) GetChannelName() string { return "chat-video" }

func (a *TaskAdaptor) BuildRequestURL(info *relaycommon.RelayInfo) (string, error) {
	if info == nil || strings.TrimSpace(info.ChannelBaseUrl) == "" {
		return "", fmt.Errorf("chat video base url is empty")
	}
	return strings.TrimRight(info.ChannelBaseUrl, "/") + "/v1/chat/completions", nil
}

func (a *TaskAdaptor) BuildRequestHeader(_ *gin.Context, req *http.Request, info *relaycommon.RelayInfo) error {
	req.Header.Set("Authorization", "Bearer "+info.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream, application/json")
	return nil
}

func (a *TaskAdaptor) BuildRequestBody(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error) {
	req, err := relaycommon.GetTaskRequest(c)
	if err != nil {
		return nil, err
	}
	modelName := strings.TrimSpace(info.UpstreamModelName)
	if modelName == "" {
		modelName = strings.TrimSpace(req.Model)
	}
	if modelName == "" || strings.TrimSpace(req.Prompt) == "" {
		return nil, fmt.Errorf("model and prompt are required")
	}

	var content any = strings.TrimSpace(req.Prompt)
	if len(req.Images) > 0 {
		parts := []map[string]any{{"type": "text", "text": strings.TrimSpace(req.Prompt)}}
		for _, imageURL := range req.Images {
			if imageURL = strings.TrimSpace(imageURL); imageURL != "" {
				parts = append(parts, map[string]any{
					"type":      "image_url",
					"image_url": map[string]any{"url": imageURL},
				})
			}
		}
		content = parts
	}

	bodyMap := map[string]any{
		"model":  modelName,
		"stream": true,
		"messages": []map[string]any{
			{"role": "user", "content": content},
		},
	}
	if duration := req.RequestedDurationSeconds(); duration > 0 {
		bodyMap["duration"] = duration
	}
	if req.AspectRatio != "" {
		bodyMap["aspect_ratio"] = req.AspectRatio
	}
	if req.Resolution != "" {
		bodyMap["resolution"] = req.Resolution
	}
	if req.GenerateAudio != nil {
		bodyMap["generate_audio"] = *req.GenerateAudio
	}
	body, err := common.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	c.Request.Header.Set("Content-Type", "application/json")
	return bytes.NewReader(body), nil
}

func (a *TaskAdaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, error) {
	return channel.DoTaskApiRequest(a, c, info, requestBody)
}

func (a *TaskAdaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (string, []byte, *dto.TaskError) {
	if resp == nil || resp.Body == nil {
		return "", nil, service.TaskErrorWrapper(fmt.Errorf("empty chat video response"), "invalid_response", http.StatusBadGateway)
	}
	defer resp.Body.Close()
	content, err := readChatContent(resp)
	if err != nil {
		return "", nil, service.TaskErrorWrapper(err, "invalid_response", http.StatusBadGateway)
	}
	videoURL := extractVideoURL(content)
	if videoURL == "" {
		return "", nil, service.TaskErrorWrapper(fmt.Errorf("chat video response does not contain a video url"), "invalid_response", http.StatusBadGateway)
	}

	data := map[string]any{
		"id":         info.PublicTaskID,
		"task_id":    info.PublicTaskID,
		"model":      info.UpstreamModelName,
		"status":     "completed",
		"progress":   100,
		"url":        videoURL,
		"result_url": videoURL,
		"video_url":  videoURL,
		"data":       []map[string]any{{"url": videoURL, "video_url": videoURL}},
	}
	taskData, err := common.Marshal(data)
	if err != nil {
		return "", nil, service.TaskErrorWrapper(err, "marshal_response_failed", http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, data)
	return info.PublicTaskID, taskData, nil
}

type chatChunk struct {
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
	Message string `json:"message"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func readChatContent(resp *http.Response) (string, error) {
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "text/event-stream") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return contentFromChunk(body)
	}

	var content strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}
		part, err := contentFromChunk([]byte(data))
		if err != nil {
			return "", err
		}
		content.WriteString(part)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content.String(), nil
}

func contentFromChunk(data []byte) (string, error) {
	var chunk chatChunk
	if err := common.Unmarshal(data, &chunk); err != nil {
		return "", err
	}
	if chunk.Error != nil && strings.TrimSpace(chunk.Error.Message) != "" {
		return "", fmt.Errorf("%s", strings.TrimSpace(chunk.Error.Message))
	}
	if len(chunk.Choices) == 0 {
		if strings.TrimSpace(chunk.Message) != "" {
			return "", fmt.Errorf("%s", strings.TrimSpace(chunk.Message))
		}
		return "", nil
	}
	return chunk.Choices[0].Delta.Content + chunk.Choices[0].Message.Content, nil
}

var (
	markdownVideoURL = regexp.MustCompile(`\((https?://[^)\s]+)\)`)
	plainVideoURL    = regexp.MustCompile(`https?://[^\s<>"']+`)
)

func extractVideoURL(content string) string {
	if match := markdownVideoURL.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return strings.TrimRight(plainVideoURL.FindString(content), ").,;]")
}
