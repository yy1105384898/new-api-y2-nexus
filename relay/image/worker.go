package image

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay/channel/task/taskcommon"
	"github.com/QuantumNous/new-api/service"
	"github.com/go-redis/redis/v8"
)

const (
	imageTaskNotifyQueue      = "new-api:image:task-notify"
	imageTaskNotifyDedup      = "new-api:image:task-notify:"
	adobeTaskNotifyQueue      = "new-api:image:task-notify:adobe"
	adobeTaskNotifyDedup      = "new-api:image:task-notify:adobe:"
	imageTaskDoneChannel      = "new-api:image:task-done:"
	defaultAdobeChannelIDList = "75"
)

type imageWorkerConfig struct {
	concurrency    int
	queueCapacity  int
	dispatchBatch  int
	dbScanInterval time.Duration
	leaseDuration  time.Duration
	maxAttempts    int
}

type imageTaskDispatcher struct {
	once            sync.Once
	name            string
	notifyKey       string
	dedupKey        string
	channelIDs      []int
	includeChannels bool
	queue           chan string
	redis           *redis.Client
	owner           string
	config          imageWorkerConfig
	mu              sync.Mutex
	queued          map[string]struct{}
	enabled         bool
	active          atomic.Int64
	completed       atomic.Int64
	failed          atomic.Int64
}

type WorkerLaneStats struct {
	Concurrency   int   `json:"concurrency"`
	QueueCapacity int   `json:"queue_capacity"`
	QueueBuffered int   `json:"queue_buffered"`
	Active        int64 `json:"active"`
	Completed     int64 `json:"completed"`
	Failed        int64 `json:"failed"`
	Backlog       int64 `json:"backlog"`
	RedisPending  int64 `json:"redis_pending"`
	DBScanMS      int64 `json:"db_scan_ms"`
}

type WorkerStats struct {
	Enabled       bool                       `json:"enabled"`
	Owner         string                     `json:"owner"`
	Concurrency   int                        `json:"concurrency"`
	QueueCapacity int                        `json:"queue_capacity"`
	QueueBuffered int                        `json:"queue_buffered"`
	Active        int64                      `json:"active"`
	Completed     int64                      `json:"completed"`
	Failed        int64                      `json:"failed"`
	GlobalBacklog int64                      `json:"global_backlog"`
	RedisPending  int64                      `json:"redis_pending"`
	DBScanMS      int64                      `json:"db_scan_ms"`
	Lanes         map[string]WorkerLaneStats `json:"lanes"`
}

func GetWorkerStats() (WorkerStats, error) {
	adobeChannelIDs := AdobeDirectChannelIDs()
	defaultBacklog, _, err := model.CountActiveImageTasksForChannels(0, adobeChannelIDs, false)
	if err != nil {
		return WorkerStats{}, err
	}
	adobeBacklog, _, err := model.CountActiveImageTasksForChannels(0, adobeChannelIDs, true)
	if err != nil {
		return WorkerStats{}, err
	}
	defaultStats := workerLaneStats(&imageDispatcher, defaultBacklog)
	adobeStats := workerLaneStats(&adobeImageDispatcher, adobeBacklog)
	stats := WorkerStats{
		Enabled:       imageDispatcher.enabled || adobeImageDispatcher.enabled,
		Owner:         imageDispatcher.owner,
		Concurrency:   defaultStats.Concurrency + adobeStats.Concurrency,
		QueueCapacity: defaultStats.QueueCapacity + adobeStats.QueueCapacity,
		QueueBuffered: defaultStats.QueueBuffered + adobeStats.QueueBuffered,
		Active:        defaultStats.Active + adobeStats.Active,
		Completed:     defaultStats.Completed + adobeStats.Completed,
		Failed:        defaultStats.Failed + adobeStats.Failed,
		GlobalBacklog: defaultBacklog + adobeBacklog,
		RedisPending:  defaultStats.RedisPending + adobeStats.RedisPending,
		DBScanMS:      imageDispatcher.config.dbScanInterval.Milliseconds(),
		Lanes: map[string]WorkerLaneStats{
			"default": defaultStats,
			"adobe":   adobeStats,
		},
	}
	return stats, nil
}

