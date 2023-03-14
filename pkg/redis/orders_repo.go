package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"store/pkg/domain"

	"github.com/go-playground/validator"
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
	if !r.ValidOrder(order) {
		return fmt.Errorf("Not a valid order")
	}
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

func (r *OrdersRepo) ValidOrder(order domain.Order) bool {
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		return false
	}
	return true
}
