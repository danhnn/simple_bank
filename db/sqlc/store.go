package db

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 1 for update")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 2 for update")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
