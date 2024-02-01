package grpc

import (
	"fmt"
	"github.com/kdrkrgz/ecomm-micro/order/config"
	"github.com/kdrkrgz/ecomm-micro/order/internal/ports"
	"github.com/kdrkrgz/ecomm-proto/golang/order"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	order.UnimplementedOrderServiceServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d : %v", a.port, err)
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	a.server = grpcServer
	order.RegisterOrderServiceServer(grpcServer, a)
	if config.GetEnv() == "dev" {
		reflection.Register(grpcServer)
	}
	log.Printf("grpc server is listening on port %d", a.port)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve grpc server on port %d : %v", a.port, err)
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}
