package clienterror

import "strings"

// Default OpenAI Video 聚合线路 (sora-2 等). Upstream: vendors/defaultvideo/adaptor.go, adaptor_test.go
// Typical raw: error string / Client specified an invalid argument + Generated video rejected...

func normalizeDefaultVideo(preferChinese bool, raw string) (string, bool) {
	lower := strings.ToLower(stripStatusCodePrefix(raw))

	if isGenericTaskFailed(raw) {
		return localized(preferChinese, GenerationFailedMessageZH, GenerationFailedMessageEN), true
	}
	if strings.Contains(lower, "client specified an invalid argument") {
		if IsContentPolicyViolation(raw) {
			return localized(preferChinese, ContentPolicyMessageZH, ContentPolicyMessageEN), true
		}
		return localized(preferChinese, InvalidRequestMessageZH, InvalidRequestMessageEN), true
	}
	if strings.Contains(lower, "invalid task_id") || strings.Contains(lower, "task_id is empty") {
		return localized(preferChinese, GenerationFailedMessageZH, GenerationFailedMessageEN), true
	}

	return "", false
}
