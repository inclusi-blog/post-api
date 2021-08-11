package helper

//go:generate mockgen -source=transaction.go -destination=./../mocks/transaction_mock.go -package=mocks
import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	SelectContext(goContext context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(goContext context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(goContext context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(goContext context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Rebind(query string) string
	Commit() error
	Rollback() error
}

type transaction struct {
	context *sqlx.Tx
}

func (transaction transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return transaction.context.QueryRowContext(ctx, query, args...)
}

func (transaction transaction) SelectContext(goContext context.Context, dest interface{}, query string, args ...interface{}) error {
	//nolint:safesql
	return transaction.context.SelectContext(goContext, dest, query, args...)
}

func (transaction transaction) GetContext(goContext context.Context, dest interface{}, query string, args ...interface{}) error {
	//nolint:safesql
	return transaction.context.GetContext(goContext, dest, query, args...)
}

func (transaction transaction) ExecContext(goContext context.Context, query string, args ...interface{}) (sql.Result, error) {
	//nolint:safesql
	return transaction.context.ExecContext(goContext, query, args...)
}

func (transaction transaction) NamedExecContext(goContext context.Context, query string, arg interface{}) (sql.Result, error) {
	//nolint:safesql
	return transaction.context.NamedExecContext(goContext, query, arg)
}

func (transaction transaction) Rebind(query string) string {
	//nolint:safesql
	return transaction.context.Rebind(query)
}

func (transaction transaction) Commit() error {
	//nolint:safesql
	return transaction.context.Commit()
}

func (transaction transaction) Rollback() error {
	//nolint:safesql
	return transaction.context.Rollback()
}

func NewTransaction(context *sqlx.Tx) Transaction {
	//nolint:safesql
	return transaction{context: context}
}
