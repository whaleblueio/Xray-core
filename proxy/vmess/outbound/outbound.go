package outbound

//go:generate go run github.com/whaleblueio/Xray-core/common/errors/errorgen

import (
	"context"
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
	"github.com/whaleblueio/Xray-core/proxy/vmess"
	"github.com/whaleblueio/Xray-core/proxy/vmess/encoding"
	"github.com/whaleblueio/Xray-core/transport"
	"github.com/whaleblueio/Xray-core/transport/internet"
)

// Handler is an outbound connection handler for VMess protocol.
type Handler struct {
	serverList    *protocol.ServerList
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
	cone          bool
}

// New creates a new VMess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Receiver {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, newError("failed to parse server spec").Base(err)
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

	err := retry.ExponentialBackoff(5, 200).On(func() error {
		rec = h.serverPicker.PickServer()
		rawConn, err := dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return newError("failed to find an available destination").Base(err).AtWarning()
	}
	defer conn.Close()

	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("target not specified").AtError()
	}

	target := outbound.Target
	newError("tunneling request to ", target, " via ", rec.Destination()).WriteToLog(session.ExportIDToError(ctx))

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	if target.Address.Family().IsDomain() && target.Address.Domain() == "v1.mux.cool" {
		command = protocol.RequestCommandMux
	}

	user := rec.PickUser()
	request := &protocol.RequestHeader{
		Version: encoding.Version,
		User:    user,
		Command: command,
		Address: target.Address,
		Port:    target.Port,
		Option:  protocol.RequestOptionChunkStream,
	}

	account := request.User.Account.(*vmess.MemoryAccount)
	request.Security = account.Security

	if request.Security == protocol.SecurityType_AES128_GCM || request.Security == protocol.SecurityType_NONE || request.Security == protocol.SecurityType_CHACHA20_POLY1305 {
		request.Option.Set(protocol.RequestOptionChunkMasking)
	}

	if shouldEnablePadding(request.Security) && request.Option.Has(protocol.RequestOptionChunkMasking) {
		request.Option.Set(protocol.RequestOptionGlobalPadding)
	}

	if request.Security == protocol.SecurityType_ZERO {
		request.Security = protocol.SecurityType_NONE
		request.Option.Clear(protocol.RequestOptionChunkStream)
		request.Option.Clear(protocol.RequestOptionChunkMasking)
	}

	input := link.Reader
	output := link.Writer

	isAEAD := false
	if !aeadDisabled && len(account.AlterIDs) == 0 {
		isAEAD = true
	}
	//var bucket *rateLimit.Bucket
	//if user != nil {
	//	bucket = protocol.GetBucket(user.Email)
	//}

	session := encoding.NewClientSession(ctx, isAEAD, protocol.DefaultIDHash)
	sessionPolicy := h.policyManager.ForLevel(request.User.Level)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)

	if request.Command == protocol.RequestCommandUDP && h.cone && request.Port != 53 && request.Port != 443 {
		request.Command = protocol.RequestCommandMux
		request.Address = net.DomainAddress("v1.mux.cool")
		request.Port = net.Port(666)
	}

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		writer := buf.NewBufferedWriter(buf.NewWriter(conn))
		if err := session.EncodeRequestHeader(request, writer); err != nil {
			return newError("failed to encode request").Base(err).AtWarning()
		}

		bodyWriter := session.EncodeRequestBody(request, writer)
		bodyWriter2 := bodyWriter
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			bodyWriter = xudp.NewPacketWriter(bodyWriter, target)
		}
		if err := buf.CopyOnceTimeoutWithLimiter(input, bodyWriter, time.Millisecond*100, nil); err != nil && err != buf.ErrNotTimeoutReader && err != buf.ErrReadTimeout {
			return newError("failed to write first payload").Base(err)
		}

		if err := writer.SetBuffered(false); err != nil {
			return err
		}

		if err := buf.CopyWithLimiter(input, bodyWriter, nil, buf.UpdateActivity(timer)); err != nil {
			return err
		}

		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if err := bodyWriter2.WriteMultiBuffer(buf.MultiBuffer{}); err != nil {
				return err
			}
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		reader := &buf.BufferedReader{Reader: buf.NewReader(conn)}
		header, err := session.DecodeResponseHeader(reader)
		if err != nil {
			return newError("failed to read header").Base(err)
		}
		h.handleCommand(rec.Destination(), header.Command)

		bodyReader := session.DecodeResponseBody(request, reader)
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			bodyReader = xudp.NewPacketReader(&buf.BufferedReader{Reader: bodyReader})
		}

		return buf.CopyWithLimiter(bodyReader, output, nil, buf.UpdateActivity(timer))
	}

	var responseDonePost = task.OnSuccess(responseDone, task.Close(output))
	if err := task.Run(ctx, requestDone, responseDonePost); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

var (
	enablePadding = false
	aeadDisabled  = false
)

func shouldEnablePadding(s protocol.SecurityType) bool {
	return enablePadding || s == protocol.SecurityType_AES128_GCM || s == protocol.SecurityType_CHACHA20_POLY1305 || s == protocol.SecurityType_AUTO
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))

	const defaultFlagValue = "NOT_DEFINED_AT_ALL"

	paddingValue := platform.NewEnvFlag("xray.vmess.padding").GetValue(func() string { return defaultFlagValue })
	if paddingValue != defaultFlagValue {
		enablePadding = true
	}

	isAeadDisabled := platform.NewEnvFlag("xray.vmess.aead.disabled").GetValue(func() string { return defaultFlagValue })
	if isAeadDisabled == "true" {
		aeadDisabled = true
	}
}
