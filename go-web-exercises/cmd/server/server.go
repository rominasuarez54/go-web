package main

import (
	"go-web-exercises/cmd/server/handlers"
	"go-web-exercises/internal/domain"
	"go-web-exercises/internal/product"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"net/http"
)
var products []domain.Product
func main() {

	//Load json data file
	_,err := loadJSONFile()
	if err != nil{
		panic (err)
	}

	repository := product.NewRepository(products)
	service := product.NewService(repository)
	productHandler := handlers.NewProductHandler(service)

	router := gin.Default()

	//Ping endpoint 
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//Products endpoint
	productGroup := router.Group("/products")
	{
		productGroup.GET("", productHandler.GetAll())
		productGroup.GET("/:id", productHandler.GetById())
		productGroup.GET("/search", productHandler.GetPriceGt())
		productGroup.POST("", productHandler.Create())
		productGroup.PUT(":id", productHandler.Update())
		productGroup.PATCH(":id", productHandler.Patch())
	}

	// Run
	if err := router.Run(); err != nil {
		panic(err)
	}
}

func loadJSONFile() ([]domain.Product, error) {
    jsonData, err := ioutil.ReadFile("/Users/romsuarez/Documents/Practica Ejercicios/go-web/go-web-exercises/products.json")
    if err != nil {
        return nil, err
    }
    json.Unmarshal([]byte(jsonData), &products)
    return products, nil
}