package vendorpatch

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
)

type gulieImagePatcher struct{}

func (gulieImagePatcher) Match(originModel string) bool {
	return matchModelPrefix(originModel, "gulie-")
}

func (gulieImagePatcher) Apply(request *dto.ImageRequest) (ImageTransformResult, error) {
	if request == nil {
		return ImageTransformResult{}, fmt.Errorf("gulie image patch: request is nil")
	}

	stripGulieUnsupportedFields(request)
	// 上游建议并发场景走非流式 JSON，SSE 长连接易在 relay 层崩坏。
	request.Stream = common.GetPointer(false)

	result := ImageTransformResult{SuppressQualityLog: true}
	if strings.EqualFold(strings.TrimSpace(request.Size), "auto") {
		request.Size = ""
	}
	return result, nil
}

func init() {
	registerImagePatcher(gulieImagePatcher{})
}

func stripGulieUnsupportedFields(request *dto.ImageRequest) {
	request.Quality = ""
	request.Background = nil
	request.Moderation = nil
	request.OutputFormat = nil
	request.OutputCompression = nil
}
