package main

import (
	"encoding/json"
	handler "go-web-exercises/cmd/server/handlers"
	"go-web-exercises/internal/domain"
	"go-web-exercises/internal/product"
	"go-web-exercises/pkg/store"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var products []domain.Product

func main() {

	_ = godotenv.Load()

	//Load json data file
	/*_,err := loadJSONFile()
	if err != nil{
		panic (err)
	}*/

	//Load json data file
	products := store.ReadJson()
	repository := product.NewRepository(products)
	service := product.NewService(repository)
	productHandler := handler.NewProductHandler(service)

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
		productGroup.PUT("/:id", productHandler.Update())
		productGroup.PATCH("/:id", productHandler.Patch())
		productGroup.DELETE("/:id", productHandler.Delete())
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
