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

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/internal/user/dal"
	user "github.com/nnieie/golanglab5/kitex_gen/user/userservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/oss"
	"github.com/nnieie/golanglab5/pkg/tracer"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func initDependencies() {
	oss.InitR2Client()
	dal.Init()
}

func main() {
	logger.InitKlog()
	config.Init(constants.UserServiceName)

	// 初始化 OpenTelemetry
	shutdown, err := tracer.InitOpenTelemetry(constants.UserServiceName, config.TelemetryEndpoint())
	if err != nil {
		logger.Fatalf("init tracer failed: %v", err)
	}

	shutdownMetrics, err := tracer.InitMetrics(constants.UserServiceName, config.TelemetryEndpoint())
	if err != nil {
		logger.Fatalf("init metrics failed: %v", err)
	}

	// 程序退出前把最后的 Trace 数据发出去
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

	initDependencies()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("newEtcdRegistry err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", config.Service.Addr)
	if err != nil {
		logger.Fatalf("resolve etcd addr err: %v", err)
	}

	// 启动 pprof HTTP 服务器
	go func() {
		logger.Infof("Starting pprof server on :6061")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe("0.0.0.0:6061", nil); err != nil {
			logger.Errorf("pprof server failed: %v", err)
		}
	}()

	svcImpl := new(UserServiceImpl)
	// TODO: Snowflake
	svcImpl.Snowflake, _ = utils.NewSnowflake(1)
	options := []server.Option{
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Service.Name}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
	}
	if config.TraceEnabled() {
		// 注入 Kitex OTel 服务端套件 解析请求头里的 TraceID，把 api 和 user 连起来
		options = append(options, server.WithSuite(kitextracing.NewServerSuite()))
	}
	svr := user.NewServer(svcImpl, options...)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("etcd start err: %v", err)
	}
}
