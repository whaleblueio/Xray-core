package tcp

import (
	"context"

	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/session"
	"github.com/whaleblueio/Xray-core/transport/internet"
	"github.com/whaleblueio/Xray-core/transport/internet/tls"
	"github.com/whaleblueio/Xray-core/transport/internet/xtls"
)

// Dial dials a new TCP connection to the given destination.
func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	newError("dialing TCP to ", dest).WriteToLog(session.ExportIDToError(ctx))
	conn, err := internet.DialSystem(ctx, dest, streamSettings.SocketSettings)
	if err != nil {
		return nil, err
	}

	if config := tls.ConfigFromStreamSettings(streamSettings); config != nil {
		tlsConfig := config.GetTLSConfig(tls.WithDestination(dest))
		if fingerprint, ok := tls.Fingerprints[config.Fingerprint]; ok {
			conn = tls.UClient(conn, tlsConfig, fingerprint)
			if err := conn.(*tls.UConn).Handshake(); err != nil {
				return nil, err
			}
		} else {
			conn = tls.Client(conn, tlsConfig)
		}
	} else if config := xtls.ConfigFromStreamSettings(streamSettings); config != nil {
		xtlsConfig := config.GetXTLSConfig(xtls.WithDestination(dest))
		conn = xtls.Client(conn, xtlsConfig)
	}

	tcpSettings := streamSettings.ProtocolSettings.(*Config)
	if tcpSettings.HeaderSettings != nil {
		headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
		if err != nil {
			return nil, newError("failed to get header settings").Base(err).AtError()
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, newError("failed to create header authenticator").Base(err).AtError()
		}
		conn = auth.Client(conn)
	}
	return internet.Connection(conn), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
