package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	orderHttp "github.com/williamchandra/kuncie-cart/order/delivery/http"
	"github.com/williamchandra/kuncie-cart/order/mocks"
	"github.com/williamchandra/kuncie-cart/models"
)

func TestFetch(t *testing.T) {
	var mockOrder models.Order
	err := faker.FakeData(&mockOrder)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	mockListOrder := make([]*models.Order, 0)
	mockListOrder = append(mockListOrder, &mockOrder)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListOrder, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/order?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := orderHttp.OrderHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchOrder(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/order?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := orderHttp.OrderHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchOrder(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockOrder models.Order
	err := faker.FakeData(&mockOrder)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockOrder.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(&mockOrder, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/order/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("order/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := orderHttp.OrderHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockOrder := models.Order{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockOrder := mockOrder
	tempMockOrder.ID = 0
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockOrder)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/order", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/order")

	handler := orderHttp.OrderHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockOrder models.Order
	err := faker.FakeData(&mockOrder)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockOrder.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/order/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("order/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := orderHttp.OrderHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
