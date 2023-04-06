package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	handler "go-web-exercises/cmd/server/handlers"
	"go-web-exercises/internal/domain"
	"go-web-exercises/internal/product"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func createServer(token string) *gin.Engine {

	if token != "" {
		err := os.Setenv("TOKEN", token)
		if err != nil {
			panic(err)
		}
	}

	products := readJson()
	repo := product.NewRepository(products)
	service := product.NewService(repo)
	productHandler := handler.NewProductHandler(service)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	pr := r.Group("/products")
	{
		pr.GET("", productHandler.GetAll())
		pr.GET("/:id", productHandler.GetById())
		pr.GET("/search", productHandler.GetPriceGt())
		pr.POST("", productHandler.Create())
		pr.PUT("/:id", productHandler.Update())
		pr.PATCH("/:id", productHandler.Patch())
		pr.DELETE("/:id", productHandler.Delete())
	}
	return r
}

func createRequestTest(method string, url string, body string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Token", token)
	}
	return req, httptest.NewRecorder()
}

func readJson() []domain.Product {
	var products []domain.Product
	jsonFile, err := os.Open("/Users/romsuarez/Documents/Practica Ejercicios/go-web/go-web-exercises/cmd/server/handlers/products_copy.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &products)
	return products
}

func writeSlice(slice []domain.Product) error {
	jsonFile, err := os.OpenFile("/Users/romsuarez/Documents/Practica Ejercicios/go-web/go-web-exercises/cmd/server/handlers/products_copy.json", os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		return errors.New("It was unable to open the file")
	}

	sliceToBytes, err := json.Marshal(slice)

	if err != nil {
		return errors.New("It was unable to write")
	}

	defer jsonFile.Close()
	jsonFile.Write(sliceToBytes)
	return nil
}

func TestProductsHandlers_GetAllProduct(t *testing.T) {
	t.Run("Test to obtain all products", func(t *testing.T) {
		//Arrange
		var expectd = response{Data: []domain.Product{}}
		r := createServer("token")
		p := readJson()
		expectd.Data = p
		actual := []domain.Product{}

		//Act
		req, rr := createRequestTest(http.MethodGet, "/products", "", "token")
		r.ServeHTTP(rr, req)

		//Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		err := json.Unmarshal(rr.Body.Bytes(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, expectd.Data, actual)
	})
}

func TestProductsHandlers_GetProductById(t *testing.T) {
	t.Run("Test to obtain an product", func(t *testing.T) {
		//Arrange
		var expectd = response{Data: domain.Product{
			ID:           1,
			Name:         "Oil - Margarine",
			Quantity:     439,
			Code_Value:   0,
			Is_Published: true,
			Expiration:   "15/12/2021",
			Price:        71.42,
		}}
		r := createServer("token")
		p := readJson()
		actual := p[0]

		//Act
		req, rr := createRequestTest(http.MethodGet, "/products/1", "", "token")
		r.ServeHTTP(rr, req)

		//Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		err := json.Unmarshal(rr.Body.Bytes(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, expectd.Data, actual)
	})
}

/*func TestProductsHandlers_PostProduct(t *testing.T) {
	t.Run("Test to create a product", func(t *testing.T) {
		//Arrange
		var expectd = response{Data: domain.Product{
			ID:           503,
			Name:         "Product insert by POST",
			Quantity:     439,
			Code_Value:   1300,
			Is_Published: false,
			Expiration:   "15/12/2023",
			Price:        71.42,
		}}

		product, _ := json.Marshal(expectd.Data)
		r := createServer("token")

		//Act
		req, rr := createRequestTest(http.MethodPost, "/products", string(product), "token")
		r.ServeHTTP(rr, req)
		actual := map[string]domain.Product{}
		_ = json.Unmarshal(rr.Body.Bytes(), &actual)

		//Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectd.Data, actual)
	})
}

func TestProductsHandlers_DeleteProduct(t *testing.T) {
	t.Run("Test delete a product", func(t *testing.T) {
		//Arrange
		r := createServer("token")
		p := readJson()
		//Act
		req, rr := createRequestTest(http.MethodDelete, "/products/500", "", "token")
		r.ServeHTTP(rr, req)
		_ = writeSlice(p)

		//Assert
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}*/

func TestProductsHandlers_BadRequest(t *testing.T) {
	t.Run("Testing bad request endpoints", func(t *testing.T) {
		//Arrange
		test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
		r := createServer("my-secret-token")
		//Act
		for _, method := range test {
			req, rr := createRequestTest(method, "/products/ABC", "{}", "my-secret-token")
			r.ServeHTTP(rr, req)
			//Assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})
}

func TestProductsHandlers_NotFound(t *testing.T) {
	t.Run("Testing not found endpoints", func(t *testing.T) {
		//Arrange
		test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
		r := createServer("my-secret-token")
		//Act
		for _, method := range test {
			req, rr := createRequestTest(method, "/products/1000", "{}", "my-secret-token")
			r.ServeHTTP(rr, req)
			//Assert
			assert.Equal(t, http.StatusNotFound, rr.Code)
		}
	})
}

func TestProductsHandlers_Unauthorized(t *testing.T) {
	t.Run("Testing unauthorized endpoints", func(t *testing.T) {
		//Arrange
		test := []string{http.MethodPut, http.MethodPatch, http.MethodDelete}
		r := createServer("my-secret-token")
		//Act
		for _, method := range test {
			req, rr := createRequestTest(method, "/products/1002", "{}", "not-my-token")
			r.ServeHTTP(rr, req)
			//Assert
			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		}
	})
}
func TestProductsHandlers_UnauthorizedPost(t *testing.T) {
	t.Run("Testing unauthorized post endpoint", func(t *testing.T) {
		//Arrange
		r := createServer("my-secret-token")
		//Act
		req, rr := createRequestTest(http.MethodPost, "/products", "{}", "not-my-token")
		r.ServeHTTP(rr, req)
		//Assert
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