var imageDispatcher imageTaskDispatcher
var adobeImageDispatcher imageTaskDispatcher
var imageTaskURLPattern = regexp.MustCompile(`https?://[^\s"']+`)

var imageTaskDoneNotifier struct {
	once    sync.Once
	mu      sync.Mutex
	ready   bool
	waiters map[string]map[chan struct{}]struct{}
}

func imageWorkerEnvInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func loadImageWorkerConfig() imageWorkerConfig {
	concurrency := imageWorkerEnvInt("IMAGE_ASYNC_MAX_CONCURRENT", 32)
	dbScanFallback := 1000
	if common.RedisEnabled && common.RDB != nil {
		dbScanFallback = 15000
	}
	return imageWorkerConfig{
		concurrency:    concurrency,
		queueCapacity:  imageWorkerEnvInt("IMAGE_ASYNC_QUEUE_CAPACITY", concurrency*4),
		dispatchBatch:  imageWorkerEnvInt("IMAGE_ASYNC_DISPATCH_BATCH", concurrency*2),
		dbScanInterval: time.Duration(imageWorkerEnvInt("IMAGE_ASYNC_DB_SCAN_INTERVAL_MS", dbScanFallback)) * time.Millisecond,
		leaseDuration:  time.Duration(imageWorkerEnvInt("IMAGE_ASYNC_LEASE_SECONDS", 180)) * time.Second,
		maxAttempts:    imageWorkerEnvInt("IMAGE_ASYNC_MAX_ATTEMPTS", 3),
	}
}

func loadAdobeImageWorkerConfig(base imageWorkerConfig) imageWorkerConfig {
	concurrency := imageWorkerEnvInt("IMAGE_ASYNC_ADOBE_MAX_CONCURRENT", 14)
	return imageWorkerConfig{
		concurrency:    concurrency,
		queueCapacity:  imageWorkerEnvInt("IMAGE_ASYNC_ADOBE_QUEUE_CAPACITY", concurrency*4),
		dispatchBatch:  imageWorkerEnvInt("IMAGE_ASYNC_ADOBE_DISPATCH_BATCH", concurrency*2),
		dbScanInterval: time.Duration(imageWorkerEnvInt("IMAGE_ASYNC_ADOBE_DB_SCAN_INTERVAL_MS", int(base.dbScanInterval/time.Millisecond))) * time.Millisecond,
		leaseDuration:  base.leaseDuration,
		maxAttempts:    base.maxAttempts,
	}
}

// AdobeDirectChannelIDs defines the deployment-specific channel lane without
// baking one database id into the scheduler. Channel 75 remains the compatible
// default for existing installations.
func AdobeDirectChannelIDs() []int {
	raw := strings.TrimSpace(os.Getenv("IMAGE_ASYNC_ADOBE_CHANNEL_IDS"))
	if raw == "" {
		raw = defaultAdobeChannelIDList
	}
	seen := make(map[int]struct{})
	ids := make([]int, 0)
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ';' || r == ' ' }) {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || id <= 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func IsAdobeDirectChannel(channelID int) bool {
	for _, candidate := range AdobeDirectChannelIDs() {
		if channelID == candidate {
			return true
		}
	}
	return false
}

func imageWorkerOwner() string {
	hostname, _ := os.Hostname()
	parts := []string{strings.TrimSpace(common.NodeName), strings.TrimSpace(hostname), strconv.Itoa(os.Getpid())}
	nonEmpty := parts[:0]
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, "/")
}

func workerLaneStats(dispatcher *imageTaskDispatcher, backlog int64) WorkerLaneStats {
	stats := WorkerLaneStats{
		Concurrency:   dispatcher.config.concurrency,
		QueueCapacity: dispatcher.config.queueCapacity,
		Active:        dispatcher.active.Load(),
		Completed:     dispatcher.completed.Load(),
		Failed:        dispatcher.failed.Load(),
		Backlog:       backlog,
		DBScanMS:      dispatcher.config.dbScanInterval.Milliseconds(),
	}
	if dispatcher.queue != nil {
		stats.QueueBuffered = len(dispatcher.queue)
	}
	if common.RedisEnabled && common.RDB != nil && dispatcher.notifyKey != "" {
		stats.RedisPending, _ = common.RDB.LLen(context.Background(), dispatcher.notifyKey).Result()
	}
	return stats
}

