package gemini

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/stretchr/testify/require"
)

func TestBuildGeminiSafetySettings(t *testing.T) {
	settings := buildGeminiSafetySettings()
	require.Len(t, settings, len(SafetySettingList))
	for i, category := range SafetySettingList {
		require.Equal(t, category, settings[i].Category)
		require.Equal(t, "OFF", settings[i].Threshold)
	}
}

func TestIsGeminiImageGenerationRequest(t *testing.T) {
	tests := []struct {
		name    string
		model   string
		request *dto.GeminiChatRequest
		want    bool
	}{
		{
			name:  "supported imagine model",
			model: "gemini-2.5-flash-image",
			want:  true,
		},
		{
			name:  "image suffix model",
			model: "0lll0-gemini-3.1-flash-lite-image",
			want:  true,
		},
		{
			name:  "banana alias",
			model: "nano-banana-pro",
			want:  true,
		},
		{
			name:  "text model without image modality",
			model: "gemini-2.5-flash",
			want:  false,
		},
		{
			name:  "response modalities include image",
			model: "gemini-2.5-flash",
			request: &dto.GeminiChatRequest{
				GenerationConfig: dto.GeminiChatGenerationConfig{
					ResponseModalities: []string{"TEXT", "IMAGE"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.request
			if req == nil {
				req = &dto.GeminiChatRequest{}
			}
			info := &relaycommon.RelayInfo{
				ChannelMeta: &relaycommon.ChannelMeta{
					UpstreamModelName: tt.model,
				},
			}
			require.Equal(t, tt.want, isGeminiImageGenerationRequest(info, req))
		})
	}
}
