package redis

import (
	"context"
	"encoding/json"
	"store/pkg/domain"

	"github.com/go-redis/redis/v9"
)

var _ domain.OrdersRepo = (*OrdersRepo)(nil)

type OrdersRepo struct {
	c *redis.Client
}

func NewOrdersRepo(addr string) *OrdersRepo {
	c := redis.NewClient(&redis.Options{Addr: addr})
	res := &OrdersRepo{c: c}
	return res
}

// Save stores an order in redis using its ID as the key
func (r *OrdersRepo) Save(ctx context.Context, order domain.Order) error {
	serlializedOrder, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = r.c.Set(ctx, order.ID, serlializedOrder, 0).Err()
	return err
}

// Get returns an order from the redis db by its ID
func (r *OrdersRepo) Get(ctx context.Context, id string) (domain.Order, error) {
	serlializedOrder, err := r.c.Get(ctx, id).Result()
	if err != nil {
		return domain.Order{}, err
	}
	var order domain.Order
	err = json.Unmarshal([]byte(serlializedOrder), &order)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}
