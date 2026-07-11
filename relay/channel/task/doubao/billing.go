package doubao

import (
	"strconv"
	"strings"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

// 官方 5s/16:9 示例价（元）→ 720p 无参考为基准 1.0，见火山方舟定价文档。
var resolutionRatio720p = map[string]map[string]float64{
	"doubao-seedance-2-0-260128": {
		"480p":  2.31 / 4.97,
		"720p":  1,
		"1080p": 12.39 / 4.97,
	},
	"doubao-seedance-2-0-fast-260128": {
		"480p": 1.86 / 4.00,
		"720p": 1,
	},
	"doubao-seedance-2-0-mini": {
		"480p": 1.16 / 2.48,
		"720p": 1,
	},
}

// videoInputCostRatio 含视频参考相对 720p 无参考的官方 5s 示例价倍率（元/秒）。
var videoInputCostRatio = map[string]float64{
	"doubao-seedance-2-0-260128":      5.44 / 4.97,
	"doubao-seedance-2-0-fast-260128": 4.28 / 4.00,
	"doubao-seedance-2-0-mini":        2.72 / 2.48,
}

func isOfficialPricingModel(modelName string) bool {
	if modelName == "" {
		return false
	}
	if _, ok := resolutionRatio720p[modelName]; ok {
		return true
	}
	_, ok := videoInputCostRatio[modelName]
	return ok
}

// resolvePricingModel 仅识别官方 doubao-seedance-* 模型名，不做渠道别名映射。
func resolvePricingModel(origin, upstream string) string {
	if isOfficialPricingModel(upstream) {
		return upstream
	}
	if isOfficialPricingModel(origin) {
		return origin
	}
	return upstream
}

func normalizeResolution(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "720p"
	}
	if !strings.HasSuffix(s, "p") {
		s += "p"
	}
	return s
}

func parseResolution(req *relaycommon.TaskSubmitReq) string {
	if req == nil {
		return "720p"
	}
	if req.Metadata != nil {
		if v, ok := req.Metadata["resolution"].(string); ok && v != "" {
			return normalizeResolution(v)
		}
	}
	if req.Size != "" {
		return normalizeResolution(req.Size)
	}
	return "720p"
}

func parseDurationSeconds(req *relaycommon.TaskSubmitReq) int {
	if req == nil {
		return 5
	}
	if seconds := req.RequestedDurationSeconds(); seconds > 0 {
		return seconds
	}
	if req.Metadata != nil {
		if v, ok := req.Metadata["duration"]; ok {
			switch d := v.(type) {
			case float64:
				if int(d) > 0 {
					return int(d)
				}
			case int:
				if d > 0 {
					return d
				}
			case string:
				if sec, err := strconv.Atoi(d); err == nil && sec > 0 {
					return sec
				}
			}
		}
	}
	return 5
}

func GetResolutionRatio(pricingModel, resolution string) float64 {
	table, ok := resolutionRatio720p[pricingModel]
	if !ok {
		return 1
	}
	res := normalizeResolution(resolution)
	if r, ok := table[res]; ok {
		return r
	}
	return 1
}

// GetVideoInputCostRatio 含视频参考相对 720p 无参考的官方案例倍率（5s 最低价）。
func GetVideoInputCostRatio(pricingModel string) (float64, bool) {
	r, ok := videoInputCostRatio[pricingModel]
	return r, ok
}
