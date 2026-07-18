package sd5

import "strings"

const modelPrefix = "cy-sd5-seedance-2.0"

func IsRelay(originModel, upstreamModel string) bool {
	for _, model := range []string{originModel, upstreamModel} {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), modelPrefix) {
			return true
		}
	}
	return false
}
