package storage

import (
	"context"
	"fmt"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type TransactionStorage struct {
	db *sqlx.DB
}

func NewTransactionStorage(db *sqlx.DB) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) Create(ctx context.Context, transaction *models.Transaction) (int, error) {
	var createTransactionQuery = fmt.Sprintf(`INSERT INTO %s 
    											(wallet_id, money, transaction_type, transaction_status, date)
												VALUES ($1,$2,$3,$4,$5) RETURNING id`, transactionsTable)

	var id int

	if err := s.db.QueryRowContext(ctx,
		createTransactionQuery,
		transaction.WalletID,
		transaction.Money,
		transaction.TransactionType,
		transaction.TransactionStatus,
		transaction.Date).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TransactionStorage) GetAll(ctx context.Context, walletID string) ([]models.Transaction, error) {
	var getAllTransactionsQuery = fmt.Sprintf(`SELECT * FROM %s WHERE wallet_id=$1`, transactionsTable)

	var transactions []models.Transaction

	if err := s.db.SelectContext(ctx, &transactions, getAllTransactionsQuery, walletID); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionStorage) Confirm(ctx context.Context, transactionID int, status string) error {
	var changeTransactionStatusQuery = fmt.Sprintf(`UPDATE %s SET transaction_status=$1 WHERE id=$2`,
		transactionsTable)

	if _, err := s.db.ExecContext(ctx, changeTransactionStatusQuery, status, transactionID); err != nil {
		return err
	}

	return nil
}

//nolint:dupl
func (s *TransactionStorage) Cancel(ctx context.Context, transactionID int, status string, walletType models.WalletType,
	credentials models.Credentials) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var returnBalanceQuery string

	if walletType == models.Personal {
		returnBalanceQuery = fmt.Sprintf(`UPDATE %s SET balance=balance+$1 WHERE id=$2`, walletsTable)
	}

	if walletType == models.Family {
		returnBalanceQuery = fmt.Sprintf(`UPDATE %s SET balance=balance+$1 WHERE id=$2`, familyWalletsTable)
	}

	if _, err = tx.ExecContext(ctx, returnBalanceQuery, credentials.Money, credentials.WalletID); err != nil {
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
