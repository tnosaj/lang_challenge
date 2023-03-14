package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"store/pkg/domain"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tj/assert"
)

func TestOrdersRepo_Save(t *testing.T) {
	repo := setup()
	defer teardown()
	type fields struct {
		c       *redis.Client
		Metrics domain.OrderMetrics
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       tt.fields.c,
				Metrics: tt.fields.Metrics,
			}
			if err := r.Save(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("OrdersRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				// we had no error

				var order domain.Order
				result, _ := r.c.Get(tt.args.ctx, tt.args.order.ID).Result()
				err = json.Unmarshal([]byte(result), &order)

				assert.True(t, reflect.DeepEqual(tt.args.order, order), "expect orders from redis to equal input")
			}
		})
	}
}

func TestOrdersRepo_Get(t *testing.T) {
	repo := setup()
	defer teardown()
	type fields struct {
		c       *redis.Client
		Metrics domain.OrderMetrics
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       tt.fields.c,
				Metrics: tt.fields.Metrics,
			}
			if tt.wantErr == false {
				// put into mock redis
				serlializedOrder, _ := json.Marshal(domain.Order{ID: tt.args.id, Status: "Tested"})
				r.c.Set(context.TODO(), tt.args.id, serlializedOrder, 0).Err()
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
	repo := setup()
	defer teardown()
	type fields struct {
		c       *redis.Client
		Metrics domain.OrderMetrics
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       repo.c,
				Metrics: repo.Metrics,
			},
			args: args{
				order: domain.Order{},
			},
			want: false,
		},
		{
			name: "nok: no status",
			fields: fields{
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       repo.c,
				Metrics: repo.Metrics,
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
				c:       tt.fields.c,
				Metrics: tt.fields.Metrics,
			}
			if got := r.ValidOrder(tt.args.order); got != tt.want {
				t.Errorf("OrdersRepo.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

var redisServer *miniredis.Miniredis

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func setup() OrdersRepo {
	redisServer = mockRedis()
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	redisRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "redis_request_duration_seconds",
		Help:    "Histogram for the runtime of a simple method function.",
		Buckets: prometheus.LinearBuckets(0.00, 0.002, 75),
	}, []string{"function"})

	redisErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_error_requests",
			Help: "The total number of failed requests",
		},
		[]string{"function"},
	)
	prometheus.MustRegister(redisRequestDuration)
	prometheus.MustRegister(redisErrorReuests)

	return OrdersRepo{
		c: redisClient,
		Metrics: domain.OrderMetrics{
			RedisLatency: redisRequestDuration,
			RedisErrors:  redisErrorReuests,
		},
	}
}
func teardown() {
	redisServer.Close()
}
