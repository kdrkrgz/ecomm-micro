package db

import (
	"context"
	"fmt"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"time"
)

type Payment struct {
	gorm.Model
	CustomerId int64
	Status     string
	OrderId    int64
	TotalPrice float32
}

type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Payment, error) {
	var paymentEntity Payment
	toContext, _ := context.WithTimeout(ctx, 5*time.Second)
	res := a.db.WithContext(toContext).First(&paymentEntity, id)
	payment := domain.Payment{
		Id:         int64(paymentEntity.ID),
		CustomerId: paymentEntity.CustomerId,
		Status:     paymentEntity.Status,
		OrderId:    paymentEntity.OrderId,
		TotalPrice: paymentEntity.TotalPrice,
		CreatedAt:  paymentEntity.CreatedAt.Unix(),
	}
	return payment, res.Error
}

func (a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	paymentEntity := Payment{
		CustomerId: payment.CustomerId,
		Status:     payment.Status,
		OrderId:    payment.OrderId,
		TotalPrice: payment.TotalPrice,
	}
	toContext, _ := context.WithTimeout(ctx, 5*time.Second)
	res := a.db.WithContext(toContext).Create(&paymentEntity)
	return res.Error
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	dbLogger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond * 100,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceUrl), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("payments"))); err != nil {
		return nil, fmt.Errorf("failed to enable otelgorm: %w", err)
	}
	migrationErr := db.AutoMigrate(&Payment{})
	if migrationErr != nil {
		return nil, fmt.Errorf("migration failed for: %w", migrationErr)
	}
	return &Adapter{db: db}, nil
}
