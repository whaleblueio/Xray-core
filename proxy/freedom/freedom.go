package freedom

//go:generate go run github.com/whaleblueio/Xray-core/common/errors/errorgen

import (
	"context"
	rateLimit "github.com/juju/ratelimit"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"time"

	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/buf"
	"github.com/whaleblueio/Xray-core/common/dice"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/retry"
	"github.com/whaleblueio/Xray-core/common/session"
	"github.com/whaleblueio/Xray-core/common/signal"
	"github.com/whaleblueio/Xray-core/common/task"
	"github.com/whaleblueio/Xray-core/core"
	"github.com/whaleblueio/Xray-core/features/dns"
	"github.com/whaleblueio/Xray-core/features/policy"
	"github.com/whaleblueio/Xray-core/features/stats"
	"github.com/whaleblueio/Xray-core/transport"
	"github.com/whaleblueio/Xray-core/transport/internet"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		h := new(Handler)
		if err := core.RequireFeatures(ctx, func(pm policy.Manager, d dns.Client) error {
			return h.Init(config.(*Config), pm, d)
		}); err != nil {
			return nil, err
		}
		return h, nil
	}))
}

// Handler handles Freedom connections.
type Handler struct {
	policyManager policy.Manager
	dns           dns.Client
	config        *Config
}

// Init initializes the Handler with necessary parameters.
func (h *Handler) Init(config *Config, pm policy.Manager, d dns.Client) error {
	h.config = config
	h.policyManager = pm
	h.dns = d

	return nil
}

func (h *Handler) policy() policy.Session {
	p := h.policyManager.ForLevel(h.config.UserLevel)
	if h.config.Timeout > 0 && h.config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(h.config.Timeout) * time.Second
	}
	return p
}

func (h *Handler) resolveIP(ctx context.Context, domain string, localAddr net.Address) net.Address {
	var option dns.IPOption = dns.IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
		FakeEnable: false,
	}
	if h.config.DomainStrategy == Config_USE_IP4 || (localAddr != nil && localAddr.Family().IsIPv4()) {
		option = dns.IPOption{
			IPv4Enable: true,
			IPv6Enable: false,
			FakeEnable: false,
		}
	} else if h.config.DomainStrategy == Config_USE_IP6 || (localAddr != nil && localAddr.Family().IsIPv6()) {
		option = dns.IPOption{
			IPv4Enable: false,
			IPv6Enable: true,
			FakeEnable: false,
		}
	}

	ips, err := h.dns.LookupIP(domain, option)
	if err != nil {
		newError("failed to get IP address for domain ", domain).Base(err).WriteToLog(session.ExportIDToError(ctx))
	}
	if len(ips) == 0 {
		return nil
	}
	return net.IPAddress(ips[dice.Roll(len(ips))])
}

func isValidAddress(addr *net.IPOrDomain) bool {
	if addr == nil {
		return false
	}

	a := addr.AsAddress()
	return a != net.AnyIP
}

