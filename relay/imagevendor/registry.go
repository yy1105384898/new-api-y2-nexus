package imagevendor

import (
	"strings"

	"github.com/QuantumNous/new-api/dto"
)

var descriptors []Descriptor

func register(d Descriptor) {
	descriptors = append(descriptors, d)
}

// ResolveRehostPolicy 按注册顺序返回首个命中的转存策略；未命中时默认仅接受 b64_json。
func ResolveRehostPolicy(originModel string) RehostPolicy {
	for _, d := range descriptors {
		if d.Match(originModel) {
			return d.Rehost
		}
	}
	return RehostPolicy{}
}

// ImageModelUsesURLRehost：异步提交上游时优先 response_format=url。
func ImageModelUsesURLRehost(originModel string) bool {
	return ResolveRehostPolicy(originModel).AsyncPreferURLResponse
}

// ImageSyncPreferUpstreamB64JSON：同步生图对客户要 url 时，对内改请求上游 b64_json。
func ImageSyncPreferUpstreamB64JSON(originModel string) bool {
	return ResolveRehostPolicy(originModel).PreferUpstreamB64JSON
}

// ImageAsyncAcceptsUpstreamURL：允许上游回 url 并转存 R2。
func ImageAsyncAcceptsUpstreamURL(originModel string) bool {
	return ResolveRehostPolicy(originModel).AcceptUpstreamURL
}

// ApplyRequestPatch 按注册顺序对首个命中的 Descriptor 执行请求补丁；无 Patch 或未命中时 no-op。
func ApplyRequestPatch(originModel string, request *dto.ImageRequest) (RequestPatchResult, error) {
	for _, d := range descriptors {
		if d.PatchRequest == nil || !d.Match(originModel) {
			continue
		}
		return d.PatchRequest(originModel, request)
	}
	return RequestPatchResult{}, nil
}

func normalizeOriginModel(originModel string) string {
	return strings.ToLower(strings.TrimSpace(originModel))
}
