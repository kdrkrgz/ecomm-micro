package ports

import "github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/domain"

import "context"

type DBPort interface {
	Get(ctx context.Context, id string) (domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) error
}
