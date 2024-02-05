package grpc

import (
	"github.com/kdrkrgz/ecomm-micro/payment/internal/ports"
	"github.com/kdrkrgz/ecomm-proto/golang/payment"
	"google.golang.org/grpc"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	payment.UnimplementedPaymentServiceServer
}