// StartWorker starts a strictly bounded local worker pool. PostgreSQL remains
// the durable queue; every worker node continuously discovers claimable jobs.
func StartWorker() {
	imageDispatcher.once.Do(func() {
		if strings.EqualFold(strings.TrimSpace(os.Getenv("IMAGE_ASYNC_WORKER_ENABLED")), "false") {
			common.SysLog("image async worker disabled on this node")
			return
		}
		defaultConfig := loadImageWorkerConfig()
		adobeConfig := loadAdobeImageWorkerConfig(defaultConfig)
		owner := imageWorkerOwner()
		adobeChannelIDs := AdobeDirectChannelIDs()
		startImageTaskDispatcher(
			&imageDispatcher, "default", imageTaskNotifyQueue, imageTaskNotifyDedup,
			owner, defaultConfig, adobeChannelIDs, false,
		)
		startImageTaskDispatcher(
			&adobeImageDispatcher, "adobe", adobeTaskNotifyQueue, adobeTaskNotifyDedup,
			owner, adobeConfig, adobeChannelIDs, true,
		)
	})
}

func startImageTaskDispatcher(
	dispatcher *imageTaskDispatcher,
	name string,
	notifyKey string,
	dedupKey string,
	owner string,
	config imageWorkerConfig,
	channelIDs []int,
	includeChannels bool,
) {
	dispatcher.name = name
	dispatcher.notifyKey = notifyKey
	dispatcher.dedupKey = dedupKey
	dispatcher.owner = owner
	dispatcher.config = config
	dispatcher.channelIDs = append([]int(nil), channelIDs...)
	dispatcher.includeChannels = includeChannels
	dispatcher.queue = make(chan string, config.queueCapacity)
	dispatcher.queued = make(map[string]struct{}, config.queueCapacity)
	if common.RedisEnabled && common.RDB != nil {
		options := *common.RDB.Options()
		if options.PoolSize < config.concurrency+2 {
			options.PoolSize = config.concurrency + 2
		}
		dispatcher.redis = redis.NewClient(&options)
	}
	dispatcher.enabled = true
	for i := 0; i < config.concurrency; i++ {
		go imageAsyncWorkerLoop(dispatcher)
	}
	go imageAsyncDispatchLoop(dispatcher)
	common.SysLog(fmt.Sprintf(
		"image async worker lane started, lane=%s owner=%s concurrency=%d queue_capacity=%d db_scan=%s lease=%s channels=%v include=%t",
		dispatcher.name, dispatcher.owner, config.concurrency, config.queueCapacity,
		config.dbScanInterval, config.leaseDuration, dispatcher.channelIDs, dispatcher.includeChannels,
	))
}

// EnqueueTask is only a wake-up hint. If the bounded local buffer is full, the
// task stays QUEUED in PostgreSQL and a dispatcher picks it up later.
func EnqueueTask(taskID string) bool {
	return enqueueImageTask(&imageDispatcher, taskID)
}

func EnqueueTaskForChannel(taskID string, channelID int) bool {
	if IsAdobeDirectChannel(channelID) {
		return enqueueImageTask(&adobeImageDispatcher, taskID)
	}
	return enqueueImageTask(&imageDispatcher, taskID)
}

func enqueueImageTask(dispatcher *imageTaskDispatcher, taskID string) bool {
	if taskID == "" {
		return false
	}
	if common.RedisEnabled && common.RDB != nil {
		notifyKey, dedupPrefix := imageDispatcherQueueKeys(dispatcher)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		dedupKey := dedupPrefix + taskID
		won, err := common.RDB.SetNX(ctx, dedupKey, "1", 30*time.Second).Result()
		if err == nil && !won {
			return true
		}
		if err == nil {
			if err = common.RDB.RPush(ctx, notifyKey, taskID).Err(); err == nil {
				return true
			}
			_ = common.RDB.Del(ctx, dedupKey).Err()
		}
	}
	return enqueueLocalImageTask(dispatcher, taskID)
}

