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
}

type ServiceImpl struct {
	repository Repository
}

//Controller
/*func NewServiceProductLocal(db []*Product, lastID int) *ServiceProduct {
	return &ServiceProduct{
		db:     db,
		lastID: lastID,
	}
}*/

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

/*type ServiceProduct struct {
	db     []*Product
	lastID int
}*/

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

/*func (sv *ServiceProduct) Save(name string,  quantity int, code_value int, is_published bool, expiration string, price float64) (newProduct *Product, err error) {
	// instance
	newProduct = &Product{
		ID: sv.lastID + 1, 
		Name : name,
		Quantity : quantity,
		Code_Value : code_value,
		Is_Published : is_published, 
		Expiration : expiration,
		Price : price,
	}
	//Chequeo que la expiracion del producto sea valida
	if(!isValidDate(newProduct.Expiration)){
		err = ErrProductExpiration
		return
	}

	//Chequeo que el code value es unico
	for _, prod := range sv.db {
		if newProduct.Code_Value == prod.Code_Value{
			err = ErrProductExpiration
			return
		}
	}
	sv.db = append(sv.db, newProduct)
	sv.lastID++
	return 
}*/



