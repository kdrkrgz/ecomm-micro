package ports

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
)

type APIPort interface {
	PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrder(ctx context.Context, orderId int64) (domain.Order, error)
}
