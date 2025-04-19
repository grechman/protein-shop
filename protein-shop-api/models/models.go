package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	Status    string      `json:"status"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	ID              uuid.UUID `json:"id"`
	OrderID         uuid.UUID `json:"order_id"`
	ProductID       uuid.UUID `json:"product_id"`
	Quantity        int       `json:"quantity"`
	PriceAtPurchase float64   `json:"price_at_purchase"`
	ProductName     string    `json:"name"`
}

type LoyaltyPoint struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Points    int       `json:"points"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}
