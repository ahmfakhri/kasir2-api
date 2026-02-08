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

	// ===== BEGIN TRANSACTION =====
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// ===== VALIDASI & HITUNG =====
	for _, item := range items {
		var categoryName string

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

		// OPSI 1: belum ada harga → subtotal = quantity
		subtotal := item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			CategoryID:   item.CategoryID,
			CategoryName: categoryName,
			Quantity:     item.Quantity,
			Subtotal:     subtotal, // ⬅️ PENTING (hindari error int vs float64)
		})
	}

	// ===== INSERT TRANSACTION =====
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

	// ===== INSERT DETAILS (PAKAI INDEX) =====
	for i := range details {
		details[i].TransactionID = transactionID

		_, err = tx.Exec(
			`INSERT INTO transaction_details
			 (transaction_id, category_id, quantity, subtotal)
			 VALUES ($1, $2, $3, $4)`,
			transactionID,
			details[i].CategoryID,
			details[i].Quantity,
			details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	// ===== COMMIT =====
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// ===== RETURN RESULT =====
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
