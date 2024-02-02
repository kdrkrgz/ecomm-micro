package payment

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
	"github.com/kdrkrgz/ecomm-proto/golang/payment"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type Adapter struct {
	payment payment.PaymentServiceClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := payment.NewPaymentServiceClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(ctx context.Context, order *domain.Order) error {
	_, err := a.payment.CreatePayment(ctx, &payment.CreatePaymentRequest{
		OrderId:  strconv.Itoa(int(order.Id)),
		Currency: "USD",
		Amount:   order.TotalPrice(),
	})
	return err
}
