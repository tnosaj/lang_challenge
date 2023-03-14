package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"store/pkg/domain"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/tj/assert"
)

func TestOrdersRepo_Save(t *testing.T) {
	setup()
	defer teardown()
	type fields struct {
		c *redis.Client
	}
	type args struct {
		ctx   context.Context
		order domain.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				c: redisClient,
			},
			args: args{
				ctx: context.TODO(),
				order: domain.Order{
					ID:     "1234567890",
					Status: "Tested",
				},
			},
			wantErr: false,
		},
		{
			name: "nok",
			fields: fields{
				c: redisClient,
			},
			args: args{
				ctx:   context.TODO(),
				order: domain.Order{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrdersRepo{
				c: tt.fields.c,
			}
			if err := r.Save(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("OrdersRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				// we had no error

				var order domain.Order
				result, _ := redisClient.Get(tt.args.ctx, tt.args.order.ID).Result()
				err = json.Unmarshal([]byte(result), &order)

				assert.True(t, reflect.DeepEqual(tt.args.order, order), "expect orders from redis to equal input")
			}
		})
	}
}

func TestOrdersRepo_Get(t *testing.T) {
	setup()
	defer teardown()
	type fields struct {
		c *redis.Client
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Order
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				c: redisClient,
			},
			args: args{
				ctx: context.TODO(),
				id:  "1234567890",
			},
			want: domain.Order{
				ID:     "1234567890",
				Status: "Tested",
			},
			wantErr: false,
		},
		{
			name: "nok",
			fields: fields{
				c: redisClient,
			},
			args: args{
				ctx: context.TODO(),
				id:  "0987654321",
			},
			want:    domain.Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrdersRepo{
				c: tt.fields.c,
			}
			if tt.wantErr == false {
				// put into mock redis
				serlializedOrder, _ := json.Marshal(domain.Order{ID: tt.args.id, Status: "Tested"})
				redisClient.Set(context.TODO(), tt.args.id, serlializedOrder, 0).Err()
			}

			got, err := r.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrdersRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrdersRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdersRepo_Valid(t *testing.T) {
	type fields struct {
		c *redis.Client
	}
	type args struct {
		order domain.Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				c: redisClient,
			},
			args: args{
				order: domain.Order{
					ID:     "myId",
					Status: "myStatus",
				},
			},
			want: true,
		},
		{
			name: "nok: empty",
			fields: fields{
				c: redisClient,
			},
			args: args{
				order: domain.Order{},
			},
			want: false,
		},
		{
			name: "nok: no status",
			fields: fields{
				c: redisClient,
			},
			args: args{
				order: domain.Order{
					ID: "myId",
				},
			},
			want: false,
		},
		{
			name: "nok: no id",
			fields: fields{
				c: redisClient,
			},
			args: args{
				order: domain.Order{
					Status: "myStatus",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrdersRepo{
				c: tt.fields.c,
			}
			if got := r.ValidOrder(tt.args.order); got != tt.want {
				t.Errorf("OrdersRepo.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

var redisServer *miniredis.Miniredis
var redisClient *redis.Client

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func setup() OrdersRepo {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	return OrdersRepo{
		c: redisClient,
	}
}
func teardown() {
	redisServer.Close()
}
