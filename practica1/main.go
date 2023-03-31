package main 

import(
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
    "net/http"
	"strconv"
	
)

//Ejercicio 1
//Cargar en una slice, desde un archivo JSON, los datos de productos.
type Product struct{
	ID         	 int64 `json:"id"`
    Name       	 string `json:"name"`
	Quantity   	 int `json:"quantity"`
	Code_Value 	 int `json:"code_value"`
	Is_Published bool `json:"is_published"`
	Expiration   string `json:"expiration"`
	Price        float64 `json:"price"`

}
func loadJSONFile() []Product {
    var products []Product
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
	c.IndentedJSON(http.StatusOK, ProductsGlobal)
}

func ProductsById(c *gin.Context){
	id := c.Param("id")

    for _, p := range ProductsGlobal {
		value, _:= strconv.ParseInt(id, 10, 64)
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
    for _, p := range ProductsGlobal {
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
 
var ProductsGlobal = loadJSONFile()

func main(){

	server := gin.Default()
	//A. Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
	server.GET("/ping", PingPong)

	//B. Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
	gopher := server.Group("/products")
	gopher.GET("", Products)

    //C. Crear una ruta /products/:id que nos devuelva un producto por su id.
	gopher.GET("/:id", ProductsById)

    //D. Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
	gopher.GET("/search", SearchProductsByPrice)

	server.Run()

}