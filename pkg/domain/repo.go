package domain

import "context"

type OrdersRepo interface {
	Save(ctx context.Context, order Order) error
	Get(ctx context.Context, id string) (Order, error)
	ValidOrder(order Order) bool
	Shutdown(ctx context.Context)
}
