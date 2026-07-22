package common

import "fmt"

// Canonical reference-media byte limits shared by relay validation, UI profiles,
// and upstream Adobe2API deployments. Keep adobe2api/core/media_limits.py in sync.
const (
	ReferenceImageMaxBytes = 30 << 20
	ReferenceVideoMaxBytes = 50 << 20
	ReferenceAudioMaxBytes = 15 << 20
)

func ReferenceImageMaxMegabytes() int {
	return int(ReferenceImageMaxBytes / (1 << 20))
}

func ReferenceVideoMaxMegabytes() int {
	return int(ReferenceVideoMaxBytes / (1 << 20))
}

func ReferenceAudioMaxMegabytes() int {
	return int(ReferenceAudioMaxBytes / (1 << 20))
}

func ReferenceImageTooLargeDetail() string {
	return fmt.Sprintf("image too large, max %dMB", ReferenceImageMaxMegabytes())
}

func ReferenceMaskTooLargeDetail() string {
	return fmt.Sprintf("mask too large, max %dMB", ReferenceImageMaxMegabytes())
}

func ReferenceFieldTooLargeDetail(fieldName string) string {
	return fmt.Sprintf("%s too large, max %dMB", fieldName, ReferenceImageMaxMegabytes())
}

func ReferenceVideoTooLargeDetail() string {
	return fmt.Sprintf("video too large, max %dMB", ReferenceVideoMaxMegabytes())
}

func ReferenceAudioTooLargeDetail() string {
	return fmt.Sprintf("audio too large, max %dMB", ReferenceAudioMaxMegabytes())
}

func ReferenceMediaLabelZH(kind string) string {
	switch kind {
	case "image":
		return "参考图"
	case "video":
		return "参考视频"
	case "audio":
		return "参考音频"
	case "mask":
		return "蒙版"
	case "entity", "entity_image", "entity image":
		return "主体参考图"
	default:
		return "参考素材"
	}
}

func ReferenceMediaLabelEN(kind string) string {
	switch kind {
	case "image":
		return "reference image"
	case "video":
		return "reference video"
	case "audio":
		return "reference audio"
	case "mask":
		return "mask"
	case "entity", "entity_image", "entity image":
		return "entity image"
	default:
		return "reference media"
	}
}

func FormatReferenceByteLimitMessageZH(kind string, maxBytes int64) string {
	return fmt.Sprintf("%s超过 %dMB，请压缩后再上传", ReferenceMediaLabelZH(kind), maxBytes/(1<<20))
}

func FormatReferenceByteLimitMessageEN(kind string, maxBytes int64) string {
	return fmt.Sprintf("%s exceeds %dMB. Please compress and retry.", ReferenceMediaLabelEN(kind), maxBytes/(1<<20))
}
