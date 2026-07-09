package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID           int64     `bun:"id,autoincrement"`
	GUID         uuid.UUID `bun:"guid,pk"`
	Name         string    `bun:"name"`
	Description  *string   `bun:"description"`
	Price        float64   `bun:"price"`
	CategoryGuid uuid.UUID `bun:"category_guid"`
	CreatedAt    time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,default:current_timestamp"`
}

type RequestProductCreate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Description  *string   `json:"description" binding:"omitempty,max=500"`
	Price        float64   `json:"price" binding:"required,gt=0"`
	CategoryGuid uuid.UUID `json:"category_guid" binding:"required,uuid"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type RequestProductUpdate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Description  *string   `json:"description" binding:"omitempty,max=500"`
	Price        float64   `json:"price" binding:"required,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"required,gt=0"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type RequestProductList struct {
	Name         *string    `json:"name"`
	CategoryGUID *uuid.UUID `json:"category_guid"`
	MinPrice     *int64     `json:"min_price"`
	MaxPrice     *int64     `json:"max_price"`
}

type ResponseProduct struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
