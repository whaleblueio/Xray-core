package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Required features. Can't remove unless there is replacements.
	_ "github.com/whaleblueio/Xray-core/app/dispatcher"
	_ "github.com/whaleblueio/Xray-core/app/proxyman/inbound"
	_ "github.com/whaleblueio/Xray-core/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/whaleblueio/Xray-core/app/commander"
	_ "github.com/whaleblueio/Xray-core/app/log/command"
	_ "github.com/whaleblueio/Xray-core/app/proxyman/command"
	_ "github.com/whaleblueio/Xray-core/app/stats/command"

	// Other optional features.
	_ "github.com/whaleblueio/Xray-core/app/dns"
	_ "github.com/whaleblueio/Xray-core/app/dns/fakedns"
	_ "github.com/whaleblueio/Xray-core/app/log"
	_ "github.com/whaleblueio/Xray-core/app/policy"
	_ "github.com/whaleblueio/Xray-core/app/reverse"
	_ "github.com/whaleblueio/Xray-core/app/router"
	_ "github.com/whaleblueio/Xray-core/app/stats"

	// Inbound and outbound proxies.
	_ "github.com/whaleblueio/Xray-core/proxy/blackhole"
	_ "github.com/whaleblueio/Xray-core/proxy/dns"
	_ "github.com/whaleblueio/Xray-core/proxy/dokodemo"
	_ "github.com/whaleblueio/Xray-core/proxy/freedom"
	_ "github.com/whaleblueio/Xray-core/proxy/http"
	_ "github.com/whaleblueio/Xray-core/proxy/mtproto"
	_ "github.com/whaleblueio/Xray-core/proxy/shadowsocks"
	_ "github.com/whaleblueio/Xray-core/proxy/socks"
	_ "github.com/whaleblueio/Xray-core/proxy/trojan"
	_ "github.com/whaleblueio/Xray-core/proxy/vless/inbound"
	_ "github.com/whaleblueio/Xray-core/proxy/vless/outbound"
	_ "github.com/whaleblueio/Xray-core/proxy/vmess/inbound"
	_ "github.com/whaleblueio/Xray-core/proxy/vmess/outbound"

	// Transports
	_ "github.com/whaleblueio/Xray-core/transport/internet/domainsocket"
	_ "github.com/whaleblueio/Xray-core/transport/internet/grpc"
	_ "github.com/whaleblueio/Xray-core/transport/internet/http"
	_ "github.com/whaleblueio/Xray-core/transport/internet/kcp"
	_ "github.com/whaleblueio/Xray-core/transport/internet/quic"
	_ "github.com/whaleblueio/Xray-core/transport/internet/tcp"
	_ "github.com/whaleblueio/Xray-core/transport/internet/tls"
	_ "github.com/whaleblueio/Xray-core/transport/internet/udp"
	_ "github.com/whaleblueio/Xray-core/transport/internet/websocket"
	_ "github.com/whaleblueio/Xray-core/transport/internet/xtls"

	// Transport headers
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/http"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/noop"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/srtp"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/tls"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/utp"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/wechat"
	_ "github.com/whaleblueio/Xray-core/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/whaleblueio/Xray-core/main/json"
	_ "github.com/whaleblueio/Xray-core/main/toml"
	_ "github.com/whaleblueio/Xray-core/main/yaml"

	// Load config from file or http(s)
	_ "github.com/whaleblueio/Xray-core/main/confloader/external"

	// Commands
	_ "github.com/whaleblueio/Xray-core/main/commands/all"
)
