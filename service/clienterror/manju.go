package clienterror

import "strings"

// Manju Sora2 (manju-openai-sora*). Upstream: vendors/manju/convert.go, convert_test.go
// 多数 Manju 中文审核文案已由 common 匹配；此处兜底 generic task failed。

func normalizeManju(preferChinese bool, raw string) (string, bool) {
	if isGenericTaskFailed(raw) {
		return localized(preferChinese, GenerationFailedMessageZH, GenerationFailedMessageEN), true
	}
	if strings.Contains(raw, "某张上传的参考图未通过平台内容审核") {
		if IsRealFaceReferenceError(raw) {
			return localized(preferChinese, ReferenceRealFaceMessageZH, ReferenceRealFaceMessageEN), true
		}
		return localized(preferChinese, ContentPolicyMessageZH, ContentPolicyMessageEN), true
	}
	return "", false
}
