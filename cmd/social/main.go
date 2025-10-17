package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/cmd/social/dal"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/config"
	social "github.com/nnieie/golanglab5/kitex_gen/social/socialservice"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func Init() {
	logger.InitKlog()
	config.Init(constants.SocialServiceName)
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

	svcImpl := new(SocialServiceImpl)
	svr := social.NewServer(svcImpl,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.Service.Name}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
