package order

import (
	"context"

	"github.com/williamchandra/kuncie-cart/models"
)

// Usecase represent the order's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.Order, string, error)
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	Update(ctx context.Context, ar *models.Order) error
	GetByTitle(ctx context.Context, title string) (*models.Order, error)
	Store(context.Context, *models.Order) error
	Delete(ctx context.Context, id int64) error
}
