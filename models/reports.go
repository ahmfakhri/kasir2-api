package models

type BestSeller struct {
	Name       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type DailyReport struct {
	TotalRevenue   float64    `json:"total_revenue"`
	TotalTransaksi int        `json:"total_transaksi"`
	ProdukTerlaris BestSeller `json:"produk_terlaris"`
}
