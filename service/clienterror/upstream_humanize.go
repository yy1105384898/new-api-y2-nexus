package clienterror

import "strings"

// Cross-channel upstream unavailable / HTTP / capacity errors.
// Vendor-specific pool humanization lives in leonardo.go, adobe.go, etc.

func humanizeUpstreamUnavailableError(preferChinese bool, raw string) (string, bool) {
	if !IsUpstreamUnavailableError(raw) {
		return "", false
	}
	lower := strings.ToLower(stripStatusCodePrefix(raw))

	switch {
	case strings.Contains(lower, "no capacity available") || strings.Contains(lower, "capacity available for model"):
		if preferChinese {
			return "模型容量已满，请稍后重试。", true
		}
		return "Model capacity is full. Please retry later.", true
	case strings.Contains(lower, "model overloaded"):
		if preferChinese {
			return "模型过载，请稍后重试。", true
		}
		return "The model is overloaded. Please retry later.", true
	case strings.Contains(lower, "no active tokens available"):
		if preferChinese {
			return "号池无可用账号，请联系管理员补充。可先缩短视频秒数或改用更省积分的模型再试。", true
		}
		return "No active pool accounts available. Contact an administrator, or try a shorter duration or economy model.", true
	case strings.Contains(lower, "all available tokens are invalid or expired"):
		if preferChinese {
			return "号池账号均已失效，请联系管理员。", true
		}
		return "All pool accounts are invalid or expired. Please contact an administrator.", true
	case strings.Contains(lower, "upstream is temporarily unavailable") || strings.Contains(lower, "upstream service temporarily unavailable"):
		return localized(preferChinese, UpstreamUnavailableMessageZH, UpstreamUnavailableMessageEN), true
	case strings.Contains(lower, "video service is temporarily unavailable"):
		return localized(preferChinese, PoolUnavailableMessageZH, PoolUnavailableMessageEN), true
	case strings.HasPrefix(raw, "status_code=502") || strings.Contains(lower, "bad response status code 502"):
		if preferChinese {
			return "网关错误（502），请稍后重试。", true
		}
		return "Gateway error (502). Please retry later.", true
	case strings.HasPrefix(raw, "status_code=503") || strings.Contains(lower, "bad response status code 503"):
		if preferChinese {
			return "服务不可用（503），请稍后重试。", true
		}
		return "Service unavailable (503). Please retry later.", true
	case strings.HasPrefix(raw, "status_code=504") || strings.Contains(lower, "bad response status code 504"):
		if preferChinese {
			return "网关超时（504），请稍后重试。", true
		}
		return "Gateway timeout (504). Please retry later.", true
	case strings.Contains(lower, "connection reset by peer"):
		if preferChinese {
			return "连接被重置，请稍后重试。", true
		}
		return "Connection was reset. Please retry later.", true
	case strings.Contains(lower, "connection refused"):
		if preferChinese {
			return "连接被拒绝，请稍后重试。", true
		}
		return "Connection was refused. Please retry later.", true
	case strings.Contains(lower, "download image failed") || strings.Contains(lower, "rehost upstream image url"):
		if preferChinese {
			return "参考图下载或转存失败，请检查图片链接后重试。", true
		}
		return "Reference image download or rehost failed. Check the image URL and retry.", true
	case strings.Contains(lower, "upstream request failed"):
		return localized(preferChinese, UpstreamUnavailableMessageZH, UpstreamUnavailableMessageEN), true
	}

	return localized(preferChinese, UpstreamUnavailableMessageZH, UpstreamUnavailableMessageEN), true
}
