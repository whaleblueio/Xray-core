package quic_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/buf"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"github.com/whaleblueio/Xray-core/common/protocol/tls/cert"
	"github.com/whaleblueio/Xray-core/common/serial"
	"github.com/whaleblueio/Xray-core/testing/servers/udp"
	"github.com/whaleblueio/Xray-core/transport/internet"
	"github.com/whaleblueio/Xray-core/transport/internet/headers/wireguard"
	"github.com/whaleblueio/Xray-core/transport/internet/quic"
	"github.com/whaleblueio/Xray-core/transport/internet/tls"
)

func TestQuicConnection(t *testing.T) {
	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			Certificate: []*tls.Certificate{
				tls.ParseCertificate(
					cert.MustGenerate(nil,
						cert.DNSNames("www.example.com"),
					),
				),
			},
		},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				common.Must2(conn.Write(b.Bytes()))
			}
		}()
	})
	common.Must(err)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			ServerName:    "www.example.com",
			AllowInsecure: true,
		},
	})
	common.Must(err)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}
}

func TestQuicConnectionWithoutTLS(t *testing.T) {
	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				common.Must2(conn.Write(b.Bytes()))
			}
		}()
	})
	common.Must(err)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
	})
	common.Must(err)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}
}

func TestQuicConnectionAuthHeader(t *testing.T) {
	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName: "quic",
		ProtocolSettings: &quic.Config{
			Header: serial.ToTypedMessage(&wireguard.WireguardConfig{}),
			Key:    "abcd",
			Security: &protocol.SecurityConfig{
				Type: protocol.SecurityType_AES128_GCM,
			},
		},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				common.Must2(conn.Write(b.Bytes()))
			}
		}()
	})
	common.Must(err)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName: "quic",
		ProtocolSettings: &quic.Config{
			Header: serial.ToTypedMessage(&wireguard.WireguardConfig{}),
			Key:    "abcd",
			Security: &protocol.SecurityConfig{
				Type: protocol.SecurityType_AES128_GCM,
			},
		},
	})
	common.Must(err)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}

	common.Must2(conn.Write(b1))

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	if r := cmp.Diff(b2.Bytes(), b1); r != "" {
		t.Error(r)
	}
}