func imageDispatcherQueueKeys(dispatcher *imageTaskDispatcher) (string, string) {
	if dispatcher != nil && dispatcher.notifyKey != "" && dispatcher.dedupKey != "" {
		return dispatcher.notifyKey, dispatcher.dedupKey
	}
	if dispatcher == &adobeImageDispatcher {
		return adobeTaskNotifyQueue, adobeTaskNotifyDedup
	}
	return imageTaskNotifyQueue, imageTaskNotifyDedup
}

func enqueueLocalImageTask(dispatcher *imageTaskDispatcher, taskID string) bool {
	if !dispatcher.enabled || dispatcher.queue == nil {
		return false
	}
	dispatcher.mu.Lock()
	if _, exists := dispatcher.queued[taskID]; exists {
		dispatcher.mu.Unlock()
		return true
	}
	dispatcher.queued[taskID] = struct{}{}
	dispatcher.mu.Unlock()

	select {
	case dispatcher.queue <- taskID:
		return true
	default:
		dispatcher.mu.Lock()
		delete(dispatcher.queued, taskID)
		dispatcher.mu.Unlock()
		return false
	}
}

func imageAsyncDispatchLoop(dispatcher *imageTaskDispatcher) {
	ticker := time.NewTicker(dispatcher.config.dbScanInterval)
	defer ticker.Stop()
	for {
		dispatchClaimableImageTasks(dispatcher)
		<-ticker.C
	}
}

func dispatchClaimableImageTasks(dispatcher *imageTaskDispatcher) {
	ids := model.GetClaimableImageAsyncTaskIDsForChannels(
		dispatcher.config.dispatchBatch, time.Now().Unix(),
		dispatcher.channelIDs, dispatcher.includeChannels,
	)
	for _, taskID := range ids {
		if !enqueueImageTask(dispatcher, taskID) {
			return
		}
	}
}

func imageAsyncWorkerLoop(dispatcher *imageTaskDispatcher) {
	for {
		taskID, ok := nextImageAsyncTaskID(dispatcher)
		if !ok {
			return
		}
		processImageAsyncTask(dispatcher, taskID)
		dispatcher.mu.Lock()
		delete(dispatcher.queued, taskID)
		dispatcher.mu.Unlock()
	}
}

// nextImageAsyncTaskID lets each idle worker compete for Redis notifications.
// Distribution therefore follows free execution slots instead of assigning an
// equal share to every node regardless of its configured concurrency.
func nextImageAsyncTaskID(dispatcher *imageTaskDispatcher) (string, bool) {
	for dispatcher.redis != nil {
		select {
		case taskID, ok := <-dispatcher.queue:
			return taskID, ok
		default:
		}
		result, err := dispatcher.redis.BLPop(context.Background(), 2*time.Second, dispatcher.notifyKey).Result()
		if err == nil && len(result) == 2 && result[1] != "" {
			return result[1], true
		}
		if err != nil && err != redis.Nil {
			common.SysError("image redis worker: " + err.Error())
			time.Sleep(time.Second)
		}
	}
	taskID, ok := <-dispatcher.queue
	return taskID, ok
}

