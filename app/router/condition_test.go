package router_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"

	. "github.com/whaleblueio/Xray-core/app/router"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/errors"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/platform"
	"github.com/whaleblueio/Xray-core/common/platform/filesystem"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"github.com/whaleblueio/Xray-core/common/protocol/http"
	"github.com/whaleblueio/Xray-core/common/session"
	"github.com/whaleblueio/Xray-core/features/routing"
	routing_session "github.com/whaleblueio/Xray-core/features/routing/session"
)

func init() {
	wd, err := os.Getwd()
	common.Must(err)

	if _, err := os.Stat(platform.GetAssetLocation("geoip.dat")); err != nil && os.IsNotExist(err) {
		common.Must(filesystem.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(wd, "..", "..", "release", "config", "geoip.dat")))
	}
	if _, err := os.Stat(platform.GetAssetLocation("geosite.dat")); err != nil && os.IsNotExist(err) {
		common.Must(filesystem.CopyFile(platform.GetAssetLocation("geosite.dat"), filepath.Join(wd, "..", "..", "release", "config", "geosite.dat")))
	}
}

func withBackground() routing.Context {
	return &routing_session.Context{}
}

func withOutbound(outbound *session.Outbound) routing.Context {
	return &routing_session.Context{Outbound: outbound}
}

func withInbound(inbound *session.Inbound) routing.Context {
	return &routing_session.Context{Inbound: inbound}
}

func withContent(content *session.Content) routing.Context {
	return &routing_session.Context{Content: content}
}

func TestRoutingRule(t *testing.T) {
	type ruleTest struct {
		input  routing.Context
		output bool
	}

	cases := []struct {
		rule *RoutingRule
		test []ruleTest
	}{
		{
			rule: &RoutingRule{
				Domain: []*Domain{
					{
						Value: "example.com",
						Type:  Domain_Plain,
					},
					{
						Value: "google.com",
						Type:  Domain_Domain,
					},
					{
						Value: "^facebook\\.com$",
						Type:  Domain_Regex,
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("example.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.example.com.www"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("example.co"), 80)}),
					output: false,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.google.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("facebook.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.facebook.com"), 80)}),
					output: false,
				},
				{
					input:  withBackground(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Cidr: []*CIDR{
					{
						Ip:     []byte{8, 8, 8, 8},
						Prefix: 32,
					},
					{
						Ip:     []byte{8, 8, 8, 8},
						Prefix: 32,
					},
					{
						Ip:     net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334").IP(),
						Prefix: 128,
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.4.4"), 80)}),
					output: false,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 80)}),
					output: true,
				},
				{
					input:  withBackground(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Geoip: []*GeoIP{
					{
						Cidr: []*CIDR{
							{
								Ip:     []byte{8, 8, 8, 8},
								Prefix: 32,
							},
							{
								Ip:     []byte{8, 8, 8, 8},
								Prefix: 32,
							},
							{
								Ip:     net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334").IP(),
								Prefix: 128,
							},
						},
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.4.4"), 80)}),
					output: false,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 80)}),
					output: true,
				},
				{
					input:  withBackground(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				SourceCidr: []*CIDR{
					{
						Ip:     []byte{192, 168, 0, 0},
						Prefix: 16,
					},
				},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{Source: net.TCPDestination(net.ParseAddress("192.168.0.1"), 80)}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Source: net.TCPDestination(net.ParseAddress("10.0.0.1"), 80)}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				UserEmail: []string{
					"admin@example.com",
				},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{User: &protocol.MemoryUser{Email: "admin@example.com"}}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{User: &protocol.MemoryUser{Email: "love@example.com"}}),
					output: false,
				},
				{
					input:  withBackground(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Protocol: []string{"http"},
			},
			test: []ruleTest{
				{
					input:  withContent(&session.Content{Protocol: (&http.SniffHeader{}).Protocol()}),
					output: true,
				},
			},
		},
		{
			rule: &RoutingRule{
				InboundTag: []string{"test", "test1"},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{Tag: "test"}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Tag: "test2"}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				PortList: &net.PortList{
					Range: []*net.PortRange{
						{From: 443, To: 443},
						{From: 1000, To: 1100},
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 443)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 1100)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 1005)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 53)}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				SourcePortList: &net.PortList{
					Range: []*net.PortRange{
						{From: 123, To: 123},
						{From: 9993, To: 9999},
					},
				},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{Source: net.UDPDestination(net.LocalHostIP, 123)}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Source: net.UDPDestination(net.LocalHostIP, 9999)}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Source: net.UDPDestination(net.LocalHostIP, 9994)}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Source: net.UDPDestination(net.LocalHostIP, 53)}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Protocol:   []string{"http"},
				Attributes: "attrs[':path'].startswith('/test')",
			},
			test: []ruleTest{
				{
					input:  withContent(&session.Content{Protocol: "http/1.1", Attributes: map[string]string{":path": "/test/1"}}),
					output: true,
				},
			},
		},
	}

	for _, test := range cases {
		cond, err := test.rule.BuildCondition()
		common.Must(err)

		for _, subtest := range test.test {
			actual := cond.Apply(subtest.input)
			if actual != subtest.output {
				t.Error("test case failed: ", subtest.input, " expected ", subtest.output, " but got ", actual)
			}
		}
	}
}

