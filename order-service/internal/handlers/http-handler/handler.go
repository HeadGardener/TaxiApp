package http_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	minute = time.Minute
)

type OrderService interface {
	GetAll(ctx context.Context) ([]models.Order, error)
}

type Handler struct {
	orderService OrderService
}

func NewHandler(orderService OrderService) *Handler {
	return &Handler{orderService: orderService}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(minute))

	r.Route("/api", func(r chi.Router) {
		r.Get("/", h.getAllOrders)
	})

	return r
}

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAll(r.Context())
	if err != nil {
		// return error response
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}