func processImageAsyncTask(dispatcher *imageTaskDispatcher, taskID string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	task, claimed, err := model.ClaimImageAsyncTaskForChannels(
		taskID, dispatcher.owner, dispatcher.config.leaseDuration,
		dispatcher.channelIDs, dispatcher.includeChannels,
	)
	if err != nil {
		common.SysError(fmt.Sprintf("image async claim failed for %s: %v", taskID, err))
		return
	}
	if !claimed || task == nil {
		return
	}
	if common.RedisEnabled && common.RDB != nil {
		_, dedupPrefix := imageDispatcherQueueKeys(dispatcher)
		_ = common.RDB.Del(context.Background(), dedupPrefix+taskID).Err()
	}
	dispatcher.active.Add(1)
	defer dispatcher.active.Add(-1)
	if task.Attempt > dispatcher.config.maxAttempts {
		failImageAsyncTask(dispatcher, ctx, task, model.TaskStatusInProgress, "image task exceeded maximum attempts")
		return
	}

	heartbeatDone := make(chan struct{})
	go imageAsyncLeaseHeartbeat(dispatcher, task.TaskID, heartbeatDone, cancel)
	defer close(heartbeatDone)

	images, _, execErr := executeTaskUpstream(ctx, task)
	if execErr != nil {
		failImageAsyncTask(dispatcher, ctx, task, model.TaskStatusInProgress, execErr.Error())
		return
	}

	resultURLs, resolveErr := resolveTaskImageResultURLs(ctx, task, images)
	if resolveErr != nil {
		failImageAsyncTask(dispatcher, ctx, task, model.TaskStatusInProgress, resolveErr.Error())
		return
	}

	meta := map[string]any{"result_urls": resultURLs}
	task.SetData(meta)
	task.PrivateData.ImageResultURLs = resultURLs
	if len(resultURLs) > 0 {
		task.PrivateData.ResultURL = resultURLs[0]
	}
	task.Status = model.TaskStatusSuccess
	task.Progress = taskcommon.ProgressComplete
	task.FinishTime = time.Now().Unix()
	task.LeaseOwner = ""
	task.LeaseExpiresAt = 0
	inputObjectKeys := imageTaskInputObjectKeys(task)
	task.ReleaseRequestSnapshot()
	won, err := model.UpdateImageTaskWithLease(task, dispatcher.owner)
	if err != nil {
		common.SysError(fmt.Sprintf("image task %s success lease CAS error: %v", task.TaskID, err))
		return
	}
	if !won {
		common.SysLog("image task success lease lost for " + task.TaskID)
		return
	}
	cleanupImageTaskInputs(task.TaskID, inputObjectKeys)
	dispatcher.completed.Add(1)
	publishImageTaskDone(task.TaskID)

	service.RecalculateTaskQuota(ctx, task, task.Quota, "image async complete")
}

func imageAsyncLeaseHeartbeat(dispatcher *imageTaskDispatcher, taskID string, done <-chan struct{}, cancel context.CancelFunc) {
	interval := dispatcher.config.leaseDuration / 3
	if interval < time.Second {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ok, err := model.RenewImageAsyncTaskLease(taskID, dispatcher.owner, dispatcher.config.leaseDuration)
			if err != nil {
				common.SysError(fmt.Sprintf("image async lease renewal failed for %s: %v", taskID, err))
				continue
			}
			if !ok {
				common.SysLog("image async lease lost for task " + taskID)
				cancel()
				return
			}
		}
	}
}

// resolveTaskImageResultURLs only publishes first-party R2 URLs. Upstream URLs
// remain process-local and are never stored in public task data.
func resolveTaskImageResultURLs(ctx context.Context, task *model.Task, images []dto.ImageData) ([]string, error) {
	return service.RehostTaskImageResultURLs(ctx, task.UserId, task.TaskID, taskUpstreamBaseURL(task), task.Properties.OriginModelName, images)
}

func failImageAsyncTask(dispatcher *imageTaskDispatcher, ctx context.Context, task *model.Task, fromStatus model.TaskStatus, reason string) {
	reason = imageTaskURLPattern.ReplaceAllString(reason, "[upstream-url-redacted]")
	task.Status = model.TaskStatusFailure
	task.Progress = taskcommon.ProgressComplete
	task.FailReason = reason
	task.FinishTime = time.Now().Unix()
	task.LeaseOwner = ""
	task.LeaseExpiresAt = 0
	inputObjectKeys := imageTaskInputObjectKeys(task)
	task.ReleaseRequestSnapshot()
	won, err := model.UpdateImageTaskWithLease(task, dispatcher.owner)
	if err != nil {
		common.SysError(fmt.Sprintf("image task %s failure lease CAS error: %v", task.TaskID, err))
		return
	}
	if !won {
		if reloaded, exist, err := model.GetByOnlyTaskId(task.TaskID); err == nil && exist {
			if reloaded.Status == model.TaskStatusSuccess {
				return
			}
		}
		return
	}
	cleanupImageTaskInputs(task.TaskID, inputObjectKeys)
	dispatcher.failed.Add(1)
	publishImageTaskDone(task.TaskID)
	service.RefundTaskQuota(ctx, task, reason)
}

