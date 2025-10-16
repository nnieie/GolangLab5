package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/cmd/video/dal"
	"github.com/nnieie/golanglab5/cmd/video/rpc"
	"github.com/nnieie/golanglab5/config"
	video "github.com/nnieie/golanglab5/kitex_gen/video/videoservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/oss"
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
	Init()
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
	)

	err = svr.Run()

	if err != nil {
		logger.Fatalf("etcd start err: %v", err)
	}
}
