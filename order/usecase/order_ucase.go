package usecase

import (
	"context"
	"time"

	"github.com/williamchand/kuncie-cart/order"

	"github.com/williamchand/kuncie-cart/models"
)

type orderUsecase struct {
	orderRepo      order.Repository
	contextTimeout time.Duration
}

// NewOrderUsecase will create new an orderUsecase object representation of order.Usecase interface
func NewOrderUsecase(a order.Repository, timeout time.Duration) order.Usecase {
	return &orderUsecase{
		orderRepo:      a,
		contextTimeout: timeout,
	}
}

func (a *orderUsecase) GetItems(c context.Context, sku []string) (result []*models.Items, err error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.orderRepo.GetItems(ctx, sku)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *orderUsecase) GetItemsById(c context.Context, id []int64) (result []*models.Items, err error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.orderRepo.GetItemsById(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a *orderUsecase) GetCart(c context.Context) (result []*models.Cart, err error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.orderRepo.GetCart(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *orderUsecase) GetPromotions(c context.Context, id int64) (*models.Promotions, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.orderRepo.GetPromotions(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *orderUsecase) CreateCart(c context.Context, m *models.Cart) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.orderRepo.CreateCart(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *orderUsecase) CreateOrder(c context.Context, m *models.Order) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.orderRepo.CreateOrder(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *orderUsecase) CreateOrderDetails(c context.Context, m *models.OrderDetails) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.orderRepo.CreateOrderDetails(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *orderUsecase) UpdateCart(c context.Context, ar *models.Cart) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.orderRepo.UpdateCart(ctx, ar)
}

func (a *orderUsecase) DeleteCart(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.orderRepo.DeleteCart(ctx)
}
