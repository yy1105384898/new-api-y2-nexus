package dto

const (
	ImageJobStatusQueued      = "queued"
	ImageJobStatusInProgress  = "in_progress"
	ImageJobStatusCompleted   = "completed"
	ImageJobStatusFailed      = "failed"
)

type OpenAIImageJob struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	Model     string            `json:"model,omitempty"`
	Status    string            `json:"status"`
	Progress  string            `json:"progress,omitempty"`
	CreatedAt int64             `json:"created_at"`
	Data      []ImageData       `json:"data,omitempty"`
	Error     *OpenAIImageError `json:"error,omitempty"`
}

type OpenAIImageError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func NewOpenAIImageJob(object string) *OpenAIImageJob {
	return &OpenAIImageJob{
		Object: object,
		Status: ImageJobStatusQueued,
	}
}
