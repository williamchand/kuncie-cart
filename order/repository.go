package order

import (
	"context"

	"github.com/williamchand/kuncie-cart/models"
)

// Repository represent the order's repository contract
type Repository interface {
	GetItems(ctx context.Context, sku []string) (res []*models.Items, err error)
	GetItemsById(ctx context.Context, id []int64) (res []*models.Items, err error)
	GetCart(ctx context.Context) (res []*models.Cart, err error)
	GetPromotions(ctx context.Context, id int64) (*models.Promotions, error)
	CreateCart(ctx context.Context, a *models.Cart) error
	UpdateItems(ctx context.Context, a *models.Items) error
	UpdateCart(ctx context.Context, a *models.Cart) error
	CreateOrder(ctx context.Context, a *models.Order) error
	CreateOrderDetails(ctx context.Context, a *models.OrderDetails) error
	DeleteCart(ctx context.Context) error
}
