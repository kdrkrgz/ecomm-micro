package grpc

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/domain"
	"github.com/kdrkrgz/ecomm-proto/golang/order"
	log "github.com/sirupsen/logrus"
)

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	log.WithContext(ctx).Info("Order Creating")
	var items []domain.OrderItem
	for _, item := range request.OrderItems {
		items = append(items, domain.OrderItem{
			ProductCode: item.Sku,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}
	newOrder := domain.NewOrder(request.UserId, items)
	res, err := a.api.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: res.Id}, nil
}

func (a Adapter) Get(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	result, err := a.api.GetOrder(ctx, request.OrderId)
	var orderItems []*order.OrderItem
	for _, orderItem := range result.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			Sku:       orderItem.ProductCode,
			UnitPrice: orderItem.UnitPrice,
			Quantity:  orderItem.Quantity,
		})
	}
	if err != nil {
		return nil, err
	}
	return &order.GetOrderResponse{UserId: result.CustomerId, OrderItems: orderItems}, nil
}
