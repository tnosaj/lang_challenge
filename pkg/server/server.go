package server

import (
	"context"
	"encoding/json"
	"net/http"
	"store/pkg/domain"
	"store/pkg/utils"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	log.Tracef("Start Create %v", r)
	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), time.Microsecond*200)
	defer cancel()

	order := domain.Order{
		ID:     utils.NewID(),
		Status: "CREATED",
	}
	err := s.repo.Save(ctx, order)
	if err != nil {
		s.Metrics.HttpErrors.WithLabelValues("Create").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Create Failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(order.ID))
	log.Tracef("End Create %v", r)
}

// Get will try to fetch the order by the ID provided in the url
func (s Server) Get(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Start Get %v", r)
	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), time.Microsecond*200)
	defer cancel()

	var id string

	// potentially use a different router to avoid this
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) == 3 {
		id = parts[2]
	} else {
		s.Metrics.HttpErrors.WithLabelValues("BadRequest").Inc()
		w.WriteHeader(http.StatusBadRequest)
		log.Infof("Get Failed: Bad request - %s", r.URL.String())
		return
	}
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			log.Infof("Get Failed: id %s not found", id)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("Get Failed: %s", err)
		}
		s.Metrics.HttpErrors.WithLabelValues("Get").Inc()
		return
	}
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(order)
	if err != nil {
		s.Metrics.HttpErrors.WithLabelValues("Get").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Get Failed: Error marshaling json - %s", err)
		return
	}
	w.Write(body)
	log.Tracef("End Get %v", r)
	return
}
