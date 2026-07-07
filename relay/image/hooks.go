package image

import (
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

var (
	getAdaptor func(apiType int) channel.Adaptor
	textRelay  func(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError
)

// SetGetAdaptor 由 relay 包 init 注入，避免 image ↔ relay 循环依赖。
func SetGetAdaptor(fn func(apiType int) channel.Adaptor) {
	getAdaptor = fn
}

// SetTextRelay 由 relay 包 init 注入，供 legacy async chat 出图重放。
func SetTextRelay(fn func(c *gin.Context, info *relaycommon.RelayInfo) *types.NewAPIError) {
	textRelay = fn
}