func loadGeoSite(country string) ([]*Domain, error) {
	geositeBytes, err := filesystem.ReadAsset("geosite.dat")
	if err != nil {
		return nil, err
	}
	var geositeList GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}

	for _, site := range geositeList.Entry {
		if site.CountryCode == country {
			return site.Domain, nil
		}
	}

	return nil, errors.New("country not found: " + country)
}

func TestChinaSites(t *testing.T) {
	domains, err := loadGeoSite("CN")
	common.Must(err)

	matcher, err := NewDomainMatcher(domains)
	common.Must(err)

	acMatcher, err := NewMphMatcherGroup(domains)
	common.Must(err)

	type TestCase struct {
		Domain string
		Output bool
	}
	testCases := []TestCase{
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "164.com",
			Output: false,
		},
		{
			Domain: "164.com",
			Output: false,
		},
	}

	for i := 0; i < 1024; i++ {
		testCases = append(testCases, TestCase{Domain: strconv.Itoa(i) + ".not-exists.com", Output: false})
	}

	for _, testCase := range testCases {
		r1 := matcher.ApplyDomain(testCase.Domain)
		r2 := acMatcher.ApplyDomain(testCase.Domain)
		if r1 != testCase.Output {
			t.Error("DomainMatcher expected output ", testCase.Output, " for domain ", testCase.Domain, " but got ", r1)
		} else if r2 != testCase.Output {
			t.Error("ACDomainMatcher expected output ", testCase.Output, " for domain ", testCase.Domain, " but got ", r2)
		}
	}
}

func BenchmarkMphDomainMatcher(b *testing.B) {
	domains, err := loadGeoSite("CN")
	common.Must(err)

	matcher, err := NewMphMatcherGroup(domains)
	common.Must(err)

	type TestCase struct {
		Domain string
		Output bool
	}
	testCases := []TestCase{
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "164.com",
			Output: false,
		},
		{
			Domain: "164.com",
			Output: false,
		},
	}

	for i := 0; i < 1024; i++ {
		testCases = append(testCases, TestCase{Domain: strconv.Itoa(i) + ".not-exists.com", Output: false})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			_ = matcher.ApplyDomain(testCase.Domain)
		}
	}
}

func BenchmarkDomainMatcher(b *testing.B) {
	domains, err := loadGeoSite("CN")
	common.Must(err)

	matcher, err := NewDomainMatcher(domains)
	common.Must(err)

	type TestCase struct {
		Domain string
		Output bool
	}
	testCases := []TestCase{
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "164.com",
			Output: false,
		},
		{
			Domain: "164.com",
			Output: false,
		},
	}

	for i := 0; i < 1024; i++ {
		testCases = append(testCases, TestCase{Domain: strconv.Itoa(i) + ".not-exists.com", Output: false})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			_ = matcher.ApplyDomain(testCase.Domain)
		}
	}
}

func BenchmarkMultiGeoIPMatcher(b *testing.B) {
	var geoips []*GeoIP

	{
		ips, err := loadGeoIP("CN")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "CN",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("JP")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "JP",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("CA")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "CA",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("US")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "US",
			Cidr:        ips,
		})
	}

	matcher, err := NewMultiGeoIPMatcher(geoips, false)
	common.Must(err)

	ctx := withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = matcher.Apply(ctx)
	}
}
