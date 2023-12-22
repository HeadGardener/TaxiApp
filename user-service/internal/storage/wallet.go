package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotOwner  = errors.New("you are not the owner of this wallet")
	ErrNotMember = errors.New("you are not the member of this wallet")
)

type WalletStorage struct {
	db *sqlx.DB
}

func NewWalletStorage(db *sqlx.DB) *WalletStorage {
	return &WalletStorage{db: db}
}

// personal wallets

func (s *WalletStorage) Create(ctx context.Context, wallet *models.Wallet) (string, error) {
	var createWalletQuery = fmt.Sprintf(`INSERT INTO %s
    										(id, user_id, card, balance)
    										VALUES($1$2,$3,$4)`, walletsTable)

	if _, err := s.db.ExecContext(ctx,
		createWalletQuery,
		wallet.ID,
		wallet.UserID,
		wallet.Card,
		wallet.Balance); err != nil {
		return "", err
	}

	return wallet.ID, nil
}

func (s *WalletStorage) GetByID(ctx context.Context, walletID string) (models.Wallet, error) {
	var getByIDQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, walletsTable)

	var wallet models.Wallet

	if err := s.db.QueryRowContext(ctx, getByIDQuery, walletID).Scan(&wallet); err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (s *WalletStorage) GetAllWallets(ctx context.Context, userID string) ([]models.Wallet, error) {
	var getAllWalletsQuery = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1`, walletsTable)

	var wallets []models.Wallet

	if err := s.db.SelectContext(ctx, &wallets, getAllWalletsQuery, userID); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *WalletStorage) GetBalance(ctx context.Context, walletType models.WalletType,
	walletID string) (float64, error) {
	var getBalanceQuery string

	if walletType == models.Personal {
		getBalanceQuery = fmt.Sprintf(`SELECT (balance) FROM %s WHERE id=$1`, walletsTable)
	}

	if walletType == models.Family {
		getBalanceQuery = fmt.Sprintf(`SELECT (balance) FROM %s WHERE id=$1`, familyWalletsTable)
	}

	var balance float64

	if err := s.db.SelectContext(ctx, &balance, getBalanceQuery, walletID); err != nil {
		return 0.0, err
	}

	return balance, nil
}

func (s *WalletStorage) Refill(ctx context.Context, transaction *models.Transaction) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var refillQuery = fmt.Sprintf(`UPDATE %s SET balance=balance+$1 WHERE id=$2`, walletsTable)

	if _, err = tx.ExecContext(ctx, refillQuery, transaction.Money, transaction.WalletID); err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return 0, err
	}

	var createTransactionQuery = fmt.Sprintf(`INSERT INTO %s 
    											(wallet_id, money, transaction_type, transaction_status, date)
												VALUES ($1,$2,$3,$4,$5) RETURNING id`, transactionsTable)

	var id int

	if err = tx.QueryRowContext(ctx,
		createTransactionQuery,
		transaction.WalletID,
		transaction.Money,
		transaction.TransactionType,
		transaction.TransactionStatus,
		transaction.Date).Scan(&id); err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("unexpected error: unable to commit: %w", err)
	}

	return id, nil
}

// family wallets

func (s *WalletStorage) CreateFamilyWallet(ctx context.Context, userID string,
	wallet *models.FamilyWallet) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	var createFamilyWalletQuery = fmt.Sprintf(`INSERT INTO %s (id, wallet_id, balance, fixed_balance)
														VALUES($1,$2,$3,$4)`, familyWalletsTable)

	if _, err = tx.ExecContext(ctx,
		createFamilyWalletQuery,
		wallet.ID,
		wallet.WalletID,
		wallet.Balance,
		wallet.FixedBalance); err != nil {
		if err = tx.Rollback(); err != nil {
			return "", fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return "", err
	}

	var addMemberQuery = fmt.Sprintf(`INSERT INTO %s (user_id, wallet_id) VALUES($1,$2)`,
		usersWalletsTable)

	if _, err = tx.ExecContext(ctx, addMemberQuery, userID, wallet.ID); err != nil {
		if err = tx.Rollback(); err != nil {
			return "", fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("unexpected error: unable to commit: %w", err)
	}

	return wallet.ID, nil
}

func (s *WalletStorage) GetFamilyWalletByID(ctx context.Context, walletID string) (models.FamilyWallet, error) {
	var getFamilyByIDQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, familyWalletsTable)

	var wallet models.FamilyWallet

	if err := s.db.QueryRowContext(ctx, getFamilyByIDQuery, walletID).Scan(&wallet); err != nil {
		return models.FamilyWallet{}, err
	}

	return wallet, nil
}

func (s *WalletStorage) DeleteFamilyWallet(ctx context.Context, walletID string) error {
	var deleteFamilyQuery = fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, familyWalletsTable)

	if _, err := s.db.ExecContext(ctx, deleteFamilyQuery, walletID); err != nil {
		return err
	}

	return nil
}

func (s *WalletStorage) SetFixedBalance(ctx context.Context, walletID string, fixedBalance float64) error {
	var setFixedBalanceQuery = fmt.Sprintf(`UPDATE %s SET fixed_balance=$1 WHERE id=$2`, familyWalletsTable)

	if _, err := s.db.ExecContext(ctx, setFixedBalanceQuery, fixedBalance, walletID); err != nil {
		return err
	}

	return nil
}

func (s *WalletStorage) AddMember(ctx context.Context, walletID, memberID string) error {
	var addMemberQuery = fmt.Sprintf(`INSERT INTO %s (user_id, wallet_id) VALUES($1,$2)`,
		usersWalletsTable)

	if _, err := s.db.ExecContext(ctx, addMemberQuery, memberID, walletID); err != nil {
		return err
	}

	return nil
}

func (s *WalletStorage) GetMembers(ctx context.Context, walletID string) ([]models.User, error) {
	var getMembersIDsQuery = fmt.Sprintf(`SELECT (user_id) FROM %s WHERE wallet_id=$1`,
		familyWalletsTable)

	var ids []string

	if err := s.db.SelectContext(ctx, &ids, getMembersIDsQuery, walletID); err != nil {
		return nil, err
	}

	var getMembersQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=ANY($1)`, usersTable)

	var users []models.User

	if err := s.db.SelectContext(ctx, &users, getMembersQuery, ids); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *WalletStorage) DeleteMember(ctx context.Context, walletID, memberID string) error {
	var deleteMemberQuery = fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1 AND wallet_id=$2`,
		usersWalletsTable)

	if _, err := s.db.ExecContext(ctx, deleteMemberQuery, memberID, walletID); err != nil {
		return err
	}

	return nil
}

func (s *WalletStorage) GetAllFamilyWalletsAsOwner(ctx context.Context, userID string) ([]models.FamilyWallet, error) {
	var getAllFamilyWalletsAsOwnerQuery = fmt.Sprintf(`SELECT * FROM %s fw INNER JOIN %s w ON
    												fw.wallet_id=w.id AND w.user_id=$1`,
		familyWalletsTable, walletsTable)

	var familyWallets []models.FamilyWallet

	if err := s.db.SelectContext(ctx, &familyWallets, getAllFamilyWalletsAsOwnerQuery, userID); err != nil {
		return nil, err
	}

	return familyWallets, nil
}

func (s *WalletStorage) GetAllFamilyWalletsAsMember(ctx context.Context, userID string) ([]models.FamilyWallet, error) {
	var getAllFamilyWalletsAsMemberQuery = fmt.Sprintf(`SELECT (id, balance) FROM %s fw INNER JOIN %s uw ON
															fw.id=uw.wallet_id AND uw.user_id=$1`,
		familyWalletsTable, usersWalletsTable)

	var familyWallets []models.FamilyWallet

	if err := s.db.SelectContext(ctx, &familyWallets, getAllFamilyWalletsAsMemberQuery, userID); err != nil {
		return nil, err
	}

	return familyWallets, nil
}

func (s *WalletStorage) RefillFamilyWallet(ctx context.Context, walletID, famWalletID string, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var refillFamilyQuery = fmt.Sprintf(`UPDATE %s SET balance=balance+$1 WHERE id=$2`, familyWalletsTable)

	if _, err = tx.ExecContext(ctx, refillFamilyQuery, amount, famWalletID); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return err
	}

	var deductQuery = fmt.Sprintf(`UPDATE %s SET balance=balance-$1 WHERE id=$2`, walletsTable)

	if _, err = tx.ExecContext(ctx, deductQuery, amount, walletID); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unexpected error: unable to commit: %w", err)
	}

	return nil
}

func (s *WalletStorage) IsWalletOwner(ctx context.Context, walletID, userID string) error {
	var isOwnerQuery = fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE id=$1 AND user_id=$2`, walletsTable)

	var isOwner int

	if err := s.db.SelectContext(ctx, &isOwner, isOwnerQuery, walletID, userID); err != nil {
		return err
	}

	if isOwner == 0 {
		return ErrNotOwner
	}

	return nil
}

