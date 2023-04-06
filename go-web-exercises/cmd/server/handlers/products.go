package handler

import (
	"errors"
	"go-web-exercises/internal/domain"
	"go-web-exercises/internal/product"
	"go-web-exercises/pkg/web"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller
type ProductHandler struct {
	service product.Service
}

// Constructor
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := h.service.GetAll()
		web.SuccessfulResponse(http.StatusOK, products, c)
	}
}

func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, "Invalid id", c)
			return
		}
		product, err := h.service.GetById(id)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, "Invalid id", c)
			return
		}
		web.SuccessfulResponse(http.StatusOK, product, c)
	}
}

func (h *ProductHandler) GetPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		valueOf := c.Query("priceGt")
		value, err := strconv.ParseFloat(valueOf, 64)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, "Invalid product price", c)
			return
		}
		product, err := h.service.GetByPriceGt(value)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, err.Error(), c)
			return
		}
		web.SuccessfulResponse(http.StatusOK, product, c)
	}
}

func (h *ProductHandler) Create() gin.HandlerFunc {
	type request struct {
		Name         string  `json:"name"`
		Quantity     int     `json:"quantity"`
		Code_Value   int     `json:"code_value"`
		Is_Published bool    `json:"is_published"`
		Expiration   string  `json:"expiration"`
		Price        float64 `json:"price"`
	}

	return func(c *gin.Context) {
		tokenFromHeader := c.GetHeader("Token")
		tokenFromEnv := os.Getenv("TOKEN")

		if tokenFromEnv != tokenFromHeader {
			web.ErrorResponse(http.StatusUnauthorized, "Token is invalid", c)
			return
		}

		var req request
		//Obtains the new product form the request body
		if err := c.ShouldBindJSON(&req); err != nil {
			web.ErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		product := &domain.Product{
			Name:         req.Name,
			Quantity:     req.Quantity,
			Code_Value:   req.Code_Value,
			Is_Published: req.Is_Published,
			Expiration:   req.Expiration,
			Price:        req.Price,
		}
		//Create the new product
		newProduct, err := h.service.Create(*product)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}
		web.SuccessfulResponse(http.StatusCreated, newProduct, c)
	}
}

func (h *ProductHandler) Update() gin.HandlerFunc {
	type request struct {
		Name         string  `json:"name"`
		Quantity     int     `json:"quantity"`
		Code_Value   int     `json:"code_value"`
		Is_Published bool    `json:"is_published"`
		Expiration   string  `json:"expiration"`
		Price        float64 `json:"price"`
	}
	return func(c *gin.Context) {
		tokenFromHeader := c.GetHeader("Token")
		tokenFromEnv := os.Getenv("TOKEN")

		if tokenFromEnv != tokenFromHeader {
			web.ErrorResponse(http.StatusUnauthorized, "Token is invalid", c)
			return
		}

		var req request
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, "Invalid id", c)
			return
		}
		//var product domain.Product
		err = c.ShouldBindJSON(&req)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, "Invalid product", c)
			return
		}

		product := domain.Product{
			Name:         req.Name,
			Quantity:     req.Quantity,
			Code_Value:   req.Code_Value,
			Is_Published: req.Is_Published,
			Expiration:   req.Expiration,
			Price:        req.Price,
		}

		p, err := h.service.Update(id, product)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, err.Error(), c)
			return
		}
		web.SuccessfulResponse(http.StatusOK, p, c)
	}
}

func validateExpiration(date string) (bool, error) {
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return false, errors.New("invalid expiration date")
	}

	if err == nil && parsedDate.After(time.Now()) {
		return true, nil
	}
	return false, errors.New("invalid expiration date")
}

func (h *ProductHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Name         string  `json:"name,omitempty"`
		Quantity     int     `json:"quantity,omitempty"`
		Code_Value   int     `json:"code_value,omitempty"`
		Is_Published bool    `json:"is_published,omitempty"`
		Expiration   string  `json:"expiration,omitempty"`
		Price        float64 `json:"price,omitempty"`
	}
	return func(c *gin.Context) {
		tokenFromHeader := c.GetHeader("Token")
		tokenFromEnv := os.Getenv("TOKEN")

		if tokenFromEnv != tokenFromHeader {
			web.ErrorResponse(http.StatusUnauthorized, "Token is invalid", c)
			return
		}
		var r Request
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, "Invalid id", c)
			return
		}
		if err := c.ShouldBindJSON(&r); err != nil {
			web.ErrorResponse(http.StatusNotFound, "Invalid request", c)
			return
		}
		update := domain.Product{
			Name:         r.Name,
			Quantity:     r.Quantity,
			Code_Value:   r.Code_Value,
			Is_Published: r.Is_Published,
			Expiration:   r.Expiration,
			Price:        r.Price,
		}
		if update.Expiration != "" {
			valid, err := validateExpiration(update.Expiration)
			if !valid {
				web.ErrorResponse(http.StatusNotFound, err.Error(), c)
				return
			}
		}
		p, err := h.service.Update(id, update)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, err.Error(), c)
			return
		}
		web.SuccessfulResponse(http.StatusOK, p, c)
	}
}

func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenFromHeader := c.GetHeader("Token")
		tokenFromEnv := os.Getenv("TOKEN")

		if tokenFromEnv != tokenFromHeader {
			web.ErrorResponse(http.StatusUnauthorized, "Token is invalid", c)
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.ErrorResponse(http.StatusBadRequest, "invalid id", c)
			return
		}
		err = h.service.Delete(id)
		if err != nil {
			web.ErrorResponse(http.StatusNotFound, err.Error(), c)
			return
		}
		web.SuccessfulResponse(http.StatusOK, gin.H{"message": "product deleted"}, c)
	}
}
