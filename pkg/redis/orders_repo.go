package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"store/pkg/domain"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var _ domain.OrdersRepo = (*OrdersRepo)(nil)

type OrdersRepo struct {
	c       *redis.Client
	Metrics domain.OrderMetrics
}

// NewOrdersRepo creates a new redis connection and metrics for orderrepo
func NewOrdersRepo(addr, password string, poolSize, timeout int) *OrdersRepo {

	redisRequestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "redis_request_duration_seconds",
		Help:    "Histogram for the runtime of a simple method function.",
		Buckets: prometheus.LinearBuckets(0.00, 0.002, 75),
	}, []string{"function"})

	redisErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_error_requests",
			Help: "The total number of failed requests",
		},
		[]string{"function"},
	)

	prometheus.MustRegister(redisRequestDuration)
	prometheus.MustRegister(redisErrorReuests)

	redisSettings := redis.Options{
		Addr:         addr,
		PoolSize:     poolSize,
		DialTimeout:  time.Second * time.Duration(timeout),
		ReadTimeout:  time.Second * time.Duration(timeout),
		WriteTimeout: time.Second * time.Duration(timeout),
	}
	if len(password) > 0 {
		redisSettings.Password = password
	}

	c := redis.NewClient(&redisSettings)

	ping := c.Ping(context.TODO()).Err()
	if ping != nil {
		log.Fatalf("Error connecting to redis: %s", ping)
	} else {
		log.Debug("Successfully connected to redis")
	}

	res := &OrdersRepo{
		c: c,
		Metrics: domain.OrderMetrics{
			RedisLatency: redisRequestDuration,
			RedisErrors:  redisErrorReuests,
		},
	}
	return res
}

// Save stores an order in redis using its ID as the key
func (r *OrdersRepo) Save(ctx context.Context, order domain.Order) error {
	log.Trace("Start Save for %v", order)
	timer := prometheus.NewTimer(r.Metrics.RedisLatency.WithLabelValues("Save"))
	if !r.ValidOrder(order) {
		r.Metrics.RedisErrors.WithLabelValues("SaveValidate").Inc()
		timer.ObserveDuration()
		return fmt.Errorf("Not a valid order")
	}
	serlializedOrder, err := json.Marshal(order)
	if err != nil {
		r.Metrics.RedisErrors.WithLabelValues("Save").Inc()
		timer.ObserveDuration()
		return err
	}

	err = r.c.Set(ctx, order.ID, serlializedOrder, 0).Err()
	timer.ObserveDuration()
	log.Trace("End Save for %v", order)
	return err
}

// Get returns an order from the redis db by its ID
func (r *OrdersRepo) Get(ctx context.Context, id string) (domain.Order, error) {
	log.Trace("Start Get for %id", id)

	timer := prometheus.NewTimer(r.Metrics.RedisLatency.WithLabelValues("Get"))
	serlializedOrder, err := r.c.Get(ctx, id).Result()
	if err == redis.Nil {
		r.Metrics.RedisErrors.WithLabelValues("EmptyGet").Inc()
		timer.ObserveDuration()
		return domain.Order{}, fmt.Errorf("not found")
	}
	if err != nil {
		r.Metrics.RedisErrors.WithLabelValues("Get").Inc()
		timer.ObserveDuration()
		return domain.Order{}, err
	}

	var order domain.Order
	err = json.Unmarshal([]byte(serlializedOrder), &order)
	if err != nil {
		r.Metrics.RedisErrors.WithLabelValues("GetUnmarshal").Inc()
		timer.ObserveDuration()
		return domain.Order{}, err
	}
	timer.ObserveDuration()
	log.Trace("End Get for %s", id)
	return order, nil
}

// ValidOrder checks if all required parts of an order exist
func (r *OrdersRepo) ValidOrder(order domain.Order) bool {
	log.Trace("Start ValidOrder for %v", order)
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		r.Metrics.RedisErrors.WithLabelValues("Validate").Inc()
		log.Errorf("Could not validate order: %v", order)
		return false
	}
	log.Trace("End ValidOrder for %v", order)
	return true
}

// Shutdown makes sure redis connections are closed correctly
func (r *OrdersRepo) Shutdown(ctx context.Context) {
	// teh lolz, thanks examples...
	// if err := c.RedisNative.FlushAll(ctx).Err(); err != nil {
	//   logrus.Fatalf("goredis - failed to flush: %v", err)
	// }
	if err := r.c.Close(); err != nil {
		log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
	}
}