func publishImageTaskDone(taskID string) {
	if taskID == "" || !common.RedisEnabled || common.RDB == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := common.RDB.Publish(ctx, imageTaskDoneChannel+taskID, "1").Err(); err != nil {
		common.SysError(fmt.Sprintf("image task completion notify failed for %s: %v", taskID, err))
	}
}

// SubscribeTaskDone multiplexes all task completion events through one Redis
// pattern subscription per API process, avoiding one Redis connection per
// synchronous waiter.
func SubscribeTaskDone(taskID string) (<-chan struct{}, func()) {
	if taskID == "" || !common.RedisEnabled || common.RDB == nil {
		return nil, func() {}
	}
	imageTaskDoneNotifier.once.Do(startImageTaskDoneNotifier)
	imageTaskDoneNotifier.mu.Lock()
	defer imageTaskDoneNotifier.mu.Unlock()
	if !imageTaskDoneNotifier.ready {
		return nil, func() {}
	}
	if imageTaskDoneNotifier.waiters == nil {
		imageTaskDoneNotifier.waiters = make(map[string]map[chan struct{}]struct{})
	}
	waiter := make(chan struct{}, 1)
	if imageTaskDoneNotifier.waiters[taskID] == nil {
		imageTaskDoneNotifier.waiters[taskID] = make(map[chan struct{}]struct{})
	}
	imageTaskDoneNotifier.waiters[taskID][waiter] = struct{}{}
	return waiter, func() {
		imageTaskDoneNotifier.mu.Lock()
		defer imageTaskDoneNotifier.mu.Unlock()
		delete(imageTaskDoneNotifier.waiters[taskID], waiter)
		if len(imageTaskDoneNotifier.waiters[taskID]) == 0 {
			delete(imageTaskDoneNotifier.waiters, taskID)
		}
	}
}

func startImageTaskDoneNotifier() {
	pubsub := common.RDB.PSubscribe(context.Background(), imageTaskDoneChannel+"*")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := pubsub.Receive(ctx)
	cancel()
	if err != nil {
		_ = pubsub.Close()
		return
	}
	imageTaskDoneNotifier.mu.Lock()
	imageTaskDoneNotifier.ready = true
	imageTaskDoneNotifier.mu.Unlock()
	go func() {
		for message := range pubsub.Channel() {
			taskID := strings.TrimPrefix(message.Channel, imageTaskDoneChannel)
			imageTaskDoneNotifier.mu.Lock()
			for waiter := range imageTaskDoneNotifier.waiters[taskID] {
				select {
				case waiter <- struct{}{}:
				default:
				}
			}
			imageTaskDoneNotifier.mu.Unlock()
		}
	}()
}

func imageTaskInputObjectKeys(task *model.Task) []string {
	if task == nil || len(task.PrivateData.RequestSnapshot) == 0 {
		return nil
	}
	return EditSnapshotObjectKeys(task.PrivateData.RequestSnapshot)
}

func cleanupImageTaskInputs(taskID string, objectKeys []string) {
	for _, objectKey := range objectKeys {
		if err := service.DeleteImageTaskInput(context.Background(), objectKey); err != nil {
			common.SysError(fmt.Sprintf("image task input cleanup failed for %s: %v", taskID, err))
		}
	}
}

func taskUpstreamBaseURL(task *model.Task) string {
	if task == nil || task.ChannelId == 0 {
		return ""
	}
	channel, err := model.GetChannelById(task.ChannelId, true)
	if err != nil || channel == nil {
		return ""
	}
	return channel.GetBaseURL()
}
