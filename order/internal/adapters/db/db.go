package db

import (
	"context"
	"fmt"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"time"
)

type Order struct {
	gorm.Model
	CustomerId int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderId     uint
}

type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(ctx context.Context, orderId int64) (domain.Order, error) {
	var order Order
	result := a.db.WithContext(ctx).Preload("OrderItems").First(&order, orderId)
	if result.Error != nil {
		return domain.Order{}, result.Error
	}

	return domain.Order{
		Id:         int64(order.ID),
		CustomerId: order.CustomerId,
		Status:     order.Status,
		OrderItems: a.mapOrderItems(order.OrderItems),
		CreatedAt:  order.CreatedAt.Unix(),
	}, nil
}

func (a Adapter) mapOrderItems(orderItems []OrderItem) []domain.OrderItem {
	var items []domain.OrderItem
	for _, item := range orderItems {
		items = append(items, domain.OrderItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}
	return items
}

func (a Adapter) Save(ctx context.Context, order *domain.Order) error {
	var items []OrderItem
	for _, item := range order.OrderItems {
		items = append(items, OrderItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}
	dbOrder := Order{
		CustomerId: order.CustomerId,
		Status:     order.Status,
		OrderItems: items,
	}
	result := a.db.WithContext(ctx).Create(&dbOrder)
	if result.Error != nil {
		return result.Error
	}
	order.Id = int64(dbOrder.ID)
	return nil
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	dbLogger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond * 100,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceUrl), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, fmt.Errorf("db conenction failed: %w", err)
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("order"))); err != nil {
		return nil, fmt.Errorf("failed to enable otelgorm: %w", err)
	}
	migrationErr := db.AutoMigrate(&Order{}, &OrderItem{})
	if migrationErr != nil {
		return nil, fmt.Errorf("migration failed: %w", migrationErr)
	}
	return &Adapter{db: db}, nil
}
