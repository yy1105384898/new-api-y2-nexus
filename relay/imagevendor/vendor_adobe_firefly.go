package imagevendor

import (
	"net/url"
	"strings"
)

const adobeGeneratedPublicHost = "eu-ai.cangyuansuanli.cn"

func init() {
	register(Descriptor{
		Name:  "adobe-firefly",
		Match: IsAdobeFireflyOriginModel,
		Rehost: RehostPolicy{
			AcceptUpstreamURL:      true,
			AsyncPreferURLResponse: true,
			TrustPublicURL:         IsTrustedAdobeGeneratedURL,
		},
	})
}

// IsAdobeFireflyOriginModel matches the internal Adobe2API image SKU family.
// Adobe2API returns generated media as URLs. Only its owned HTTPS MinIO path
// may be exposed directly; every other upstream URL still follows R2 rehosting.
func IsAdobeFireflyOriginModel(originModel string) bool {
	return strings.HasPrefix(normalizeOriginModel(originModel), "adobe-firefly-")
}

func IsTrustedAdobeGeneratedURL(rawURL string) bool {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Scheme != "https" || u.Host != adobeGeneratedPublicHost || u.User != nil {
		return false
	}
	return strings.HasPrefix(u.EscapedPath(), "/generated/") && len(u.EscapedPath()) > len("/generated/")
}
