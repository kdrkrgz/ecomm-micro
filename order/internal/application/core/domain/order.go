package domain

import "time"

type Order struct {
	Id         int64       `json:"id"`
	CustomerId int64       `json:"customer_id"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt  int64       `json:"created_at"`
}

type OrderItem struct {
	Id          int64   `json:"id"`
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

func NewOrder(customerId int64, orderItems []OrderItem) Order {
	return Order{
		CustomerId: customerId,
		OrderItems: orderItems,
		Status:     "Pending",
		CreatedAt:  time.Now().Unix(),
	}
}

func (order *Order) TotalPrice(orderItems []OrderItem) float32 {
	var totalPrice float32
	for _, orderItem := range orderItems {
		totalPrice += orderItem.UnitPrice * float32(orderItem.Quantity)
	}
	return totalPrice
}
