package services

import (
	"context"
	"errors"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotEnoughBalance     = errors.New("not enough balance")
	ErrAddedMember          = errors.New("member was already added")
	ErrInvalidMember        = errors.New("this user is not the member")
	ErrDeleteOwner          = errors.New("you can't delete your self, you are the owner")
	ErrNotConnectedToWallet = errors.New("you not the owner or member of this wallet")
)

type WalletStorage interface {
	Create(ctx context.Context, wallet *models.Wallet) (string, error)
	GetByID(ctx context.Context, walletID string) (models.Wallet, error)
	GetAllWallets(ctx context.Context, userID string) ([]models.Wallet, error)
	GetBalance(ctx context.Context, walletType models.WalletType, walletID string) (float64, error)
	Refill(ctx context.Context, transaction *models.Transaction) (int, error)

	CreateFamilyWallet(ctx context.Context, userID string, wallet *models.FamilyWallet) (string, error)
	GetFamilyWalletByID(ctx context.Context, walletID string) (models.FamilyWallet, error)
	DeleteFamilyWallet(ctx context.Context, walletID string) error
	SetFixedBalance(ctx context.Context, walletID string, fixedBalance float64) error
	AddMember(ctx context.Context, walletID, memberID string) error
	GetMembers(ctx context.Context, walletID string) ([]models.User, error)
	DeleteMember(ctx context.Context, walletID, memberID string) error
	GetAllFamilyWalletsAsOwner(ctx context.Context, userID string) ([]models.FamilyWallet, error)
	GetAllFamilyWalletsAsMember(ctx context.Context, userID string) ([]models.FamilyWallet, error)
	RefillFamilyWallet(ctx context.Context, walletID, famWalletID string, amount float64) error

	IsWalletOwner(ctx context.Context, walletID, userID string) error
	IsFamilyWalletOwner(ctx context.Context, walletID, userID string) error
	IsFamilyWalletMember(ctx context.Context, walletID, userID string) error

	Withdraw(ctx context.Context, transactionID int, status string, walletType models.WalletType,
		credentials models.Credentials) error
}

type UserProvider interface {
	GetByPhone(ctx context.Context, phone string) (models.User, error)
}

type WalletService struct {
	walletStorage WalletStorage
	userProvider  UserProvider
}

func NewWalletService(walletStorage WalletStorage, userProvider UserProvider) *WalletService {
	return &WalletService{
		walletStorage: walletStorage,
		userProvider:  userProvider,
	}
}

// personal wallets

func (s *WalletService) Create(ctx context.Context, userID, card string) (string, error) {
	wallet := &models.Wallet{
		ID:      uuid.NewString(),
		UserID:  userID,
		Card:    card,
		Balance: 0.0,
	}

	return s.walletStorage.Create(ctx, wallet)
}

// TopUp method for refilling personal wallet balance
func (s *WalletService) TopUp(ctx context.Context, userID, walletID string,
	money float64) (int, error) {
	if err := s.walletStorage.IsWalletOwner(ctx, walletID, userID); err != nil {
		return 0, err
	}

	transaction, err := models.BuildTransaction(models.Refill, models.Success, money)
	if err != nil {
		return 0, err
	}

	transaction.WalletID = walletID

	return s.walletStorage.Refill(ctx, &transaction)
}

func (s *WalletService) GetByID(ctx context.Context, userID, walletID string) (models.Wallet, error) {
	if err := s.walletStorage.IsWalletOwner(ctx, walletID, userID); err != nil {
		return models.Wallet{}, err
	}

	return s.walletStorage.GetByID(ctx, walletID)
}

func (s *WalletService) ViewAll(ctx context.Context, userID string) ([]models.Wallet, error) {
	return s.walletStorage.GetAllWallets(ctx, userID)
}

// family wallets

func (s *WalletService) CreateFamilyWallet(ctx context.Context, userID, walletID string) (string, error) {
	if err := s.walletStorage.IsWalletOwner(ctx, walletID, userID); err != nil {
		return "", err
	}

	wallet := models.FamilyWallet{
		ID:           uuid.NewString(),
		WalletID:     walletID,
		Balance:      0.0,
		FixedBalance: 0.0,
	}

	return s.walletStorage.CreateFamilyWallet(ctx, userID, &wallet)
}