// Process implements proxy.Outbound.
func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	outbound := session.OutboundFromContext(ctx)

	if outbound == nil || !outbound.Target.IsValid() {
		return newError("target not specified.")
	}
	destination := outbound.Target
	UDPOverride := net.UDPDestination(nil, 0)
	if h.config.DestinationOverride != nil {
		server := h.config.DestinationOverride.Server
		if isValidAddress(server.Address) {
			destination.Address = server.Address.AsAddress()
			UDPOverride.Address = destination.Address
		}
		if server.Port != 0 {
			destination.Port = net.Port(server.Port)
			UDPOverride.Port = destination.Port
		}
	}
	inboundSession := session.InboundFromContext(ctx)
	newError(" from:", inboundSession.Source.Address.IP().String(), " opening connection to ", destination, " sequeceId:", common.GetSequenceId()).WriteToLog(session.ExportIDToError(ctx))

	input := link.Reader
	output := link.Writer
	var conn internet.Connection
	err := retry.ExponentialBackoff(5, 100).On(func() error {
		dialDest := destination
		if h.config.useIP() && dialDest.Address.Family().IsDomain() {
			ip := h.resolveIP(ctx, dialDest.Address.Domain(), dialer.Address())
			if ip != nil {
				dialDest = net.Destination{
					Network: dialDest.Network,
					Address: ip,
					Port:    dialDest.Port,
				}
				newError("dialing to ", dialDest, " sequenceId:", common.GetSequenceId()).WriteToLog(session.ExportIDToError(ctx))
			}
		}

		rawConn, err := dialer.Dial(ctx, dialDest)
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		return newError("failed to open connection to ", destination).Base(err)
	}

	defer conn.Close()
	plcy := h.policy()
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, plcy.Timeouts.ConnectionIdle)

	var user *protocol.MemoryUser

	if inboundSession != nil || inboundSession.User != nil {
		user = inboundSession.User
	}
	var bucket *rateLimit.Bucket
	if user != nil {
		//user.IpCounter.Add(string(dialer.Address().IP()))
		bucket = protocol.GetBucket(user.Email)
	} else {
		newError("user is nil").WriteToLog()
	}
	requestDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.DownlinkOnly)

		var writer buf.Writer
		if destination.Network == net.Network_TCP {
			writer = buf.NewWriter(conn)
		} else {
			writer = NewPacketWriter(conn, h, ctx, UDPOverride)
		}

		if err := buf.CopyWithLimiter(input, writer, bucket, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process request").Base(err)
		}
		user.IpCounter.Add(inboundSession.Source.Address.IP().String())
		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.UplinkOnly)

		var reader buf.Reader
		if destination.Network == net.Network_TCP {
			reader = buf.NewReader(conn)
		} else {
			reader = NewPacketReader(conn, UDPOverride)
		}
		if err := buf.CopyWithLimiter(reader, output, bucket, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process response").Base(err)
		}
		user.IpCounter.Add(inboundSession.Source.Address.IP().String())
		return nil
	}

	if err := task.Run(ctx, requestDone, task.OnSuccess(responseDone, task.Close(output))); err != nil {
		return newError("connection ends,sequenceId:", common.GetSequenceId()).Base(err)
	}

	return nil
}

func NewPacketReader(conn net.Conn, UDPOverride net.Destination) buf.Reader {
	iConn := conn
	statConn, ok := iConn.(*internet.StatCouterConnection)
	if ok {
		iConn = statConn.Connection
	}
	var counter stats.Counter
	if statConn != nil {
		counter = statConn.ReadCounter
	}
	if c, ok := iConn.(*internet.PacketConnWrapper); ok && UDPOverride.Address == nil && UDPOverride.Port == 0 {
		return &PacketReader{
			PacketConnWrapper: c,
			Counter:           counter,
		}
	}
	return &buf.PacketReader{Reader: conn}
}

func NewPacketReaderWithRateLimiter(conn net.Conn, UDPOverride net.Destination, speed int64) buf.Reader {
	iConn := conn
	statConn, ok := iConn.(*internet.StatCouterConnection)
	if ok {
		iConn = statConn.Connection
	}
	var counter stats.Counter
	if statConn != nil {
		counter = statConn.ReadCounter
	}
	var bucket *rateLimit.Bucket
	if speed > 0 {
		bucket = rateLimit.NewBucketWithQuantum(time.Second, speed, speed)
	}
	if c, ok := iConn.(*internet.PacketConnWrapper); ok && UDPOverride.Address == nil && UDPOverride.Port == 0 {
		return &PacketReader{
			PacketConnWrapper: c,
			Counter:           counter,
			Bucket:            bucket,
		}
	}
	//return &buf.PacketReader{Reader: conn}
	return buf.NewPacketReaderWithRateLimiter(conn, speed)
}

