package graphql

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/williamchand/kuncie-cart/models"
	"github.com/williamchand/kuncie-cart/order"
)

type Resolver interface {
	AddCart(params graphql.ResolveParams) (interface{}, error)
	ConfirmOrder(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	orderService order.Usecase
}

func (r resolver) ConfirmOrder(params graphql.ResolveParams) (interface{}, error) {
	var ()

	ctx := context.Background()
	carts, err := r.orderService.GetCart(ctx)
	promotion_carts := make([]*models.Cart, 0)
	order := make([]*models.OrderDetails, 0)
	for i := range carts {
		promotion, err := r.orderService.GetPromotions(ctx, carts[i].ItemsID)
		if err != nil {
			return nil, err
		}
		item_list := []int64{carts[i].ItemsID}
		item_detail, err := r.orderService.GetItemsById(ctx, item_list)
		if err != nil {
			return nil, err
		}
		if promotion.PromoType == "free_items" {
			promotion_quantity := int64(math.Floor(float64(carts[i].Quantity) / float64(promotion.QuantityRequirement)))
			if promotion_quantity > 0 {
				found_promotion := false
				for j := range promotion_carts {
					if promotion_carts[j].ItemsID == carts[i].ItemsID {
						found_promotion = true
						promotion_carts[i].Quantity += promotion_quantity
						break
					}
				}
				if !found_promotion {
					promotion_carts = append(promotion_carts, &models.Cart{
						ItemsID:  carts[i].ItemsID,
						Quantity: promotion_quantity,
					})
				}
			}
			price := item_detail[0].Price * float64(carts[i].Quantity)
			order = append(order, &models.OrderDetails{
				SKU:       item_detail[0].SKU,
				Name:      item_detail[0].Name,
				Price:     price,
				Quantity:  carts[i].Quantity,
				PromoType: "",
			})
		} else if promotion.PromoType == "bonus_price" {
			promoValue, err := strconv.ParseFloat(promotion.Promo, 64)
			if err != nil {
				return nil, err
			}
			price := 0.0
			if promotion.PromoType == "bonus_price" {
				price = item_detail[0].Price*float64(carts[i].Quantity%promotion.QuantityRequirement) + promoValue*math.Floor(float64(carts[i].Quantity)/float64(promotion.QuantityRequirement))
			} else if promotion.PromoType == "discount_items" {
				price = item_detail[0].Price * float64(carts[i].Quantity)
				if promotion.QuantityRequirement >= carts[i].Quantity {
					price = price * promoValue
				}
			}
			promo_type := ""
			if promotion.QuantityRequirement >= carts[i].Quantity {
				promo_type = promotion.PromoType
			}
			order = append(order, &models.OrderDetails{
				SKU:       item_detail[0].SKU,
				Name:      item_detail[0].Name,
				Price:     price,
				Quantity:  carts[i].Quantity,
				PromoType: promo_type,
			})
		}
	}
	for i := range promotion_carts {
		item_list := []int64{promotion_carts[i].ItemsID}
		item_detail, err := r.orderService.GetItemsById(ctx, item_list)
		if err != nil {
			return nil, err
		}
		order = append(order, &models.OrderDetails{
			SKU:       item_detail[0].SKU,
			Name:      item_detail[0].Name,
			Price:     0.0,
			Quantity:  promotion_carts[i].Quantity,
			PromoType: "free_items",
		})
	}
	if err != nil {
		return nil, err
	}
	createOrder := &models.Order{
		TotalPrice: 0.0,
	}
	for i := range order {
		createOrder.TotalPrice += order[i].Price
	}

	if err := r.orderService.CreateOrder(ctx, createOrder); err != nil {
		return nil, err
	}

	for i := range order {
		order[i].OrderID = createOrder.ID
		if err := r.orderService.CreateOrderDetails(ctx, order[i]); err != nil {
			return nil, err
		}
	}

	return *createOrder, nil
}

func (r resolver) AddCart(params graphql.ResolveParams) (interface{}, error) {
	var (
		sku      string
		quantity int64
		ok       bool
	)

	ctx := context.Background()

	if sku, ok = params.Args["sku"].(string); !ok || sku == "" {
		return nil, fmt.Errorf("sku is empty or not string")
	}
	skuList := []string{sku}
	items, err := r.orderService.GetItems(ctx, skuList)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("sku is not valid value")
	}
	if quantity, ok = params.Args["quantity"].(int64); !ok || quantity == 0 {
		return nil, fmt.Errorf("quantity is not integer or zero")
	}

	carts, err := r.orderService.GetCart(ctx)
	if err != nil {
		return nil, err
	}
	found := int64(0)
	for i := range carts {
		if carts[i].ItemsID == items[0].ID {
			carts[i].Quantity += quantity
			found = carts[i].Quantity
			break
		}
	}
	if found == 0 {
		carts = append(carts, &models.Cart{
			ItemsID:  items[0].ID,
			Quantity: quantity,
		})
	}
	promotion_carts := make([]*models.Cart, 0)
	for i := range carts {
		promotion, err := r.orderService.GetPromotions(ctx, carts[i].ItemsID)
		if err != nil {
			return nil, err
		}
		if promotion.PromoType == "free_items" {
			promotion_quantity := int64(math.Floor(float64(carts[i].Quantity) / float64(promotion.QuantityRequirement)))
			if promotion_quantity > 0 {
				promotion_carts = append(promotion_carts, &models.Cart{
					ItemsID:  carts[i].ItemsID,
					Quantity: promotion_quantity,
				})
			}
		}
	}
	for i := range promotion_carts {
		found_promotion := false
		for j := range carts {
			if carts[j].ItemsID == items[0].ID {
				found_promotion = true
				carts[i].Quantity += promotion_carts[i].Quantity
				break
			}
		}
		if !found_promotion {
			carts = append(carts, &models.Cart{
				ItemsID:  promotion_carts[i].ItemsID,
				Quantity: quantity,
			})
		}
	}
	item_list := make([]int64, len(carts))
	for i := range carts {
		item_list[i] = carts[i].ItemsID
	}
	items_availability, err := r.orderService.GetItemsById(ctx, item_list)
	if err != nil {
		return nil, err
	}
	for i := range items_availability {
		for j := range carts {
			if items_availability[i].ID == carts[j].ItemsID {
				if items_availability[i].InventoryQuantity < carts[j].Quantity {
					return nil, fmt.Errorf("cannot add the items")
				}
				break
			}
		}
	}
	cartsAns := &models.Cart{
		ItemsID:  items[0].ID,
		Quantity: quantity,
	}
	if found == 0 {
		err = r.orderService.CreateCart(ctx, cartsAns)
	} else {
		cartsAns.Quantity = found
		err = r.orderService.UpdateCart(ctx, cartsAns)
	}
	if err != nil {
		return nil, err
	}

	return *cartsAns, nil
}

func NewResolver(orderService order.Usecase) Resolver {
	return &resolver{
		orderService: orderService,
	}
}
