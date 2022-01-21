package usecase

import (
	"context"
	"time"

	"github.com/williamchand/kuncie-cart/order"
	"golang.org/x/sync/errgroup"

	"github.com/williamchandra/kuncie-cart/author"
	"github.com/williamchandra/kuncie-cart/models"
)

type orderUsecase struct {
	orderRepo      order.Repository
	authorRepo     author.Repository
	contextTimeout time.Duration
}

// NewOrderUsecase will create new an orderUsecase object representation of order.Usecase interface
func NewOrderUsecase(a order.Repository, ar author.Repository, timeout time.Duration) order.Usecase {
	return &orderUsecase{
		orderRepo:      a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *orderUsecase) fillAuthorDetails(c context.Context, data []*models.Order) ([]*models.Order, error) {

	g, _ := errgroup.WithContext(c)

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *orderUsecase) Fetch(c context.Context, cursor string, num int64) ([]*models.Order, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listOrder, nextCursor, err := a.orderRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	listOrder, err = a.fillAuthorDetails(ctx, listOrder)
	if err != nil {
		return nil, "", err
	}

	return listOrder, nextCursor, nil
}

func (a *orderUsecase) GetByID(c context.Context, id int64) (*models.Order, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor
	return res, nil
}

func (a *orderUsecase) Update(c context.Context, ar *models.Order) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.orderRepo.Update(ctx, ar)
}

func (a *orderUsecase) GetByTitle(c context.Context, title string) (*models.Order, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.orderRepo.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor

	return res, nil
}

func (a *orderUsecase) Store(c context.Context, m *models.Order) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedOrder, _ := a.GetByTitle(ctx, m.Title)
	if existedOrder != nil {
		return models.ErrConflict
	}

	err := a.orderRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *orderUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedOrder, err := a.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedOrder == nil {
		return models.ErrNotFound
	}
	return a.orderRepo.Delete(ctx, id)
}
