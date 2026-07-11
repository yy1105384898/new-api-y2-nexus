package chatvideo

import "strings"

// IsRelay identifies internal routes whose upstream still exposes video
// generation through chat/completions. Public model names are intentionally
// not matched here; routing uses the internal channel prefix contract.
func IsRelay(originModel string) bool {
	model := strings.ToLower(strings.TrimSpace(originModel))
	return strings.HasPrefix(model, "cy-vid2-") ||
		strings.HasPrefix(model, "yunwu-") ||
		model == "cy-sd1-grok-video" ||
		model == "oairegbox-grok-video"
}
