// +build !linux,!freebsd

package tcp

import (
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}
