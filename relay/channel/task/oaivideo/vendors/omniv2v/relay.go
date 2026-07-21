package omniv2v

import "strings"

func IsRelay(originModel, upstreamModel string) bool {
	upstream := strings.ToLower(strings.TrimSpace(upstreamModel))
	if strings.Contains(upstream, "omni-fast-v2v") {
		return true
	}
	origin := strings.ToLower(strings.TrimSpace(originModel))
	return strings.Contains(origin, "omni-v2v")
}
