package main

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/cmd/user/dal"
	"github.com/nnieie/golanglab5/config"
	user "github.com/nnieie/golanglab5/kitex_gen/user/userservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/oss"
	"github.com/nnieie/golanglab5/pkg/tracer"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func Init() {
	logger.InitKlog()
	config.Init(constants.UserServiceName)
	oss.InitR2Client()
	dal.Init()
}

func main() {
	// 初始化 OpenTelemetry
	shutdown, err := tracer.InitOpenTelemetry(constants.UserServiceName, constants.OpenTelemetryCollectorEndpoint)
	if err != nil {
		logger.Fatalf("init tracer failed: %v", err)
	}

	// 程序退出前把最后的 Trace 数据发出去
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			logger.Errorf("shutdown tracer failed: %v", err)
		}
	}()

	Init()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("newEtcdRegistry err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", config.Service.Addr)
	if err != nil {
		logger.Fatalf("resolve etcd addr err: %v", err)
	}

	svcImpl := new(UserServiceImpl)
	// TODO: Snowflake
	svcImpl.Snowflake, _ = utils.NewSnowflake(1)
	svr := user.NewServer(svcImpl,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Service.Name}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),

		// 注入 Kitex OTel 服务端套件 解析请求头里的 TraceID，把 api 和 user 连起来
		server.WithSuite(kitextracing.NewServerSuite()),
	)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("etcd start err: %v", err)
	}
}
