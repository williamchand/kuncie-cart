package graphql

import "github.com/graphql-go/graphql"

// OrderGraphQL holds order information with graphql object
var OrderGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
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

// OrderEdgeGraphQL holds order edge information with graphql object
var OrderEdgeGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OrderEdge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: OrderGraphQL,
			},
			"cursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// OrderResultGraphQL holds order result information with graphql object
var OrderResultGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OrderResult",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(OrderEdgeGraphQL),
			},
			"pageInfo": &graphql.Field{
				Type: pageInfoGraphQL,
			},
		},
	},
)

var pageInfoGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PageInfo",
		Fields: graphql.Fields{
			"endCursor": &graphql.Field{
				Type: graphql.String,
			},
			"hasNextPage": &graphql.Field{
				Type: graphql.Boolean,
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
			"FetchOrder": &graphql.Field{
				Type:        OrderResultGraphQL,
				Description: "Fetch Order",
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"after": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.orderResolver.FetchOrder,
			},
			"GetOrderByID": &graphql.Field{
				Type:        OrderGraphQL,
				Description: "Get Order By ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.orderResolver.GetOrderByID,
			},
			"GetOrderByTitle": &graphql.Field{
				Type:        OrderGraphQL,
				Description: "Get Order By Title",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.orderResolver.GetOrderByTitle,
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
			"UpdateOrder": &graphql.Field{
				Type:        graphql.String,
				Description: "Update order by certain ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.orderResolver.UpdateOrder,
			},
			"StoreOrder": &graphql.Field{
				Type:        graphql.String,
				Description: "Store a new order",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.orderResolver.StoreOrder,
			},
			"DeleteOrder": &graphql.Field{
				Type:        graphql.String,
				Description: "Delete an order by its ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.orderResolver.DeleteOrder,
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
