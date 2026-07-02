package setting

import (
	"github.com/QuantumNous/new-api/common"
)

var CanvasBaseURL = "https://canvas.yangyangnj.top"

func InitCanvasSetting() {
	CanvasBaseURL = common.GetEnvOrDefaultString("CANVAS_BASE_URL", CanvasBaseURL)
}