func (s *WalletStorage) IsFamilyWalletOwner(ctx context.Context, walletID, userID string) error {
	var isFamilyOwnerQuery = fmt.Sprintf(`SELECT COUNT(1) FROM %s fw WHERE fw.id=$1 
                           							INNER JOIN %s w ON fw.wallet_id=w.id AND w.user_id=$2`,
		familyWalletsTable, walletsTable)

	var isOwner int

	if err := s.db.SelectContext(ctx, &isOwner, isFamilyOwnerQuery, walletID, userID); err != nil {
		return err
	}

	if isOwner == 0 {
		return ErrNotOwner
	}

	return nil
}

func (s *WalletStorage) IsFamilyWalletMember(ctx context.Context, walletID, userID string) error {
	var isMemberQuery = fmt.Sprintf(`SELECT COUNT(1) FROM %s uw WHERE uw.wallet_id=$1 AND uw.user_id=$2
                           					INNER JOIN %s fw ON uw.wallet_id=fw.id`,
		usersWalletsTable, familyWalletsTable)

	var isMember int

	if err := s.db.SelectContext(ctx, &isMember, isMemberQuery, walletID, userID); err != nil {
		return err
	}

	if isMember == 0 {
		return ErrNotMember
	}

	return nil
}

func (s *WalletStorage) Withdraw(ctx context.Context, transactionID int, status string, walletType models.WalletType,
	credentials models.Credentials) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var withdrawQuery string

	if walletType == models.Personal {
		withdrawQuery = fmt.Sprintf(`UPDATE %s SET balance=balance-$1 WHERE id=$2`,
			walletsTable)
	}

	if walletType == models.Family {
		withdrawQuery = fmt.Sprintf(`UPDATE %s SET balance=balance-$1 WHERE id=$2`,
			familyWalletsTable)
	}

	if _, err = tx.ExecContext(ctx, withdrawQuery, credentials.Money, credentials.WalletID); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return err
	}

	var changeTransactionStatusQuery = fmt.Sprintf(`UPDATE %s SET transaction_status=$1 WHERE id=$2`,
		transactionsTable)

	if _, err = tx.ExecContext(ctx, changeTransactionStatusQuery, status, transactionID); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("unexpected error: unable to rollback: %w", err)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unexpected error: unable to commit: %w", err)
	}

	return nil
}
