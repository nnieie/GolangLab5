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

	"github.com/nnieie/golanglab5/cmd/interaction/dal"
	"github.com/nnieie/golanglab5/cmd/interaction/kafka"
	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/config"
	interaction "github.com/nnieie/golanglab5/kitex_gen/interaction/interactionservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

func Init() {
	logger.InitKlog()
	config.Init(constants.InteractionServiceName)
	dal.Init()
	rpc.InitVideoRPC()
	kafka.InitKafka()
	go kafka.ConsumeLikeEvent()
}

func main() {
	shutdown, err := tracer.InitOpenTelemetry(constants.InteractionServiceName, constants.OpenTelemetryCollectorEndpoint)
	if err != nil {
		logger.Fatalf("init tracer failed: %v", err)
	}

	shutdownMetrics, err := tracer.InitMetrics(constants.InteractionServiceName, constants.OpenTelemetryCollectorEndpoint)
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
	defer kafka.CloseKafka()

	// 启动 pprof HTTP 服务器
	go func() {
		logger.Infof("Starting pprof server on :6065")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe("0.0.0.0:6065", nil); err != nil {
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

	svcImpl := new(InteractionServiceImpl)
	svr := interaction.NewServer(svcImpl,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Service.Name}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithSuite(kitextracing.NewServerSuite()),
	)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("interaction service run err: %v", err)
	}
}
