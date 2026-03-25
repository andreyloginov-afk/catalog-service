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
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGuid uuid.UUID `json:"category_guid"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

func (r RequestProductCreate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParametrs
	}

	if r.CategoryGuid.IsNil() {
		return ErrIncorrectParametrs
	}

	if r.Price <= 0 {
		return ErrIncorrectParametrs
	}

	return nil
}

type RequestProductUpdate struct {
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

func (r RequestProductUpdate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParametrs
	}

	if r.CategoryGUID.IsNil() {
		return ErrIncorrectParametrs
	}

	if r.Price <= 0 {
		return ErrIncorrectParametrs
	}

	return nil
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
