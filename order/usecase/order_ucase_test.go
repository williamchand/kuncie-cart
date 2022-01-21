package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/williamchand/kuncie-cart/models"
	"github.com/williamchand/kuncie-cart/order/mocks"
	ucase "github.com/williamchand/kuncie-cart/order/usecase"
)

func TestFetch(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := &models.Order{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]*models.Order, 0)
	mockListArtilce = append(mockListArtilce, mockOrder)

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := &models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := models.Order{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := &models.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockOrder, nil).Once()
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockOrder.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected")).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockOrder.ID)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := models.Order{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockOrder := mockOrder
		tempMockOrder.ID = 0
		mockOrderRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(nil, models.ErrNotFound).Once()
		mockOrderRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockOrder)

		assert.NoError(t, err)
		assert.Equal(t, mockOrder.Title, tempMockOrder.Title)
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingOrder := mockOrder
		mockOrderRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(&existingOrder, nil).Once()
		mockAuthor := &models.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		err := u.Store(context.TODO(), &mockOrder)

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := models.Order{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockOrder, nil).Once()

		mockOrderRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockOrder.ID)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("order-is-not-exist", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, nil).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockOrder.ID)

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected Error")).Once()

		u := ucase.NewOrderUsecase(mockOrderRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockOrder.ID)

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := models.Order{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("Update", mock.Anything, &mockOrder).Once().Return(nil)

		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2)

		err := u.Update(context.TODO(), &mockOrder)
		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
