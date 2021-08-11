package helper

//go:generate mockgen -source=transaction_manager.go -destination=./../mocks/transaction_manager_mock.go -package=mocks

import (
	"github.com/jmoiron/sqlx"
)

type TransactionManager interface {
	NewTransaction() Transaction
}

type transactionManager struct {
	db *sqlx.DB
}

func (transactionManager transactionManager) NewTransaction() Transaction {
	return NewTransaction(transactionManager.db.MustBegin())
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return transactionManager{db: db}
}
