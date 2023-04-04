package handlers

import (
	"go-web-exercises/internal/domain"
	"go-web-exercises/internal/product"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

//Controller
type ProductHandler struct {
	service product.Service
}
//Constructor
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		service : service,
	}
}

func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := h.service.GetAll()
		c.JSON(http.StatusOK, products)
	}
}

func (h *ProductHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
		product, err := h.service.GetById(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func (h *ProductHandler) GetPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		valueOf := c.Query("priceGt")
		value, err := strconv.ParseFloat(valueOf, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Invalid product price")})
			return
		}
		product, err := h.service.GetByPriceGt(value)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func (h *ProductHandler) Create() gin.HandlerFunc {
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

		//Obtains the new product form the request body
		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product := &domain.Product{
			Name:  req.Name,
			Quantity:   req.Quantity,
			Code_Value:  req.Code_Value,
			Is_Published: req.Is_Published,
			Expiration: req.Expiration,
			Price: req.Price,
		}
		//Create the new product
		newProduct, err := h.service.Create(*product)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		c.JSON(http.StatusCreated,  newProduct)
	}
}

func (h *ProductHandler) Update() gin.HandlerFunc {
	type request struct {
		Name       	 string `json:"name"`
		Quantity   	 int `json:"quantity"`
		Code_Value 	 int `json:"code_value"`
		Is_Published bool `json:"is_published"`
		Expiration   string `json:"expiration"`
		Price        float64 `json:"price"`
	}
	return func(ctx *gin.Context) {
		var req request
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		//var product domain.Product
		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid product"})
			return
		}

		product := domain.Product{
			Name:  req.Name,
			Quantity:   req.Quantity,
			Code_Value:  req.Code_Value,
			Is_Published: req.Is_Published,
			Expiration: req.Expiration,
			Price: req.Price,
		}
		
		p, err := h.service.Update(id, product)
		if err != nil {
			ctx.JSON(409, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
	
}