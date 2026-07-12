package grok

import "strings"

func IsRelay(originModel, upstreamModel string) bool {
	origin := strings.ToLower(strings.TrimSpace(originModel))
	upstream := strings.ToLower(strings.TrimSpace(upstreamModel))
	return strings.HasPrefix(origin, "cy-gv1-grok-video") ||
		strings.HasPrefix(origin, "119337-grok-video") ||
		upstream == "grok-image-video" || upstream == "grok-video-1.5"
}

func isGrok15(originModel, upstreamModel string) bool {
	return strings.Contains(strings.ToLower(originModel), "1.5") || strings.EqualFold(strings.TrimSpace(upstreamModel), "grok-video-1.5")
}
