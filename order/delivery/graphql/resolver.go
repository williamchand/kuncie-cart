package graphql

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/williamchandra/kuncie-cart/order"
	"github.com/williamchandra/kuncie-cart/order/repository"
	"github.com/williamchandra/kuncie-cart/models"
	"time"
)

// OrderEdge holds information of order edge.
type OrderEdge struct {
	Node   models.Order
	Cursor string
}

// OrderResult holds information of order result.
type OrderResult struct {
	Edges    []OrderEdge
	PageInfo PageInfo
}

// PageInfo holds information of page info.
type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type Resolver interface {
	FetchOrder(params graphql.ResolveParams) (interface{}, error)
	GetOrderByID(params graphql.ResolveParams) (interface{}, error)
	GetOrderByTitle(params graphql.ResolveParams) (interface{}, error)

	UpdateOrder(params graphql.ResolveParams) (interface{}, error)
	StoreOrder(params graphql.ResolveParams) (interface{}, error)
	DeleteOrder(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	orderService order.Usecase
}

func (r resolver) FetchOrder(params graphql.ResolveParams) (interface{}, error) {
	ctx := context.Background()
	num := 0
	cursor := ""
	if cursorFromClient, ok := params.Args["after"].(string); ok {
		cursor = cursorFromClient
	}

	if numFromClient, ok := params.Args["first"].(int); ok {
		num = numFromClient
	}

	results, cursorFromService, err := r.orderService.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return nil, err
	}

	edges := make([]OrderEdge, len(results))
	for index, result := range results {
		if result != nil {
			edges[index] = OrderEdge{
				Node: *result,
				Cursor: repository.EncodeCursor(result.CreatedAt),
			}
		}
	}

	isHasNextPage := false
	if len(results) > 0 {
		results, _, err := r.orderService.Fetch(ctx, cursorFromService, int64(1))
		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			isHasNextPage = true
		}
	}

	return OrderResult{
		Edges: edges,
		PageInfo:PageInfo{
			EndCursor: cursorFromService,
			HasNextPage:isHasNextPage,
		},
	}, nil
}

func (r resolver) GetOrderByID(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	result, err := r.orderService.GetByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (r resolver) GetOrderByTitle(params graphql.ResolveParams) (interface{}, error) {
	var (
		title string
		ok bool
	)

	ctx := context.Background()

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	result, err := r.orderService.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func (r resolver) UpdateOrder(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		title, content string
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	if content, ok = params.Args["content"].(string); !ok {
		return nil, fmt.Errorf("content is not string")
	}

	updatedOrder := &models.Order{
		ID: int64(id),
		Title: title,
		Content: content,
		UpdatedAt: time.Now(),
	}

	if err := r.orderService.Update(ctx, updatedOrder); err != nil {
		return nil, err
	}

	return *updatedOrder, nil
}

func (r resolver) StoreOrder(params graphql.ResolveParams) (interface{}, error) {
	var (
		title, content string
		ok bool
	)

	ctx := context.Background()

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	if content, ok = params.Args["content"].(string); !ok {
		return nil, fmt.Errorf("content is not string")
	}

	storedOrder := &models.Order{
		Content: content,
		Title:title,
	}

	if err := r.orderService.Store(ctx, storedOrder); err != nil {
		return nil, err
	}

	return *storedOrder, nil
}

func (r resolver) DeleteOrder(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	if err := r.orderService.Delete(ctx, int64(id)); err != nil {
		return nil, err
	}

	return id, nil
}

func NewResolver(orderService order.Usecase) Resolver {
	return &resolver{
		orderService:orderService,
	}
}
