package image

import (
	"math"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/dto"
)

// channelImageSize4KDowngradeID：Gulie 渠道 72，发上游前按 OpenAI gpt-image-2 标准压到 2K。
const channelImageSize4KDowngradeID = 72

// openAIImage2KMaxPixels：OpenAI 2K 可靠上限（2560×1440）。
// 超过此总像素视为 4K/experimental，见 OpenAI Image generation 文档。
const openAIImage2KMaxPixels = 2560 * 1440

// openAIImage4KTo2KSize：OpenAI 官方 popular sizes 中的 4K → 2K 对。
var openAIImage4KTo2KSize = map[string]string{
	"3840x2160": "2048x1152",
	"2160x3840": "1152x2048",
	"2880x2880": "2048x2048",
	"4096x4096": "2048x2048",
	"4k":        "2048x1152",
	"16:9-4k":   "2048x1152",
	"9:16-4k":   "1152x2048",
	"1:1-4k":    "2048x2048",
}

// ApplyChannelImageSizeDowngrade 渠道 72：OpenAI 4K size 映射或缩到 2K 像素上限内。
func ApplyChannelImageSizeDowngrade(channelID int, request *dto.ImageRequest) bool {
	if channelID != channelImageSize4KDowngradeID || request == nil {
		return false
	}
	downgraded, ok := downgradeImageSizeForChannel72(request.Size)
	if !ok {
		return false
	}
	request.Size = downgraded
	return true
}

func downgradeImageSizeForChannel72(size string) (string, bool) {
	trimmed := strings.TrimSpace(size)
	if trimmed == "" {
		return "", false
	}
	normalized := strings.ToLower(trimmed)
	if mapped, ok := openAIImage4KTo2KSize[normalized]; ok {
		if mapped == trimmed {
			return trimmed, false
		}
		return mapped, true
	}
	if mapped, ok := openAIImage4KTo2KSize[trimmed]; ok {
		if mapped == trimmed {
			return trimmed, false
		}
		return mapped, true
	}

	width, height, ok := parseImagePixelSize(trimmed)
	if !ok {
		return "", false
	}
	return scaleImagePixelsToOpenAI2KMax(width, height)
}

// scaleImagePixelsToOpenAI2KMax：w×h 超过 OpenAI 2K 上限时，等比缩到 3686400 并 16px 对齐。
func scaleImagePixelsToOpenAI2KMax(width, height int) (string, bool) {
	if width <= 0 || height <= 0 {
		return "", false
	}
	area := width * height
	if area <= openAIImage2KMaxPixels {
		return "", false
	}
	scale := math.Sqrt(float64(openAIImage2KMaxPixels) / float64(area))
	scaledWidth := alignOpenAIPixelDimension(int(math.Round(float64(width) * scale)))
	scaledHeight := alignOpenAIPixelDimension(int(math.Round(float64(height) * scale)))
	if scaledWidth < 16 || scaledHeight < 16 {
		return "", false
	}
	// 对齐后可能略超上限，再收一档
	for scaledWidth*scaledHeight > openAIImage2KMaxPixels && scaledWidth > 16 && scaledHeight > 16 {
		if scaledWidth >= scaledHeight {
			scaledWidth -= 16
		} else {
			scaledHeight -= 16
		}
	}
	if scaledWidth == width && scaledHeight == height {
		return "", false
	}
	return formatImagePixelSize(scaledWidth, scaledHeight), true
}

func alignOpenAIPixelDimension(n int) int {
	if n < 16 {
		return 16
	}
	return (n / 16) * 16
}

func parseImagePixelSize(size string) (int, int, bool) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(size)), "x")
	if len(parts) != 2 {
		return 0, 0, false
	}
	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || width <= 0 {
		return 0, 0, false
	}
	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || height <= 0 {
		return 0, 0, false
	}
	return width, height, true
}

func formatImagePixelSize(width, height int) string {
	return strconv.Itoa(width) + "x" + strconv.Itoa(height)
}
