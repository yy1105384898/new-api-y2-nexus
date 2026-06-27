package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskListSelectLoadsResultURLWithoutSnapshot(t *testing.T) {
	truncateTables(t)

	largeSnapshot := []byte(strings.Repeat("A", 1024*512))
	task := &Task{
		UserId:     1,
		TaskID:     "task_list_select",
		Platform:   "openai",
		Status:     TaskStatusSuccess,
		SubmitTime: time.Now().Unix(),
		PrivateData: TaskPrivateData{
			ResultURL:       "https://example.com/result.png",
			RequestSnapshot: largeSnapshot,
		},
		Data: json.RawMessage(`{"status":"ok"}`),
	}
	insertTask(t, task)

	items := TaskGetAllTasks(0, 10, SyncTaskQueryParams{})
	require.Len(t, items, 1)
	assert.Equal(t, "https://example.com/result.png", items[0].GetResultURL())
	assert.Empty(t, items[0].PrivateData.RequestSnapshot)
	assert.JSONEq(t, `{"status":"ok"}`, string(items[0].Data))
}

func TestTaskListSelectUserQueryOmitsChannelID(t *testing.T) {
	truncateTables(t)

	task := &Task{
		UserId:     42,
		ChannelId:  99,
		TaskID:     "task_user_list",
		Platform:   "openai",
		Status:     TaskStatusSuccess,
		SubmitTime: time.Now().Unix(),
		PrivateData: TaskPrivateData{
			ResultURL: "https://example.com/video.mp4",
		},
		Data: json.RawMessage(`{}`),
	}
	insertTask(t, task)

	items := TaskGetAllUserTask(42, 0, 10, SyncTaskQueryParams{})
	require.Len(t, items, 1)
	assert.Equal(t, 0, items[0].ChannelId)
	assert.Equal(t, "https://example.com/video.mp4", items[0].GetResultURL())
}
