package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type BalanceProcessor struct {
	db *sqlx.DB
}

func NewBalanceProcessor(db *sqlx.DB) *BalanceProcessor {
	return &BalanceProcessor{db: db}
}

func (p *BalanceProcessor) Add(ctx context.Context, driverID string, money float64) error {
	if _, err := p.db.ExecContext(ctx, `UPDATE driver SET balance+=$1 WHERE id=$2`, money, driverID); err != nil {
		return err
	}

	return nil
}

func (p *BalanceProcessor) Withdraw(ctx context.Context, driverID string, money float64) error {
	if _, err := p.db.ExecContext(ctx, `UPDATE driver SET balance-=$1 WHERE id=$2`, money, driverID); err != nil {
		return err
	}

	return nil
}
