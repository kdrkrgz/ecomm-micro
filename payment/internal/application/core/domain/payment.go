package domain

import "time"

type Payment struct {
	Id         int64   `json:"id"`
	CustomerId int64   `json:"customer_id"`
	Status     string  `json:"status"`
	OrderId    int64   `json:"order_id"`
	TotalPrice float32 `json:"total_price"`
	CreatedAt  int64   `json:"created_at"`
}

func NewPayment(customerId int64, orderId int64, totalPrice float32) Payment {
	return Payment{
		Status:     "PENDING",
		CustomerId: customerId,
		OrderId:    orderId,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now().Unix(),
	}
}
