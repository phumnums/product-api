package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description *string
	SalePrice   *float64
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
