package main

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nnieie/golanglab5/cmd/video/dal"
	"github.com/nnieie/golanglab5/cmd/video/rpc"
	"github.com/nnieie/golanglab5/config"
	video "github.com/nnieie/golanglab5/kitex_gen/video/videoservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/oss"
	"github.com/nnieie/golanglab5/pkg/tracer"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func Init() {
	logger.InitKlog()
	config.Init(constants.VideoServiceName)
	oss.InitR2Client()
	dal.Init()
	rpc.InitUserRPC()
}

func main() {
	shutdown, err := tracer.InitOpenTelemetry(constants.VideoServiceName, constants.OpenTelemetryCollectorEndpoint)
	if err != nil {
		logger.Fatalf("init tracer failed: %v", err)
	}

	shutdownMetrics, err := tracer.InitMetrics(constants.VideoServiceName, constants.OpenTelemetryCollectorEndpoint)
	if err != nil {
		logger.Fatalf("init metrics failed: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			logger.Errorf("shutdown tracer failed: %v", err)
		}
		if err := shutdownMetrics(ctx); err != nil {
			logger.Errorf("shutdown metrics failed: %v", err)
		}
	}()
	Init()

	// 启动 pprof HTTP 服务器
	go func() {
		logger.Infof("Starting pprof server on :6062")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe("0.0.0.0:6062", nil); err != nil {
			logger.Errorf("pprof server failed: %v", err)
		}
	}()

	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("newEtcdRegistry err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", config.Service.Addr)
	if err != nil {
		logger.Fatalf("resolve etcd addr err: %v", err)
	}

	svcImpl := new(VideoServiceImpl)
	// TODO: Snowflake
	svcImpl.Snowflake, _ = utils.NewSnowflake(1)
	svr := video.NewServer(svcImpl,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Service.Name}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithSuite(kitextracing.NewServerSuite()),
	)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("etcd start err: %v", err)
	}
}
