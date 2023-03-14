package server

import (
	"encoding/json"
	"net/http"
	"store/pkg/domain"
	"store/pkg/utils"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	repo    domain.OrdersRepo
	Metrics ServerMetrics
}

type ServerMetrics struct {
	HttpErrors *prometheus.CounterVec
}

// NewServer returns a Server with a repo for the OrdersRepo domain
func NewServer(repo domain.OrdersRepo) *Server {

	httpErrorReuests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_requests",
			Help: "The total number of failed http requests",
		},
		[]string{"function"},
	)

	prometheus.MustRegister(httpErrorReuests)

	res := &Server{
		repo: repo,
		Metrics: ServerMetrics{
			HttpErrors: httpErrorReuests,
		},
	}
	return res
}

// Create will create a new order and save it to redis, returning the order.ID
func (s Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	order := domain.Order{
		ID:     utils.NewID(),
		Status: "CREATED",
	}
	err := s.repo.Save(ctx, order)
	if err != nil {
		s.Metrics.HttpErrors.WithLabelValues("Create").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(order.ID))
}

// Get will try to fetch the order by the ID provided in the url
func (s Server) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var id string
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) > 2 {
		id = parts[2]
	}
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		s.Metrics.HttpErrors.WithLabelValues("Get").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(order)
	if err != nil {
		s.Metrics.HttpErrors.WithLabelValues("Get").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
