package middleware

import (
	"net/http"
	"sync"

	"github.com/QuantumNous/new-api/constant"
	"github.com/gin-gonic/gin"
)

type imageInFlightCounter struct {
	mu     sync.Mutex
	global int
	users  map[int]int
}

var activeImageRequests = imageInFlightCounter{users: make(map[int]int)}

func (l *imageInFlightCounter) tryAcquire(userID, perUserLimit, globalLimit int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if globalLimit > 0 && l.global >= globalLimit {
		return false
	}
	if perUserLimit > 0 && l.users[userID] >= perUserLimit {
		return false
	}

	l.global++
	l.users[userID]++
	return true
}

func (l *imageInFlightCounter) release(userID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.global > 0 {
		l.global--
	}
	if count := l.users[userID]; count <= 1 {
		delete(l.users, userID)
	} else {
		l.users[userID] = count - 1
	}
}

func isImageRelayRequest(c *gin.Context) bool {
	if c.Request.Method != http.MethodPost {
		return false
	}
	switch c.Request.URL.Path {
	case "/v1/images/generations", "/v1/images/edits", "/v1/edits":
		return true
	default:
		return false
	}
}

// ImageInFlightLimit bounds active image generation/edit requests before any
// middleware parses or copies a potentially large multipart request body.
func ImageInFlightLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isImageRelayRequest(c) {
			c.Next()
			return
		}

		perUserLimit := constant.ImageMaxInFlightPerUser
		globalLimit := constant.ImageMaxInFlightGlobal
		if perUserLimit <= 0 && globalLimit <= 0 {
			c.Next()
			return
		}

		userID := c.GetInt("id")
		if !activeImageRequests.tryAcquire(userID, perUserLimit, globalLimit) {
			abortWithOpenAiMessage(c, http.StatusTooManyRequests, "too many concurrent image requests; please retry later")
			return
		}
		defer activeImageRequests.release(userID)

		c.Next()
	}
}
