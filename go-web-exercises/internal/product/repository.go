package product

import(
	"errors"
	"go-web-exercises/internal/domain"
	"time"
)

//Interface of the product service
type Repository interface {
	GetAll() []domain.Product
	GetById(id int) (domain.Product, error)
	GetByPriceGt(price float64) []domain.Product
	Create(product domain.Product) (domain.Product, error)
}

//Implementation of repository service
type RepositoryImpl struct {
	products []domain.Product
}

//Returns a new instance of the repository.
func NewRepository(products []domain.Product) Repository {
	return &RepositoryImpl{
		products: products,
	}
}

//Returns list of all available products 
func(r *RepositoryImpl) GetAll() []domain.Product{
	return r.products
}

//Returns a product filter by id
func (r *RepositoryImpl) GetById(id int) (domain.Product, error){
    for _, p := range r.products {
        if p.ID == id {
            return p, nil
        }
    }
    return domain.Product{}, errors.New("Product does not exist")
}

//Returns a list of products whith the price grater than the given price
func (r *RepositoryImpl) GetByPriceGt(price float64) []domain.Product{
    var totalProducts []domain.Product

    for _, p := range r.products {
        if p.Price > price {
			totalProducts = append(totalProducts, p)
		}
	  }
    return totalProducts 
}

func (r *RepositoryImpl) Create(newProduct domain.Product) (domain.Product, error){
	//Check valid expiration date 
	if !isValidDate(newProduct.Expiration){
		return domain.Product{}, errors.New("Product expiration date is invalid")
	}

	//Check valid code_value 
	if !r.isValidCodeValue(newProduct.Code_Value) {
		return domain.Product{}, errors.New("Product code cannot be duplicated")
	}
	/*8for _, prod := range r.products {
		if newProduct.Code_Value == prod.Code_Value{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product code cannot be duplicated "})
			return
		}
	}*/
	newProduct.ID = len(r.products) + 1
	r.products = append(r.products, newProduct)

	return newProduct, nil
}

//Check if the given Expiration date is in a correct format. 
func isValidDate(date string) bool{
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil{
		return false
	}

	if err == nil && parsedDate.After(time.Now()){
		return true
	}
	return false
}

func (r *RepositoryImpl) isValidCodeValue(codeValue int) bool{
	for _, prod := range r.products {
		if prod.Code_Value == codeValue{
			return false
		}
	}
	return true
}






