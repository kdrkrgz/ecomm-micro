package db

import "gorm.io/gorm"

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
