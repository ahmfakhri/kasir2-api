package repositories

import (
	"database/sql"
	"kasir2-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	var report models.DailyReport

	// total revenue hari ini
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalRevenue)
	if err != nil {
		return nil, err
	}

	// total transaksi hari ini
	err = repo.db.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// produk terlaris
	err = repo.db.QueryRow(`
		SELECT c.name, SUM(td.quantity) AS qty
		FROM transaction_details td
		JOIN categories c ON c.id = td.category_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY c.name
		ORDER BY qty DESC
		LIMIT 1
	`).Scan(
		&report.ProdukTerlaris.Name,
		&report.ProdukTerlaris.QtyTerjual,
	)

	// kalau belum ada transaksi hari ini
	if err == sql.ErrNoRows {
		report.ProdukTerlaris = models.BestSeller{}
		return &report, nil
	}
	if err != nil {
		return nil, err
	}

	return &report, nil
}
