package usecase

import "jwt-go/internal/domain"

type ProductUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(r domain.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: r}
}

func (u *ProductUsecase) GetProducts() ([]domain.Product, error) {
	return u.repo.ListProducts()
}
