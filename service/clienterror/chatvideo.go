package clienterror

import "strings"

// Chat-completions video (cy-vid2-sora-2 等). Upstream: vendors/chatvideo/adaptor.go

func normalizeChatVideo(preferChinese bool, raw string) (string, bool) {
	lower := strings.ToLower(stripStatusCodePrefix(raw))

	if strings.Contains(lower, "empty chat video response") {
		return localized(preferChinese, GenerationFailedMessageZH, GenerationFailedMessageEN), true
	}
	if strings.Contains(lower, "chat video response does not contain a video url") {
		if preferChinese {
			return "视频生成失败，未返回视频地址，请稍后重试。", true
		}
		return "Video generation failed: upstream did not return a video URL. Please retry later.", true
	}

	return "", false
}
