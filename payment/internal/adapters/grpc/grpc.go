package grpc

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/domain"
	"github.com/kdrkrgz/ecomm-proto/golang/payment"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (a Adapter) Create(ctx context.Context, request *domain.Payment) (*payment.CreatePaymentResponse, error) {
	log.WithContext(ctx).Info("Create payment request")
	newPayment := domain.NewPayment(request.CustomerId, request.OrderId, request.TotalPrice)
	res, err := a.api.Charge(ctx, newPayment)
	if err != nil {
		return nil, err
	}
	transID := strconv.Itoa(int(res.Id))
	return &payment.CreatePaymentResponse{PaymentId: transID}, nil

}
