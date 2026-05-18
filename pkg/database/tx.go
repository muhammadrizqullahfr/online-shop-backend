package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Gunakan key unik untuk menyimpan TX di context
type txKey struct{}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	GetTx(ctx context.Context) *gorm.DB
}

type TransactionManagerImpl struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManagerImpl {
	return &TransactionManagerImpl{db: db}
}

func (t *TransactionManagerImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Masukkan TX ke dalam context
	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			fmt.Println("ROLLBACK")
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = fn(txCtx)
	return
}

// GetTx digunakan oleh Repository untuk mengambil TX jika ada,
// jika tidak ada, gunakan DB standar.
func (t *TransactionManagerImpl) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return t.db.WithContext(ctx)
}
