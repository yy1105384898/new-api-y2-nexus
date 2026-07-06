package imagevendor

import "strings"

func init() {
	register(Descriptor{
		Name: "large-url-image",
		Match: func(originModel string) bool {
			name := normalizeOriginModel(originModel)
			return strings.HasSuffix(name, "-4k") || strings.HasPrefix(name, "flux-")
		},
		Rehost: RehostPolicy{
			AcceptUpstreamURL:      true,
			AsyncPreferURLResponse: true,
		},
	})
}
