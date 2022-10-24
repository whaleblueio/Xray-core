package command

//go:generate go run github.com/whaleblueio/Xray-core/common/errors/errorgen

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/whaleblueio/Xray-core/common/protocol"
	"regexp"
	"runtime"
	"time"

	grpc "google.golang.org/grpc"

	"github.com/whaleblueio/Xray-core/app/stats"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/strmatcher"
	"github.com/whaleblueio/Xray-core/core"
	feature_stats "github.com/whaleblueio/Xray-core/features/stats"
)

// statsServer is an implementation of StatsService.
type statsServer struct {
	stats     feature_stats.Manager
	startTime time.Time
}

func NewStatsServer(manager feature_stats.Manager) StatsServiceServer {
	return &statsServer{
		stats:     manager,
		startTime: time.Now(),
	}
}

func (s *statsServer) GetStats(ctx context.Context, request *GetStatsRequest) (*GetStatsResponse, error) {
	c := s.stats.GetCounter(request.Name)
	if c == nil {
		return nil, newError(request.Name, " not found.")
	}
	var value int64
	if request.Reset_ {
		value = c.Set(0)
	} else {
		value = c.Value()
	}

	return &GetStatsResponse{
		Stat: &Stat{
			Name:  request.Name,
			Value: value,
		},
	}, nil
}

func (s *statsServer) QueryStats(ctx context.Context, request *QueryStatsRequest) (*QueryStatsResponse, error) {
	matcher, err := strmatcher.Substr.New(request.Pattern)
	if err != nil {
		return nil, err
	}
	response := &QueryStatsResponse{}

	manager, ok := s.stats.(*stats.Manager)
	if !ok {
		return nil, newError("QueryStats only works its own stats.Manager.")
	}
	reg := regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(v2rayService.com)?")

	manager.VisitCounters(func(name string, c feature_stats.Counter) bool {
		if matcher.Match(name) {
			var value int64
			if request.Reset_ {
				value = c.Set(0)
			} else {
				value = c.Value()
			}
			email := reg.FindString(name)
			ips := protocol.GetIPs(email)
			logrus.Debugf("QueryStats get user:%s connected ips:%s", email, len(ips))
			response.Stat = append(response.Stat, &Stat{
				Name:  name,
				Value: value,
				Ips:   ips,
			})
		}
		return true
	})

	return response, nil
}

func (s *statsServer) GetSysStats(ctx context.Context, request *SysStatsRequest) (*SysStatsResponse, error) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	uptime := time.Since(s.startTime)

	response := &SysStatsResponse{
		Uptime:       uint32(uptime.Seconds()),
		NumGoroutine: uint32(runtime.NumGoroutine()),
		Alloc:        rtm.Alloc,
		TotalAlloc:   rtm.TotalAlloc,
		Sys:          rtm.Sys,
		Mallocs:      rtm.Mallocs,
		Frees:        rtm.Frees,
		LiveObjects:  rtm.Mallocs - rtm.Frees,
		NumGC:        rtm.NumGC,
		PauseTotalNs: rtm.PauseTotalNs,
	}

	return response, nil
}

func (s *statsServer) mustEmbedUnimplementedStatsServiceServer() {}

type service struct {
	statsManager feature_stats.Manager
}

func (s *service) Register(server *grpc.Server) {
	ss := NewStatsServer(s.statsManager)
	RegisterStatsServiceServer(server, ss)

	// For compatibility purposes
	vCoreDesc := StatsService_ServiceDesc
	vCoreDesc.ServiceName = "v2ray.core.app.stats.command.StatsService"
	server.RegisterService(&vCoreDesc, ss)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := new(service)

		core.RequireFeatures(ctx, func(sm feature_stats.Manager) {
			s.statsManager = sm
		})

		return s, nil
	}))
}
