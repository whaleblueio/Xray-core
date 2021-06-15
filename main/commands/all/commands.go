package all

import (
	"github.com/whaleblueio/Xray-core/main/commands/all/api"
	"github.com/whaleblueio/Xray-core/main/commands/all/tls"
	"github.com/whaleblueio/Xray-core/main/commands/base"
)

// go:generate go run github.com/whaleblueio/Xray-core/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		//cmdConvert,
		tls.CmdTLS,
		cmdUUID,
	)
}
