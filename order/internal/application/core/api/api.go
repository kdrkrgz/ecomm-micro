package api

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
	"github.com/kdrkrgz/ecomm-micro/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (app *Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := app.db.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}
	payErr := app.payment.Charge(ctx, &order)
	if payErr != nil {
		st, _ := status.FromError(payErr)
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: st.Message(),
		}

		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "payment error")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}
	return order, nil
}
func (app *Application) GetOrder(ctx context.Context, orderId int64) (domain.Order, error) {
	return app.db.Get(ctx, orderId)
}
