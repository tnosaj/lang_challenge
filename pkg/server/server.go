package server

import (
	"encoding/json"
	"net/http"
	"store/pkg/domain"
	"store/pkg/utils"
	"strings"
)

type Server struct {
	repo domain.OrdersRepo
}

func NewServer(repo domain.OrdersRepo) *Server {
	res := &Server{
		repo: repo,
	}
	return res
}

func (s Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	order := domain.Order{
		ID:     utils.NewID(),
		Status: "CREATED",
	}
	err := s.repo.Save(ctx, order)
	if err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
