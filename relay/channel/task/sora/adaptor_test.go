package sora

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
)

func TestParseTaskResult_GZFormat(t *testing.T) {
	adaptor := &TaskAdaptor{}

	t.Run("running", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"running","videoUrl":null,"error":null}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusInProgress {
			t.Fatalf("expected IN_PROGRESS, got %s", result.Status)
		}
	})

	t.Run("succeeded with videoUrl", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"succeeded","videoUrl":"https://example.com/a.mp4","error":null}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusSuccess {
			t.Fatalf("expected SUCCESS, got %s", result.Status)
		}
		if result.Url != "https://example.com/a.mp4" {
			t.Fatalf("expected video url, got %q", result.Url)
		}
	})

	t.Run("failed with string error", func(t *testing.T) {
		result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"failed","videoUrl":null,"error":"content policy violation"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusFailure {
			t.Fatalf("expected FAILURE, got %s", result.Status)
		}
		if result.Reason != "content policy violation" {
			t.Fatalf("expected error reason, got %q", result.Reason)
		}
	})

	t.Run("grok moderation without status", func(t *testing.T) {
		body := []byte(`{"code":"Client specified an invalid argument","error":"Generated video rejected by content moderation.","id":"task_upstream","task_id":"task_upstream","model":"grok-image-video"}`)
		result, err := adaptor.ParseTaskResult(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != model.TaskStatusFailure {
			t.Fatalf("expected FAILURE, got %s", result.Status)
		}
		if result.Reason != "Generated video rejected by content moderation." {
			t.Fatalf("expected moderation reason, got %q", result.Reason)
		}
	})
}

func TestParseTaskResult_OpenAIFormat(t *testing.T) {
	adaptor := &TaskAdaptor{}

	result, err := adaptor.ParseTaskResult([]byte(`{"id":"vid1","status":"completed","usage":{"seconds":8}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != model.TaskStatusSuccess {
		t.Fatalf("expected SUCCESS, got %s", result.Status)
	}
	if result.CompletionTokens != 8 {
		t.Fatalf("expected 8 seconds, got %d", result.CompletionTokens)
	}
}
