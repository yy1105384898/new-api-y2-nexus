package clienterror

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Grok video (119337 / cy-gv1) and Geeknow Grok. Upstream:
//   - relay/channel/task/oaivideo/vendors/grok/adaptor.go
//   - relay/channel/task/oaivideo/vendors/geeknowgrok/adaptor.go

var grokReferenceImagesRe = regexp.MustCompile(`(?i)(?:grok|geeknow grok) video supports at most (\d+) reference images`)

func normalizeGrok(preferChinese bool, raw string) (string, bool) {
	lower := strings.ToLower(stripStatusCodePrefix(raw))

	if strings.Contains(lower, "grok-video-1.5 requires exactly one reference image") ||
		strings.Contains(lower, "grok-imagine-video-1.5-preview supports at most one reference image") {
		if preferChinese {
			return "该 Grok 模型仅支持 1 张参考图，请减少后重试。", true
		}
		return "This Grok model requires exactly one reference image.", true
	}
	if strings.Contains(lower, "grok-video-1.5 does not support video references") ||
		strings.Contains(lower, "grok-imagine-video-1.5-preview does not support video references") {
		if preferChinese {
			return "该 Grok 模型不支持参考视频，请移除视频参考。", true
		}
		return "This Grok model does not support reference videos.", true
	}
	if m := grokReferenceImagesRe.FindStringSubmatch(raw); len(m) == 2 {
		max, _ := strconv.Atoi(m[1])
		if max > 0 {
			if preferChinese {
				return fmt.Sprintf("参考图最多 %d 张，请减少后重试。", max), true
			}
			return fmt.Sprintf("At most %d reference images allowed. Please remove extras and retry.", max), true
		}
	}

	return "", false
}