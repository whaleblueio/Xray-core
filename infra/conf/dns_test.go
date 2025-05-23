package conf_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/whaleblueio/Xray-core/app/dns"
	"github.com/whaleblueio/Xray-core/app/router"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/platform"
	"github.com/whaleblueio/Xray-core/common/platform/filesystem"
	. "github.com/whaleblueio/Xray-core/infra/conf"
)

func init() {
	wd, err := os.Getwd()
	common.Must(err)

	if _, err := os.Stat(platform.GetAssetLocation("geoip.dat")); err != nil && os.IsNotExist(err) {
		common.Must(filesystem.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(wd, "..", "..", "resources", "geoip.dat")))
	}

	geositeFilePath := filepath.Join(wd, "geosite.dat")
	os.Setenv("xray.location.asset", wd)
	geositeFile, err := os.OpenFile(geositeFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	common.Must(err)
	defer geositeFile.Close()

	list := &router.GeoSiteList{
		Entry: []*router.GeoSite{
			{
				CountryCode: "TEST",
				Domain: []*router.Domain{
					{Type: router.Domain_Full, Value: "example.com"},
				},
			},
		},
	}

	listBytes, err := proto.Marshal(list)
	common.Must(err)
	common.Must2(geositeFile.Write(listBytes))
}
func TestDNSConfigParsing(t *testing.T) {
	geositePath := platform.GetAssetLocation("geosite.dat")
	defer func() {
		os.Remove(geositePath)
		os.Unsetenv("xray.location.asset")
	}()

	parserCreator := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(DNSConfig)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"servers": [{
					"address": "8.8.8.8",
					"port": 5353,
					"domains": ["domain:example.com"]
				}],
				"hosts": {
					"example.com": "127.0.0.1",
					"domain:example.com": "google.com",
					"geosite:test": "10.0.0.1",
					"keyword:google": "8.8.8.8",
					"regexp:.*\\.com": "8.8.4.4"
				},
				"clientIp": "10.0.0.1"
			}`,
			Parser: parserCreator(),
			Output: &dns.Config{
				NameServer: []*dns.NameServer{
					{
						Address: &net.Endpoint{
							Address: &net.IPOrDomain{
								Address: &net.IPOrDomain_Ip{
									Ip: []byte{8, 8, 8, 8},
								},
							},
							Network: net.Network_UDP,
							Port:    5353,
						},
						PrioritizedDomain: []*dns.NameServer_PriorityDomain{
							{
								Type:   dns.DomainMatchingType_Subdomain,
								Domain: "example.com",
							},
						},
						OriginalRules: []*dns.NameServer_OriginalRule{
							{
								Rule: "domain:example.com",
								Size: 1,
							},
						},
					},
				},
				StaticHosts: []*dns.Config_HostMapping{
					{
						Type:          dns.DomainMatchingType_Subdomain,
						Domain:        "example.com",
						ProxiedDomain: "google.com",
					},
					{
						Type:   dns.DomainMatchingType_Full,
						Domain: "example.com",
						Ip:     [][]byte{{127, 0, 0, 1}},
					},
					{
						Type:   dns.DomainMatchingType_Full,
						Domain: "example.com",
						Ip:     [][]byte{{10, 0, 0, 1}},
					},
					{
						Type:   dns.DomainMatchingType_Keyword,
						Domain: "google",
						Ip:     [][]byte{{8, 8, 8, 8}},
					},
					{
						Type:   dns.DomainMatchingType_Regex,
						Domain: ".*\\.com",
						Ip:     [][]byte{{8, 8, 4, 4}},
					},
				},
				ClientIp: []byte{10, 0, 0, 1},
			},
		},
	})
}
