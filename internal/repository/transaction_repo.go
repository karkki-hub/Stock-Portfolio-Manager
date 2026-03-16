package repository

import (
	"database/sql"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	query := `
	INSERT INTO transactions (user_id, stock_id, transaction_type, qty, price) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(
		query,
		transaction.UserID,
		transaction.StockID,
		transaction.Type,
		transaction.Quantity,
		transaction.Price,
	)
	return err
}

func (r *TransactionRepository) GetByUserID(userID uint) ([]*models.Transaction, error) {
	query := `
	SELECT id, user_id, stock_id, transaction_type, qty, price, tot_amt, created_at
	FROM transactions WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.StockID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.TotalAmt,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetHolding(userID uint, stockID uint) (int, error) {
	query := `
	SELECT COALESCE(SUM(CASE WHEN transaction_type = 'buy' THEN qty END), 0) - COALESCE(SUM(CASE WHEN transaction_type = 'sell' THEN qty END), 0)
	FROM transactions WHERE user_id = ? AND stock_id = ?`

	var qty float64
	err := r.DB.QueryRow(query, userID, stockID).Scan(&qty)
	if err != nil {
		return 0, err
	}
	return int(qty), nil
}
