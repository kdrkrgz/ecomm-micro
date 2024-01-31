package ports

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
)

type PaymentPort interface {
	Charge(context.Context, *domain.Order) error
}
