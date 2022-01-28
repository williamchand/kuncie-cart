package graphql

import "github.com/graphql-go/graphql"

// OrderGraphQL holds order information with graphql object
var OrderGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"total_price": &graphql.Field{
				Type: graphql.Float,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// CartGraphQL holds order information with graphql object
var CartGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Cart",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"items_id": &graphql.Field{
				Type: graphql.Int,
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// Schema is struct which has method for Query and Mutation. Please init this struct using constructor function.
type Schema struct {
	orderResolver Resolver
}

// Query initializes config schema query for graphql server.
func (s Schema) Query() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Placeholder": &graphql.Field{
				Type:        graphql.String,
				Description: "Confirm all order at the cart",
				Args:        graphql.FieldConfigArgument{},
				Resolve:     s.orderResolver.Placeholder,
			},
		},
	}

	return graphql.NewObject(objectConfig)
}

// Mutation initializes config schema mutation for graphql server.
func (s Schema) Mutation() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"ConfirmOrder": &graphql.Field{
				Type:        OrderGraphQL,
				Description: "Confirm all order at the cart",
				Args: graphql.FieldConfigArgument{
					"placeholder": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.orderResolver.ConfirmOrder,
			},
			"AddCart": &graphql.Field{
				Type:        CartGraphQL,
				Description: "Store a new order",
				Args: graphql.FieldConfigArgument{
					"sku": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"quantity": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.orderResolver.AddCart,
			},
		},
	}

	return graphql.NewObject(objectConfig)
}

// NewSchema initializes Schema struct which takes resolver as the argument.
func NewSchema(orderResolver Resolver) Schema {
	return Schema{
		orderResolver: orderResolver,
	}
}
