package services

import (
	"kasir2-api/models"
	"kasir2-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Category) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
