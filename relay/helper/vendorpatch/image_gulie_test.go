package vendorpatch

import (
	"encoding/json"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/require"
)

func TestGulieImagePatcherMatch(t *testing.T) {
	p := gulieImagePatcher{}
	require.True(t, p.Match("gulie-gpt-image-2"))
	require.False(t, p.Match("kedaya-gpt-image-2"))
}

func TestGulieImagePatcherApplyStripsUnsupportedFields(t *testing.T) {
	p := gulieImagePatcher{}
	request := &dto.ImageRequest{
		Model:             "gpt-image-2",
		Prompt:            "test",
		Size:              "1536x1024",
		Quality:           "high",
		Background:        json.RawMessage(`"transparent"`),
		OutputFormat:      json.RawMessage(`"png"`),
		OutputCompression: json.RawMessage(`80`),
		Moderation:        json.RawMessage(`"low"`),
	}

	result, err := p.Apply(request)
	require.NoError(t, err)
	require.True(t, result.SuppressQualityLog)
	require.Equal(t, "1536x1024", request.Size)
	require.Empty(t, request.Quality)
	require.Nil(t, request.Background)
	require.Nil(t, request.OutputFormat)
	require.Nil(t, request.OutputCompression)
	require.Nil(t, request.Moderation)
	require.NotNil(t, request.Stream)
	require.False(t, *request.Stream)
}

func TestGulieImagePatcherApplyAutoSize(t *testing.T) {
	p := gulieImagePatcher{}
	request := &dto.ImageRequest{Size: "auto"}
	_, err := p.Apply(request)
	require.NoError(t, err)
	require.Empty(t, request.Size)
}
