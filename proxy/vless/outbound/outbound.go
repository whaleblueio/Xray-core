package outbound

//go:generate go run github.com/whaleblueio/Xray-core/common/errors/errorgen

import (
	"context"
	"syscall"
	"time"

	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/buf"
	"github.com/whaleblueio/Xray-core/common/net"
	"github.com/whaleblueio/Xray-core/common/platform"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"github.com/whaleblueio/Xray-core/common/retry"
	"github.com/whaleblueio/Xray-core/common/session"
	"github.com/whaleblueio/Xray-core/common/signal"
	"github.com/whaleblueio/Xray-core/common/task"
	"github.com/whaleblueio/Xray-core/common/xudp"
	core "github.com/whaleblueio/Xray-core/core"
	"github.com/whaleblueio/Xray-core/features/policy"
	"github.com/whaleblueio/Xray-core/features/stats"
	"github.com/whaleblueio/Xray-core/proxy/vless"
	"github.com/whaleblueio/Xray-core/proxy/vless/encoding"
	"github.com/whaleblueio/Xray-core/transport"
	"github.com/whaleblueio/Xray-core/transport/internet"
	"github.com/whaleblueio/Xray-core/transport/internet/xtls"
)

var (
	xtls_show = false
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))

	const defaultFlagValue = "NOT_DEFINED_AT_ALL"

	xtlsShow := platform.NewEnvFlag("xray.vless.xtls.show").GetValue(func() string { return defaultFlagValue })
	if xtlsShow == "true" {
		xtls_show = true
	}
}

// Handler is an outbound connection handler for VLess protocol.
type Handler struct {
	serverList    *protocol.ServerList
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
	cone          bool
}

// New creates a new VLess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Vnext {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, newError("failed to parse server spec").Base(err).AtError()
		}
		serverList.AddServer(s)
	}

	v := core.MustFromContext(ctx)
	handler := &Handler{
		serverList:    serverList,
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		cone:          ctx.Value("cone").(bool),
	}

	return handler, nil
}

