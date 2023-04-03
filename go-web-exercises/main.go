package main 

import(
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
    "net/http"
	"strconv"
	"time"
)

var products []Product
var autoId = 0
 
//Ejercicio 1
//Cargar en una slice, desde un archivo JSON, los datos de productos.
type Product struct{
	ID         	 int `json:"id"`
    Name       	 string `json:"name"`
	Quantity   	 int `json:"quantity"`
	Code_Value 	 int `json:"code_value"`
	Is_Published bool `json:"is_published"`
	Expiration   string `json:"expiration"`
	Price        float64 `json:"price"`

}
func loadJSONFile() []Product {
    jsonData, err := ioutil.ReadFile("./products.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    json.Unmarshal([]byte(jsonData), &products)
    return products
}

func PingPong(c *gin.Context){
	mensaje := fmt.Sprintf("Pong - status 200 OK")
	c.JSON(http.StatusOK, gin.H{"mensaje": mensaje })
}

func Products(c *gin.Context){
	c.IndentedJSON(http.StatusOK, products)
}

func ProductsById(c *gin.Context){
	id := c.Param("id")
	value, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
    for _, p := range products {
		
        if p.ID == value {
            c.IndentedJSON(http.StatusOK, p)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "El producto no existe "})

}

func SearchProductsByPrice(c *gin.Context){
	var totalProducts []Product

	valueOf := c.Query("priceGt")
    for _, p := range products {
		value, _ := strconv.ParseFloat(valueOf, 64)
        if p.Price > value {
			totalProducts = append(totalProducts, p)
			c.IndentedJSON(http.StatusOK, totalProducts)
		}
	  }
	  if(len(totalProducts)==0){
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No existen productos"})
	  }
    }

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

func CreateProduct() gin.HandlerFunc{
	type request struct {
		Name       	 string `json:"name"`
		Quantity   	 int `json:"quantity"`
		Code_Value 	 int `json:"code_value"`
		Is_Published bool `json:"is_published"`
		Expiration   string `json:"expiration"`
		Price        float64 `json:"price"`
	}
	return func (c *gin.Context){
		var req request 
		//var newProduct Product
		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		newProduct := &Product{
			ID: autoId + 1, 
			Quantity : req.Quantity,
			Code_Value : req.Code_Value,
			Is_Published : req.Is_Published, 
			Expiration : req.Expiration,
			Price : req.Price,
		}
		//Chequeo que la expiracion del producto sea valida
		if(!isValidDate(newProduct.Expiration)){
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product expiration date is invalid"})
			return
		}

		//Chequeo que el code value es unico
		for _, prod := range products {
			if newProduct.Code_Value == prod.Code_Value{
				c.JSON(http.StatusBadRequest, gin.H{"error": "Product code cannot be duplicated "})
				return
			}
		}
		products = append(products, *newProduct)
		c.JSON(http.StatusCreated, newProduct)
		}
	}
 
func main(){

	loadJSONFile()
	autoId = products[len(products)-1].ID

	server := gin.Default()
	//A. Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
	server.GET("/ping", PingPong)

	//B. Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
	productGroup := server.Group("/products")
	productGroup.GET("", Products)

    //C. Crear una ruta /products/:id que nos devuelva un producto por su id.
	productGroup.GET("/:id", ProductsById)

    //D. Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
	productGroup.GET("/search", SearchProductsByPrice)

	productGroup.POST("", CreateProduct())

	server.Run()

}