type PacketReader struct {
	*internet.PacketConnWrapper
	stats.Counter
	Bucket *rateLimit.Bucket
}

func (r *PacketReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b := buf.New()
	b.Resize(0, buf.Size)
	n, d, err := r.PacketConnWrapper.ReadFrom(b.Bytes())
	if err != nil {
		b.Release()
		return nil, err
	}
	b.Resize(0, int32(n))
	b.UDP = &net.Destination{
		Address: net.IPAddress(d.(*net.UDPAddr).IP),
		Port:    net.Port(d.(*net.UDPAddr).Port),
		Network: net.Network_UDP,
	}
	if r.Bucket != nil {
		r.Bucket.Wait(int64(n))
	}
	if r.Counter != nil {
		r.Counter.Add(int64(n))
	}
	return buf.MultiBuffer{b}, nil
}

func NewPacketWriter(conn net.Conn, h *Handler, ctx context.Context, UDPOverride net.Destination) buf.Writer {
	iConn := conn
	statConn, ok := iConn.(*internet.StatCouterConnection)
	if ok {
		iConn = statConn.Connection
	}
	var counter stats.Counter
	if statConn != nil {
		counter = statConn.WriteCounter
	}
	if c, ok := iConn.(*internet.PacketConnWrapper); ok {
		return &PacketWriter{
			PacketConnWrapper: c,
			Counter:           counter,
			Handler:           h,
			Context:           ctx,
			UDPOverride:       UDPOverride,
		}
	}
	return &buf.SequentialWriter{Writer: conn}
}
func NewPacketWriterWithRateLimiter(conn net.Conn, h *Handler, ctx context.Context, UDPOverride net.Destination, speed int64) buf.Writer {
	iConn := conn
	statConn, ok := iConn.(*internet.StatCouterConnection)
	if ok {
		iConn = statConn.Connection
	}
	var counter stats.Counter
	if statConn != nil {
		counter = statConn.WriteCounter
	}
	var bucket *rateLimit.Bucket
	if speed > 0 {
		bucket = rateLimit.NewBucketWithQuantum(time.Second, speed, speed)
	}

	if c, ok := iConn.(*internet.PacketConnWrapper); ok {
		return &PacketWriter{
			PacketConnWrapper: c,
			Counter:           counter,
			Handler:           h,
			Context:           ctx,
			UDPOverride:       UDPOverride,
			Bucket:            bucket,
		}
	}
	return &buf.SequentialWriter{Writer: conn, Bucket: bucket}
}

type PacketWriter struct {
	*internet.PacketConnWrapper
	stats.Counter
	*Handler
	context.Context
	UDPOverride net.Destination
	Bucket      *rateLimit.Bucket
}

func (w *PacketWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	for {
		mb2, b := buf.SplitFirst(mb)
		mb = mb2
		if b == nil {
			break
		}
		var n int
		var err error
		if b.UDP != nil {
			if w.UDPOverride.Address != nil {
				b.UDP.Address = w.UDPOverride.Address
			}
			if w.UDPOverride.Port != 0 {
				b.UDP.Port = w.UDPOverride.Port
			}
			if w.Handler.config.useIP() && b.UDP.Address.Family().IsDomain() {
				ip := w.Handler.resolveIP(w.Context, b.UDP.Address.Domain(), nil)
				if ip != nil {
					b.UDP.Address = ip
				}
			}
			destAddr, _ := net.ResolveUDPAddr("udp", b.UDP.NetAddr())
			if destAddr == nil {
				b.Release()
				continue
			}
			n, err = w.PacketConnWrapper.WriteTo(b.Bytes(), destAddr)
		} else {
			n, err = w.PacketConnWrapper.Write(b.Bytes())
		}
		b.Release()
		if err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
		if w.Bucket != nil {
			w.Bucket.Wait(int64(n))

		}
		if w.Counter != nil {
			w.Counter.Add(int64(n))
		}
	}
	return nil
}
