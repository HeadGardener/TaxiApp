package http_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	walletIDParam          = "wallet_id"
	transactionIDParam     = "transaction_id"
	phoneQuery             = "phone"
	walletTypeQuery        = "wallet_type"
	moneyQuery             = "money"
	transactionStatusQuery = "status"
	orderIDParam           = "order_id"
)

const (
	minute = time.Minute
)

type IdentityService interface {
	SignUp(ctx context.Context, user *models.User) (string, error)
	SignIn(ctx context.Context, phone, password string) (string, error)
	Check(ctx context.Context, userID, token string) error
	LogOut(ctx context.Context, userID string) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID string) (models.UserProfile, error)
	Update(ctx context.Context, userID string, userUpdate *models.User) error
	SetInactive(ctx context.Context, userID string) error
}

type WalletService interface {
	Create(ctx context.Context, userID, card string) (string, error)
	TopUp(ctx context.Context, userID, walletID string, money float64) (int, error)
	GetByID(ctx context.Context, userID, walletID string) (models.Wallet, error)
	ViewAll(ctx context.Context, userID string) ([]models.Wallet, error)

	CreateFamilyWallet(ctx context.Context, userID, walletID string) (string, error)
	GetFamilyWalletByID(ctx context.Context, userID, walletID string) (models.FamilyWallet, error)
	DeleteFamilyWallet(ctx context.Context, userID, walletID string) error
	SetFixedBalance(ctx context.Context, userID, walletID string, fixedBalance float64) error
	AddFamilyBalance(ctx context.Context, userID, walletID, famWalletID string, amount float64) error
	AddMember(ctx context.Context, userID, walletID, phone string) error
	ViewMembers(ctx context.Context, userID, walletID string) ([]string, error)
	DeleteMember(ctx context.Context, userID, walletID, phone string) error
	ViewAllFamily(ctx context.Context, userID string) ([]models.FamilyWallet, error)
	ViewMemberships(ctx context.Context, userID string) ([]models.FamilyWallet, error)

	PickWalletAndPay(ctx context.Context, userID string, walletType models.WalletType,
		transactionID int, credentials models.Credentials) error
}

type TransactionService interface {
	Create(ctx context.Context, money float64) (int, error)
	ViewAll(ctx context.Context, userID, walletID string, walletType models.WalletType) ([]models.Transaction, error)
	Confirm(ctx context.Context, walletType models.WalletType, transactionID int) error
	Cancel(ctx context.Context, walletType models.WalletType, transactionID int, credentials models.Credentials) error
}

type TripService interface {
	ViewAll(ctx context.Context, userID string) ([]models.Trip, error)
}

type OrderService interface {
	SendOrder(ctx context.Context, userID string, order *models.Order) (string, error)
	SendComment(ctx context.Context, orderID, comment string) error
	Get(orderID, userID string) (*models.UserOrder, error)
	Delete(orderID, userID string) error
}

type Handler struct {
	identityService    IdentityService
	userService        UserService
	walletService      WalletService
	transactionService TransactionService
	tripService        TripService
	orderService       OrderService
}

func NewHandler(identityService IdentityService,
	userService UserService,
	walletService WalletService,
	transactionService TransactionService,
	tripService TripService,
	orderService OrderService) *Handler {
	return &Handler{
		identityService:    identityService,
		userService:        userService,
		walletService:      walletService,
		transactionService: transactionService,
		tripService:        tripService,
		orderService:       orderService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:9090"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/sign-up", h.signUp)
			r.Post("/sign-in", h.signIn)
			r.Route("/profile", func(r chi.Router) {
				r.Use(h.identifyUser)
				r.Post("/logout", h.logout)
				r.Get("/", h.profile)
				r.Put("/update", h.updateProfile)
				r.Put("/delete", h.setInactive)
			})
		})

		r.Route("/wallets", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Route("/personal", func(r chi.Router) {
				r.Post("/", h.createWallet)
				r.Get("/{wallet_id}", h.viewWallet)
				r.Get("/", h.viewWallets)

				r.Route("/{wallet_id}", func(r chi.Router) {
					r.Put("/balance", h.topUp)
				})
			})

			r.Route("/family", func(r chi.Router) {
				r.Post("/{wallet_id}", h.createFamilyWallet)
				r.Get("/{wallet_id}", h.viewFamilyWallet)
				r.Get("/", h.viewFamilyWallets)
				r.Get("/memberships", h.viewMemberships)
				r.Delete("/{wallet_id}", h.deleteFamilyWallet)

				r.Route("/{wallet_id}", func(r chi.Router) {
					r.Put("/balance", h.topUpFamily)
					r.Put("/set-balance", h.setFixedBalance)
					r.Route("/members", func(r chi.Router) {
						r.Post("/", h.addMember)
						r.Get("/", h.viewMembers)
						r.Delete("/", h.deleteMember)
					})
				})
			})

			r.Post("/pay/{transaction_id}", h.pay)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Get("/{wallet_id}", h.viewTransactions)
			r.Put("/{transaction_id}", h.processTransaction)
		})

		r.Route("/order", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Get("/", h.viewTrips)
			r.Post("/", h.sendOrder)
			r.Post("/comment", h.sendComment)
			r.Get("/{order_id}", h.awaitingOrders)
		})
	})

	return r
}
