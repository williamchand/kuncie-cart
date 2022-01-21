package order

import (
	"context"

	"github.com/williamchandra/kuncie-cart/models"
)

// Repository represent the order's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Order, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	GetByTitle(ctx context.Context, title string) (*models.Order, error)
	Update(ctx context.Context, ar *models.Order) error
	Store(ctx context.Context, a *models.Order) error
	Delete(ctx context.Context, id int64) error
}
