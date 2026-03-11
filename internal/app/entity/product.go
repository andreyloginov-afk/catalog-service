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
