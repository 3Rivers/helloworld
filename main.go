package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"github.com/3Rivers/helloworld/handler"
	"github.com/3Rivers/helloworld/subscriber"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"

	helloworld "github.com/3Rivers/helloworld/proto/helloworld"
)

var etcdReg registry.Registry

func  init()  {
	//新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
	etcdReg = etcd.NewRegistry(
		registry.Addrs("192.168.2.254:12379", "192.168.2.254:22379", "192.168.2.254:32379"),
	)
}

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.helloworld"),
		micro.Version("latest"),
		micro.Registry(etcdReg),
		micro.WrapHandler(limiter.NewHandlerWrapper(1)),//设置qps
	)

	// Initialise service
	service.Init()

	// Register Handler
	helloworld.RegisterHelloworldHandler(service.Server(), new(handler.Helloworld))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.helloworld", service.Server(), new(subscriber.Helloworld))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
