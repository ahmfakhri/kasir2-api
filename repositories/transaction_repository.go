package repositories

import (
	"database/sql"
	"fmt"
	"kasir2-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(
	items []models.CheckoutItem,
) (*models.Transaction, error) {

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var categoryName string

		// ambil category (produk)
		err := tx.QueryRow(
			`SELECT name FROM categories WHERE id = $1`,
			item.CategoryID,
		).Scan(&categoryName)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category id %d tidak ditemukan", item.CategoryID)
		}
		if err != nil {
			return nil, err
		}

		// sementara: subtotal = quantity (karena belum ada harga)
		subtotal := item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			CategoryID:   item.CategoryID,
			CategoryName: categoryName,
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
		})
	}

	// insert transaksi
	var transactionID int
	err = tx.QueryRow(
		`INSERT INTO transactions (total_amount)
		 VALUES ($1)
		 RETURNING id`,
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert detail
	for _, d := range details {
		_, err = tx.Exec(
			`INSERT INTO transaction_details
			(transaction_id, category_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4)`,
			transactionID, d.CategoryID, d.Quantity, d.Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
