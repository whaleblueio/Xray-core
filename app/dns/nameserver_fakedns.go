package dns

import (
	"context"

	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/core"
	"github.com/whaleblueio/Xray-core/features/dns"
)

type FakeDNSServer struct {
	fakeDNSEngine dns.FakeDNSEngine
}

func NewFakeDNSServer() *FakeDNSServer {
	return &FakeDNSServer{}
}

func (FakeDNSServer) Name() string {
	return "FakeDNS"
}

func (f *FakeDNSServer) QueryIP(ctx context.Context, domain string, _ dns.IPOption) ([]net.IP, error) {
	if f.fakeDNSEngine == nil {
		if err := core.RequireFeatures(ctx, func(fd dns.FakeDNSEngine) {
			f.fakeDNSEngine = fd
		}); err != nil {
			return nil, newError("Unable to locate a fake DNS Engine").Base(err).AtError()
		}
	}
	ips := f.fakeDNSEngine.GetFakeIPForDomain(domain)

	netIP := toNetIP(ips)
	if netIP == nil {
		return nil, newError("Unable to convert IP to net ip").AtError()
	}

	newError(f.Name(), " got answer: ", domain, " -> ", ips).AtInfo().WriteToLog()

	return netIP, nil
}
