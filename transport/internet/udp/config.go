package udp

import (
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
