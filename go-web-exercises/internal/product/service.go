package product

import (

	"go-web-exercises/internal/domain"
	"errors"
)

type Service interface {
	GetAll() []domain.Product
	GetById(id int) (domain.Product, error)
	GetByPriceGt(price float64) ([]domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
	Update(id int, p domain.Product) (domain.Product, error)
}

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &ServiceImpl{
		repository : repository,
	}
}

//Defino errores
var (
	ErrProductExpiration = errors.New("Product expiration date is invalid")
	ErrDuplicatedProduct= errors.New("Product code cannot be duplicated")
)

func (s *ServiceImpl) GetAll() []domain.Product {
	return s.repository.GetAll()
}

func (s *ServiceImpl) GetById(id int) (domain.Product, error) {
	product, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *ServiceImpl) GetByPriceGt(price float64) ([]domain.Product, error) {
	products := s.repository.GetByPriceGt(price)
	if len(products) == 0 {
		return []domain.Product{}, errors.New("product does not exist")
	}
	return products, nil
}

func (s *ServiceImpl) Create(product domain.Product) (domain.Product, error) {
	newProduct, err := s.repository.Create(product)
	if err != nil {
		return domain.Product{}, err
	}
	return newProduct, nil
}

func (s *ServiceImpl) Update(id int, u domain.Product) (domain.Product, error) {
	p, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}
	if u.Name != "" {
		p.Name = u.Name
	}
	if u.Expiration != "" {
		p.Expiration = u.Expiration
	}
	if u.Quantity > 0 {
		p.Quantity = u.Quantity
	}
	if u.Price > 0 {
		p.Price = u.Price
	}
	p, err = s.repository.Update(id, p)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}
