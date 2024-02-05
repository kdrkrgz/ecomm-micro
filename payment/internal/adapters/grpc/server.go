package grpc

import (
	"fmt"
	"github.com/kdrkrgz/ecomm-micro/payment/config"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/ports"
	"github.com/kdrkrgz/ecomm-proto/golang/payment"
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
	payment.UnimplementedPaymentServiceServer
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
		panic(fmt.Sprintf("failed to listen on port %d : %v", a.port, err))
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	a.server = grpcServer
	payment.RegisterPaymentServiceServer(grpcServer, a)
	if config.GetEnv() == "dev" {
		reflection.Register(grpcServer)
	}
	log.Println("starting payment service on port", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		panic(fmt.Sprintf("failed to serve grpc server on port %d : %v", a.port, err))
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}
