package scenarios

import (
	"fmt"
	"testing"
	"time"

	"github.com/whaleblueio/Xray-core/app/dns"
	"github.com/whaleblueio/Xray-core/app/proxyman"
	"github.com/whaleblueio/Xray-core/app/router"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/serial"
	"github.com/whaleblueio/Xray-core/core"
	"github.com/whaleblueio/Xray-core/proxy/blackhole"
	"github.com/whaleblueio/Xray-core/proxy/freedom"
	"github.com/whaleblueio/Xray-core/proxy/socks"
	"github.com/whaleblueio/Xray-core/testing/servers/tcp"
	xproxy "golang.org/x/net/proxy"
)

func TestResolveIP(t *testing.T) {
	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close()

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dns.Config{
				Hosts: map[string]*net.IPOrDomain{
					"google.com": net.NewIPOrDomain(dest.Address),
				},
			}),
			serial.ToTypedMessage(&router.Config{
				DomainStrategy: router.Config_IpIfNonMatch,
				Rule: []*router.RoutingRule{
					{
						Cidr: []*router.CIDR{
							{
								Ip:     []byte{127, 0, 0, 0},
								Prefix: 8,
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "direct",
						},
					},
				},
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_NO_AUTH,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&blackhole.Config{}),
			},
			{
				Tag: "direct",
				ProxySettings: serial.ToTypedMessage(&freedom.Config{
					DomainStrategy: freedom.Config_USE_IP,
				}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig)
	common.Must(err)
	defer CloseAllServers(servers)

	{
		noAuthDialer, err := xproxy.SOCKS5("tcp", net.TCPDestination(net.LocalHostIP, serverPort).NetAddr(), nil, xproxy.Direct)
		common.Must(err)
		conn, err := noAuthDialer.Dial("tcp", fmt.Sprintf("google.com:%d", dest.Port))
		common.Must(err)
		defer conn.Close()

		if err := testTCPConn2(conn, 1024, time.Second*5)(); err != nil {
			t.Error(err)
		}
	}
}
