package conf

import (
	"github.com/golang/protobuf/proto"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/proxy/dns"
)

type DNSOutboundConfig struct {
	Network Network  `json:"network"`
	Address *Address `json:"address"`
	Port    uint16   `json:"port"`
}

func (c *DNSOutboundConfig) Build() (proto.Message, error) {
	config := &dns.Config{
		Server: &net.Endpoint{
			Network: c.Network.Build(),
			Port:    uint32(c.Port),
		},
	}
	if c.Address != nil {
		config.Server.Address = c.Address.Build()
	}
	return config, nil
}
