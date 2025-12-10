package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/cmd/interaction/dal"
	"github.com/nnieie/golanglab5/cmd/interaction/kafka"
	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/config"
	interaction "github.com/nnieie/golanglab5/kitex_gen/interaction/interactionservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

const (
	waitTimeForGracefulShutdown = 5 * time.Second
)

func Init() {
	logger.InitKlog()
	config.Init(constants.InteractionServiceName)
	dal.Init()
	rpc.InitVideoRPC()
	kafka.InitKafka()
}

func main() {
	shutdown, err := tracer.InitOpenTelemetry(constants.InteractionServiceName, constants.OpenTelemetryCollectorEndpoint)
	if err != nil {
		logger.Fatalf("init tracer failed: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			logger.Errorf("shutdown tracer failed: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	Init()
	go kafka.StartInteractionConsumer(ctx)
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

	go func() {
		err = svr.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Infof("Shutting down...")
	cancel()

	err = svr.Stop()
	if err != nil {
		logger.Errorf("Error stopping server: %v", err)
	}
	logger.Infof("Waiting for 5 seconds to gracefully shutdown...")
	time.Sleep(waitTimeForGracefulShutdown)
}
