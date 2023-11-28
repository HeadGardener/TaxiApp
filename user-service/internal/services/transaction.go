package services

import (
	"context"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
)

type TransactionStorage interface {
	Create(ctx context.Context, transactionInput *models.Transaction) (int, error)
	GetAll(ctx context.Context, walletID string) ([]models.Transaction, error)
	Confirm(ctx context.Context, transactionID int, status string) error
	Cancel(ctx context.Context, transactionID int, status string, walletType models.WalletType,
		credentials models.Credentials) error
}

type WalletProvider interface {
	IsOwner(ctx context.Context, walletID, userID string) error
	IsFamilyOwner(ctx context.Context, walletID, userID string) error
	IsMember(ctx context.Context, walletID, userID string) error
}

type TransactionService struct {
	transactionStorage TransactionStorage
	walletProvider     WalletProvider
}

func NewTransactionService(transactionStorage TransactionStorage, walletProvider WalletProvider) *TransactionService {
	return &TransactionService{
		transactionStorage: transactionStorage,
		walletProvider:     walletProvider,
	}
}

func (s *TransactionService) Create(ctx context.Context, money float64) (int, error) {
	transaction, err := models.BuildTransaction(models.Spent, models.Create, money)
	if err != nil {
		return 0, err
	}

	return s.transactionStorage.Create(ctx, &transaction)
}

func (s *TransactionService) ViewAll(ctx context.Context, userID, walletID string,
	walletType models.WalletType) ([]models.Transaction, error) {
	if walletType == models.Personal {
		if err := s.walletProvider.IsOwner(ctx, walletID, userID); err != nil {
			return nil, err
		}

		return s.transactionStorage.GetAll(ctx, walletID)
	}

	if walletType == models.Family {
		if err := s.walletProvider.IsFamilyOwner(ctx, walletID, userID); err != nil {
			return nil, err
		}

		return s.transactionStorage.GetAll(ctx, walletID)
	}

	return nil, models.ErrInvalidWalletType
}

func (s *TransactionService) Confirm(ctx context.Context, walletType models.WalletType, transactionID int) error {
	if walletType == models.Personal || walletType == models.Family {
		return s.transactionStorage.Confirm(ctx, transactionID, models.Success)
	}

	return models.ErrInvalidWalletType
}

func (s *TransactionService) Cancel(ctx context.Context, walletType models.WalletType, transactionID int,
	credentials models.Credentials) error {
	if walletType == models.Personal || walletType == models.Family {
		return s.transactionStorage.Cancel(ctx, transactionID, models.Canceled, walletType, credentials)
	}

	return models.ErrInvalidWalletType
}
