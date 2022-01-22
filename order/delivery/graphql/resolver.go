package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/williamchand/kuncie-cart/models"
	"github.com/williamchand/kuncie-cart/order"
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
	AddCart(params graphql.ResolveParams) (interface{}, error)
	ConfirmOrder(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	orderService order.Usecase
}

func (r resolver) ConfirmOrder(params graphql.ResolveParams) (interface{}, error) {
	var (
		id             int
		title, content string
		ok             bool
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
		ID:        int64(id),
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now(),
	}

	if err := r.orderService.Update(ctx, updatedOrder); err != nil {
		return nil, err
	}

	return *updatedOrder, nil
}

func (r resolver) AddCart(params graphql.ResolveParams) (interface{}, error) {
	var (
		title, content string
		ok             bool
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
		Title:   title,
	}

	if err := r.orderService.Store(ctx, storedOrder); err != nil {
		return nil, err
	}

	return *storedOrder, nil
}

// func (r resolver) DeleteOrder(params graphql.ResolveParams) (interface{}, error) {
// 	var (
// 		id int
// 		ok bool
// 	)

// 	ctx := context.Background()
// 	if id, ok = params.Args["id"].(int); !ok || id == 0 {
// 		return nil, fmt.Errorf("id is not integer or zero")
// 	}

// 	if err := r.orderService.Delete(ctx, int64(id)); err != nil {
// 		return nil, err
// 	}

// 	return id, nil
// }

func NewResolver(orderService order.Usecase) Resolver {
	return &resolver{
		orderService: orderService,
	}
}
