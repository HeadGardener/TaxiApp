package http_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	moneyQuery        = "money"
	driverStatusQuery = "status"
)

const (
	minute = time.Minute
)

type IdentityService interface {
	SignUp(ctx context.Context, driver *models.Driver) (string, error)
	SignIn(ctx context.Context, phone, password string) (string, error)
	Check(ctx context.Context, driverID, token string) error
	LogOut(ctx context.Context, driverID string) error
}

type DriverService interface {
	GetProfile(ctx context.Context, driverID string) (*models.Driver, error)
	Update(ctx context.Context, driverID string, driverUpdate *models.Driver) error
	ChangeStatus(ctx context.Context, driverID string, status string) error
	SetInactive(ctx context.Context, driverID string) error
}

type TripService interface {
	ViewAll(ctx context.Context, driverID string) ([]*models.Trip, error)
}

type BalanceService interface {
	Add(ctx context.Context, driverID string, money float64) error
	CashOut(ctx context.Context, driverID string, credentials models.Credentials) error
}

type OrderService interface {
	GetInLine(ctx context.Context, driverID string, taxiType models.TaxiType) error
	ProcessOrder(ctx context.Context, driverID, orderID string, status models.AcceptOrderStatus) error
	Complete(ctx context.Context, driverID, orderID string, status models.CompleteOrderStatus, rating float64) error
	CurrentOrder(driverID string) (models.Order, error)
}

type Handler struct {
	identityService IdentityService
	driverService   DriverService
	tripService     TripService
	balanceService  BalanceService
	orderService    OrderService
}

func NewHandler(identityService IdentityService,
	driverService DriverService,
	tripService TripService,
	balanceService BalanceService,
	orderService OrderService) *Handler {
	return &Handler{
		identityService: identityService,
		driverService:   driverService,
		tripService:     tripService,
		balanceService:  balanceService,
		orderService:    orderService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(minute))

	r.Route("/api", func(r chi.Router) {
		r.Route("/drivers", func(r chi.Router) {
			r.Post("/sign-up", h.signUp)
			r.Post("/sign-in", h.signIn)
			r.Route("/profile", func(r chi.Router) {
				r.Use(h.identifyDriver)
				r.Post("/logout", h.logout)
				r.Get("/", h.profile)
				r.Put("/update", h.updateProfile)
				r.Put("/change-status", h.changeStatus)
				r.Put("/delete", h.setInactive)
			})
		})

		r.Use(h.identifyDriver)
		r.Route("/balance", func(r chi.Router) {
			r.Put("/add", h.addBalance)
			r.Put("/cash-out", h.cashOut)
		})

		r.Route("/trips", func(r chi.Router) {
			r.Get("/", h.viewTrips)
		})

		r.Route("/order", func(r chi.Router) {
			r.Post("/", h.getInLine)
			r.Post("/process", h.processOrder)
			r.Post("/complete", h.completeOrder)
			r.Get("/orders", h.awaitingOrder)
		})
	})

	return r
}
