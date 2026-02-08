package repositories

import (
	"database/sql"
	"errors"
	"kasir2-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// =====================
// GET ALL
// =====================
func (repo *CategoryRepository) GetAll(name string) ([]models.Category, error) {
	query := `SELECT id, name, description FROM categories`

	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// =====================
// CREATE
// =====================
func (repo *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`

	return repo.db.
		QueryRow(query, category.Name, category.Description).
		Scan(&category.ID)
}

// =====================
// GET BY ID
// =====================
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var c models.Category
	err := repo.db.QueryRow(query, id).
		Scan(&c.ID, &c.Name, &c.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// =====================
// UPDATE
// =====================
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`

	result, err := repo.db.Exec(
		query,
		category.Name,
		category.Description,
		category.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}

// =====================
// DELETE
// =====================
func (repo *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}
