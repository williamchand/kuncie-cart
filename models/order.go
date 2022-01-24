package models

import (
	"time"
)

// Order represent the order model
type OrderDetails struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id" validate:"required"`
	SKU       string    `json:"sku" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Price     float64   `json:"price" validate:"required"`
	Quantity  int64     `json:"quantity" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Order represent the order model
type Order struct {
	ID         int64     `json:"id"`
	TotalPrice float64   `json:"total_price" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// Promotions represent the promotion model
type Promotions struct {
	ID                  int64  `json:"id"`
	ItemsID             int64  `json:"items_id" validate:"required"`
	PromoType           string `json:"promo_type" validate:"required"`
	Promo               string `json:"promo" validate:"required"`
	QuantityRequirement int64  `json:"quantity_requirement" validate:"required"`
}

type Items struct {
	ID                int64     `json:"id"`
	SKU               string    `json:"sku" validate:"required"`
	Name              string    `json:"name" validate:"required"`
	Price             float64   `json:"price" validate:"required"`
	InventoryQuantity int64     `json:"inventory_quantity" validate:"required"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Cart struct {
	ID        int64     `json:"id"`
	ItemsID   int64     `json:"items_id" validate:"required"`
	Quantity  int64     `json:"quantity" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
