package shared

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
)

func TestParseResponseTaskAcceptsIntegerAndFractionalUnixTimestamps(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		createdAt   int64
		completedAt int64
		expiresAt   int64
	}{
		{
			name:        "integer seconds",
			body:        `{"id":"task-1","created_at":1783852874,"completed_at":1783853041,"expires_at":1783856641}`,
			createdAt:   1783852874,
			completedAt: 1783853041,
			expiresAt:   1783856641,
		},
		{
			name:        "fractional seconds",
			body:        `{"id":"task-2","created_at":1783852874.2331223,"completed_at":1783853041.3966722,"expires_at":1783856641.9}`,
			createdAt:   1783852874,
			completedAt: 1783853041,
			expiresAt:   1783856641,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResponseTask([]byte(tt.body))
			if err != nil {
				t.Fatalf("ParseResponseTask() error = %v", err)
			}
			if got := int64(result.CreatedAt); got != tt.createdAt {
				t.Fatalf("CreatedAt = %d, want %d", got, tt.createdAt)
			}
			if got := int64(result.CompletedAt); got != tt.completedAt {
				t.Fatalf("CompletedAt = %d, want %d", got, tt.completedAt)
			}
			if got := int64(result.ExpiresAt); got != tt.expiresAt {
				t.Fatalf("ExpiresAt = %d, want %d", got, tt.expiresAt)
			}
		})
	}
}

func TestFractionalUnixTimestampMarshalsAsIntegerSeconds(t *testing.T) {
	result, err := ParseResponseTask([]byte(`{"id":"task-1","created_at":1783852874.2331223}`))
	if err != nil {
		t.Fatalf("ParseResponseTask() error = %v", err)
	}
	body, err := common.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var decoded map[string]any
	if err := common.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got := decoded["created_at"]; got != float64(1783852874) {
		t.Fatalf("created_at = %#v, want integer seconds", got)
	}
}
