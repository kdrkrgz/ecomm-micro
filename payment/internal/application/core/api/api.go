package api

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/domain"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a Application) Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	err := a.db.Save(ctx, &payment)
	if err != nil {
		return domain.Payment{}, nil
	}
	return payment, nil
}
