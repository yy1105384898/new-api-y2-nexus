package relay

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay/channel/task/taskcommon"
	"github.com/QuantumNous/new-api/service"
)

var imageAsyncQueue chan string

func StartImageAsyncWorker() {
	maxConcurrent := 32
	if v := strings.TrimSpace(os.Getenv("IMAGE_ASYNC_MAX_CONCURRENT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxConcurrent = n
		}
	}
	imageAsyncQueue = make(chan string, maxConcurrent*2)
	for i := 0; i < maxConcurrent; i++ {
		go imageAsyncWorkerLoop()
	}
	common.SysLog(fmt.Sprintf("image async worker started, concurrency=%d", maxConcurrent))
	go recoverPendingImageAsyncTasks()
}

func recoverPendingImageAsyncTasks() {
	tasks := model.GetPendingImageAsyncTasks(64)
	if len(tasks) == 0 {
		return
	}
	common.SysLog(fmt.Sprintf("image async recovering %d pending tasks", len(tasks)))
	for _, task := range tasks {
		if task != nil && task.TaskID != "" {
			EnqueueImageAsyncTask(task.TaskID)
		}
	}
}

func EnqueueImageAsyncTask(taskID string) {
	if imageAsyncQueue == nil {
		StartImageAsyncWorker()
	}
	select {
	case imageAsyncQueue <- taskID:
	default:
		common.SysLog("image async queue full, running inline for task " + taskID)
		go processImageAsyncTask(taskID)
	}
}

func imageAsyncWorkerLoop() {
	for taskID := range imageAsyncQueue {
		processImageAsyncTask(taskID)
	}
}

func processImageAsyncTask(taskID string) {
	ctx := context.Background()
	task, exist, err := model.GetByOnlyTaskId(taskID)
	if err != nil || !exist || task == nil {
		return
	}
	if task.Properties.TaskKind != constant.TaskKindImage {
		return
	}
	if task.Status == model.TaskStatusSuccess || task.Status == model.TaskStatusFailure {
		return
	}

	snap := task.Snapshot()
	task.Status = model.TaskStatusInProgress
	task.Progress = taskcommon.ProgressInProgress
	if task.StartTime == 0 {
		task.StartTime = time.Now().Unix()
	}
	if _, err := task.UpdateWithStatus(snap.Status); err != nil {
		logger.LogError(ctx, "image async update in_progress failed: "+err.Error())
		return
	}

	images, _, execErr := ExecuteImageTaskUpstream(task)
	if execErr != nil {
		failImageAsyncTask(ctx, task, execErr.Error())
		return
	}

	publicURLs, uploadErr := uploadTaskImagesToR2(ctx, task, images)
	if uploadErr != nil {
		failImageAsyncTask(ctx, task, uploadErr.Error())
		return
	}

	meta := map[string]any{"r2_urls": publicURLs}
	task.SetData(meta)
	task.PrivateData.ImageResultURLs = publicURLs
	if len(publicURLs) > 0 {
		task.PrivateData.ResultURL = publicURLs[0]
	}
	task.Status = model.TaskStatusSuccess
	task.Progress = taskcommon.ProgressComplete
	task.FinishTime = time.Now().Unix()
	if _, err := task.UpdateWithStatus(model.TaskStatusInProgress); err != nil {
		logger.LogError(ctx, "image async mark success failed: "+err.Error())
		return
	}

	service.RecalculateTaskQuota(ctx, task, task.Quota, "image async complete")
}

func uploadTaskImagesToR2(ctx context.Context, task *model.Task, images []dto.ImageData) ([]string, error) {
	channelBaseURL := ""
	if channel, err := model.GetChannelById(task.ChannelId, true); err == nil && channel != nil {
		channelBaseURL = channel.GetBaseURL()
	}
	publicURLs := make([]string, 0, len(images))
	for index, item := range images {
		data, remoteURL, err := DecodeImageDataItemExported(item)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			uploaded, err := service.UploadGeneratedImageBytes(ctx, task.UserId, task.TaskID, index, data, "image/png")
			if err != nil {
				return nil, err
			}
			publicURLs = append(publicURLs, uploaded.PublicURL)
			continue
		}
		if remoteURL != "" {
			remoteURL = rewriteLoopbackUpstreamImageURL(channelBaseURL, remoteURL)
			uploaded, err := uploadGeneratedImageFromURLWithRetry(ctx, task.UserId, task.TaskID, index, remoteURL)
			if err != nil {
				return nil, err
			}
			publicURLs = append(publicURLs, uploaded.PublicURL)
		}
	}
	if len(publicURLs) == 0 {
		return nil, fmt.Errorf("no images uploaded to R2")
	}
	return publicURLs, nil
}

func failImageAsyncTask(ctx context.Context, task *model.Task, reason string) {
	snap := task.Snapshot()
	task.Status = model.TaskStatusFailure
	task.Progress = taskcommon.ProgressComplete
	task.FailReason = reason
	task.FinishTime = time.Now().Unix()
	if _, err := task.UpdateWithStatus(snap.Status); err != nil {
		common.SysLog("image async mark failure failed: " + err.Error())
	}
	service.RefundTaskQuota(ctx, task, reason)
}

func uploadGeneratedImageFromURLWithRetry(ctx context.Context, userID int, taskID string, index int, imageURL string) (*service.R2UploadResult, error) {
	const maxAttempts = 3
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		uploaded, err := service.UploadGeneratedImageFromURL(ctx, userID, taskID, index, imageURL)
		if err == nil {
			return uploaded, nil
		}
		lastErr = err
		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
	}
	return nil, lastErr
}
