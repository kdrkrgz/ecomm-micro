package ports

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/domain"
)

type APIPort interface {
	Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
