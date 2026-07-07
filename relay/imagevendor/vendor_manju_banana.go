package imagevendor

import "strings"

func init() {
	register(Descriptor{
		Name:  "manju-banana",
		Match: IsManjuBananaOriginModel,
		Rehost: RehostPolicy{
			AcceptUpstreamURL: true,
		},
	})
}

// IsManjuBananaOriginModel：Manju Gemini Banana 渠道 internal 名（manju-gemini-banana-*）。
func IsManjuBananaOriginModel(originModel string) bool {
	name := normalizeOriginModel(originModel)
	return strings.HasPrefix(name, "manju-gemini-banana")
}