// Process implements proxy.Outbound.Process().
func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	var rec *protocol.ServerSpec
	var conn internet.Connection

	if err := retry.ExponentialBackoff(5, 200).On(func() error {
		rec = h.serverPicker.PickServer()
		var err error
		conn, err = dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return newError("failed to find an available destination").Base(err).AtWarning()
	}
	defer conn.Close()

	iConn := conn
	statConn, ok := iConn.(*internet.StatCouterConnection)
	if ok {
		iConn = statConn.Connection
	}

	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("target not specified").AtError()
	}

	target := outbound.Target
	newError("tunneling request to ", target, " via ", rec.Destination()).AtInfo().WriteToLog(session.ExportIDToError(ctx))

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	if target.Address.Family().IsDomain() && target.Address.Domain() == "v1.mux.cool" {
		command = protocol.RequestCommandMux
	}

	request := &protocol.RequestHeader{
		Version: encoding.Version,
		User:    rec.PickUser(),
		Command: command,
		Address: target.Address,
		Port:    target.Port,
	}

	account := request.User.Account.(*vless.MemoryAccount)

	requestAddons := &encoding.Addons{
		Flow: account.Flow,
	}

	var rawConn syscall.RawConn
	var sctx context.Context

	allowUDP443 := false
	switch requestAddons.Flow {
	case vless.XRO + "-udp443", vless.XRD + "-udp443", vless.XRS + "-udp443":
		allowUDP443 = true
		requestAddons.Flow = requestAddons.Flow[:16]
		fallthrough
	case vless.XRO, vless.XRD, vless.XRS:
		switch request.Command {
		case protocol.RequestCommandMux:
			return newError(requestAddons.Flow + " doesn't support Mux").AtWarning()
		case protocol.RequestCommandUDP:
			if !allowUDP443 && request.Port == 443 {
				return newError(requestAddons.Flow + " stopped UDP/443").AtInfo()
			}
			requestAddons.Flow = ""
		case protocol.RequestCommandTCP:
			if xtlsConn, ok := iConn.(*xtls.Conn); ok {
				xtlsConn.RPRX = true
				xtlsConn.SHOW = xtls_show
				xtlsConn.MARK = "XTLS"
				if requestAddons.Flow == vless.XRS {
					sctx = ctx
					requestAddons.Flow = vless.XRD
				}
				if requestAddons.Flow == vless.XRD {
					xtlsConn.DirectMode = true
					if sc, ok := xtlsConn.Connection.(syscall.Conn); ok {
						rawConn, _ = sc.SyscallConn()
					}
				}
			} else {
				return newError(`failed to use ` + requestAddons.Flow + `, maybe "security" is not "xtls"`).AtWarning()
			}
		}
	default:
		if _, ok := iConn.(*xtls.Conn); ok {
			panic(`To avoid misunderstanding, you must fill in VLESS "flow" when using XTLS.`)
		}
	}

	sessionPolicy := h.policyManager.ForLevel(request.User.Level)
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)

	clientReader := link.Reader // .(*pipe.Reader)
	clientWriter := link.Writer // .(*pipe.Writer)

	if request.Command == protocol.RequestCommandUDP && h.cone && request.Port != 53 && request.Port != 443 {
		request.Command = protocol.RequestCommandMux
		request.Address = net.DomainAddress("v1.mux.cool")
		request.Port = net.Port(666)
	}
	//user := request.User
	//var bucket *rateLimit.Bucket
	//if user != nil {
	//	bucket = protocol.GetBucket(user.Email)
	//} else {
	//	newError("user is nil").WriteToLog()
	//}
	postRequest := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		bufferWriter := buf.NewBufferedWriter(buf.NewWriter(conn))
		if err := encoding.EncodeRequestHeader(bufferWriter, request, requestAddons); err != nil {
			return newError("failed to encode request header").Base(err).AtWarning()
		}

		// default: serverWriter := bufferWriter
		serverWriter := encoding.EncodeBodyAddons(bufferWriter, request, requestAddons)
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			serverWriter = xudp.NewPacketWriter(serverWriter, target)
		}
		if err := buf.CopyOnceTimeoutWithLimiter(clientReader, serverWriter, time.Millisecond*100, nil); err != nil && err != buf.ErrNotTimeoutReader && err != buf.ErrReadTimeout {
			return err // ...
		}

		// Flush; bufferWriter.WriteMultiBufer now is bufferWriter.writer.WriteMultiBuffer
		if err := bufferWriter.SetBuffered(false); err != nil {
			return newError("failed to write A request payload").Base(err).AtWarning()
		}

		// from clientReader.ReadMultiBuffer to serverWriter.WriteMultiBufer
		if err := buf.CopyWithLimiter(clientReader, serverWriter, nil, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transfer request payload").Base(err).AtInfo()
		}

		// Indicates the end of request payload.
		switch requestAddons.Flow {
		default:
		}
		return nil
	}

	getResponse := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		responseAddons, err := encoding.DecodeResponseHeader(conn, request)
		if err != nil {
			return newError("failed to decode response header").Base(err).AtInfo()
		}

		// default: serverReader := buf.NewReader(conn)
		serverReader := encoding.DecodeBodyAddons(conn, request, responseAddons)
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			serverReader = xudp.NewPacketReader(conn)
		}

		if rawConn != nil {
			var counter stats.Counter
			if statConn != nil {
				counter = statConn.ReadCounter
			}
			err = encoding.ReadV(serverReader, clientWriter, timer, iConn.(*xtls.Conn), rawConn, counter, sctx)
		} else {
			// from serverReader.ReadMultiBuffer to clientWriter.WriteMultiBufer
			err = buf.CopyWithLimiter(serverReader, clientWriter, nil, buf.UpdateActivity(timer))
		}

		if err != nil {
			return newError("failed to transfer response payload").Base(err).AtInfo()
		}

		return nil
	}

	if err := task.Run(ctx, postRequest, task.OnSuccess(getResponse, task.Close(clientWriter))); err != nil {
		return newError("connection ends").Base(err).AtInfo()
	}

	return nil
}
