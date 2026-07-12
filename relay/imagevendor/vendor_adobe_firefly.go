package imagevendor

import "strings"

func init() {
	register(Descriptor{
		Name:  "adobe-firefly",
		Match: IsAdobeFireflyOriginModel,
		Rehost: RehostPolicy{
			AcceptUpstreamURL:      true,
			AsyncPreferURLResponse: true,
		},
	})
}

// IsAdobeFireflyOriginModel matches the internal Adobe2API image SKU family.
// Adobe2API returns generated media as URLs, which must be downloaded and
// rehosted to R2 before the task result is exposed to clients.
func IsAdobeFireflyOriginModel(originModel string) bool {
	return strings.HasPrefix(normalizeOriginModel(originModel), "adobe-firefly-")
}
