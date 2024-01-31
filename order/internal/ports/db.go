package ports

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, orderId int64) (domain.Order, error)
	Save(context.Context, *domain.Order) error
}
