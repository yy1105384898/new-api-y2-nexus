package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetModelRequestVideoFetchDoesNotSelectChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		path string
	}{
		{name: "video generations", path: "/v1/video/generations/task_abc123"},
		{name: "openai videos", path: "/v1/videos/task_abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, tt.path, nil)
			c.Params = gin.Params{{Key: "task_id", Value: "task_abc123"}}

			_, shouldSelectChannel, err := getModelRequest(c)
			require.NoError(t, err)
			require.False(t, shouldSelectChannel)
			require.Equal(t, relayconstant.RelayModeVideoFetchByID, c.GetInt("relay_mode"))
		})
	}
}

func TestGetModelRequestVideoSubmitSelectsChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := strings.NewReader(`{"model":"grok-video","prompt":"test"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/video/generations", body)
	c.Request.Header.Set("Content-Type", "application/json")

	_, shouldSelectChannel, err := getModelRequest(c)
	require.NoError(t, err)
	require.True(t, shouldSelectChannel)
}

func TestDistributeVideoTaskFetchSkipsModelLimitAndNilChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	reached := false
	router.GET("/v1/video/generations/:task_id",
		func(c *gin.Context) {
			common.SetContextKey(c, constant.ContextKeyTokenModelLimitEnabled, true)
			common.SetContextKey(c, constant.ContextKeyTokenModelLimit, map[string]bool{"grok-video": true})
		},
		Distribute(),
		func(c *gin.Context) {
			reached = true
			c.Status(http.StatusOK)
		},
	)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/video/generations/task_test123", nil)
	router.ServeHTTP(rec, req)

	require.True(t, reached, "task fetch should reach handler without channel selection")
	require.Equal(t, http.StatusOK, rec.Code)
}
