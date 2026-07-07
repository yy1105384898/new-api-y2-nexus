package openai

import (
	"encoding/json"

	tasksora "github.com/QuantumNous/new-api/relay/channel/task/sora"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

// IsManjuSora2OriginModel Manju 渠道 #70 视频模型。
func IsManjuSora2OriginModel(originModel string) bool {
	return tasksora.IsManjuSora2Relay(originModel, "")
}

func convertManjuSora2OpenAIChatRequest(request *dto.GeneralOpenAIRequest, info *relaycommon.RelayInfo) (any, error) {
	raw, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	return tasksora.ConvertManjuSora2ChatBody(body, info.UpstreamModelName)
}

// manjuSora2PassthroughIfNeeded chat/completions 上游若返回 Manju task 对象，原样透传；失败时转为 OpenAI error。
func manjuSora2PassthroughIfNeeded(info *relaycommon.RelayInfo, responseBody []byte) ([]byte, bool) {
	if info == nil || !IsManjuSora2OriginModel(info.OriginModelName) {
		return responseBody, false
	}
	if errBody, ok := tasksora.BuildManjuSoraOpenAIErrorResponse(responseBody); ok {
		return errBody, true
	}
	if !tasksora.IsManjuSora2Response(responseBody) {
		if reason := tasksora.ExtractManjuSoraFailReasonForChat(responseBody); reason != "" {
			out, err := json.Marshal(map[string]any{
				"error": map[string]any{
					"message": reason,
					"type":    "upstream_error",
				},
			})
			if err == nil {
				return out, true
			}
		}
		return responseBody, false
	}
	return responseBody, true
}
