package conf_test

import (
	"testing"

	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"github.com/whaleblueio/Xray-core/common/serial"
	. "github.com/whaleblueio/Xray-core/infra/conf"
	"github.com/whaleblueio/Xray-core/proxy/shadowsocks"
)

func TestShadowsocksServerConfigParsing(t *testing.T) {
	creator := func() Buildable {
		return new(ShadowsocksServerConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"method": "aes-128-gcm",
				"password": "xray-password"
			}`,
			Parser: loadJSON(creator),
			Output: &shadowsocks.ServerConfig{
				Users: []*protocol.User{{
					Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: shadowsocks.CipherType_AES_128_GCM,
						Password:   "xray-password",
					}),
				}},
				Network: []net.Network{net.Network_TCP},
			},
		},
	})
}