func (s *WalletService) GetFamilyWalletByID(ctx context.Context, userID, walletID string) (models.FamilyWallet, error) {
	wallet, err := s.walletStorage.GetFamilyWalletByID(ctx, walletID)
	if err != nil {
		return models.FamilyWallet{}, err
	}

	if err = s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err == nil {
		return wallet, nil
	}

	err = s.walletStorage.IsFamilyWalletMember(ctx, walletID, userID)
	if err == nil {
		wallet.WalletID = ""
		return wallet, nil
	}

	return models.FamilyWallet{}, ErrNotConnectedToWallet
}

func (s *WalletService) DeleteFamilyWallet(ctx context.Context, userID, walletID string) error {
	if err := s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err != nil {
		return err
	}

	return s.walletStorage.DeleteFamilyWallet(ctx, walletID)
}

func (s *WalletService) SetFixedBalance(ctx context.Context, userID, walletID string, fixedBalance float64) error {
	if err := s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err != nil {
		return err
	}

	return s.walletStorage.SetFixedBalance(ctx, walletID, fixedBalance)
}

func (s *WalletService) AddFamilyBalance(ctx context.Context, userID, walletID, famWalletID string,
	amount float64) error {
	if err := s.walletStorage.IsWalletOwner(ctx, walletID, userID); err != nil {
		return err
	}

	if err := s.walletStorage.IsFamilyWalletOwner(ctx, famWalletID, userID); err != nil {
		return err
	}

	balance, err := s.walletStorage.GetBalance(ctx, models.Family, walletID)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrNotEnoughBalance
	}

	return s.walletStorage.RefillFamilyWallet(ctx, walletID, famWalletID, amount)
}

func (s *WalletService) AddMember(ctx context.Context, userID, walletID, phone string) error {
	if err := s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err != nil {
		return err
	}

	member, err := s.userProvider.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}

	if err = s.walletStorage.IsFamilyWalletMember(ctx, walletID, member.ID); err == nil {
		return ErrAddedMember
	}

	return s.walletStorage.AddMember(ctx, walletID, member.ID)
}

func (s *WalletService) ViewMembers(ctx context.Context, userID, walletID string) ([]string, error) {
	if err := s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err != nil {
		return nil, err
	}

	members, err := s.walletStorage.GetMembers(ctx, walletID)
	if err != nil {
		return nil, err
	}

	phones := make([]string, len(members))
	for i := range members {
		phones[i] = members[i].Phone
	}

	return phones, nil
}

func (s *WalletService) DeleteMember(ctx context.Context, userID, walletID, phone string) error {
	if err := s.walletStorage.IsFamilyWalletOwner(ctx, walletID, userID); err != nil {
		return err
	}

	member, err := s.userProvider.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}

	if err = s.walletStorage.IsFamilyWalletMember(ctx, walletID, member.ID); err != nil {
		return ErrInvalidMember
	}

	if userID == member.ID {
		return ErrDeleteOwner
	}

	return s.walletStorage.DeleteMember(ctx, walletID, member.ID)
}

func (s *WalletService) ViewAllFamily(ctx context.Context, userID string) ([]models.FamilyWallet, error) {
	return s.walletStorage.GetAllFamilyWalletsAsOwner(ctx, userID)
}

func (s *WalletService) ViewMemberships(ctx context.Context, userID string) ([]models.FamilyWallet, error) {
	return s.walletStorage.GetAllFamilyWalletsAsMember(ctx, userID)
}

// pay section

func (s *WalletService) PickWalletAndPay(ctx context.Context, userID string, walletType models.WalletType,
	transactionID int, credentials models.Credentials) error {
	if walletType == models.Personal {
		if err := s.walletStorage.IsWalletOwner(ctx, credentials.WalletID, userID); err != nil {
			return err
		}

		balance, err := s.walletStorage.GetBalance(ctx, walletType, credentials.WalletID)
		if err != nil {
			return err
		}

		if balance < credentials.Money {
			return ErrNotEnoughBalance
		}

		return s.walletStorage.Withdraw(ctx, transactionID, models.Blocked, walletType, credentials)
	}

	if walletType == models.Family {
		if err := s.walletStorage.IsFamilyWalletMember(ctx, credentials.WalletID, userID); err != nil {
			return err
		}

		balance, err := s.walletStorage.GetBalance(ctx, walletType, credentials.WalletID)
		if err != nil {
			return err
		}

		if balance < credentials.Money {
			return ErrNotEnoughBalance
		}

		return s.walletStorage.Withdraw(ctx, transactionID, models.Blocked, walletType, credentials)
	}

	return models.ErrInvalidWalletType
}
