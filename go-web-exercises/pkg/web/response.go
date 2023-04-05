package web

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
 }
 type response struct {
	Data interface{} `json:"data"`
 }
 
 func ErrorResponse(status int, message string, ctx *gin.Context)  {
	ctx.JSON(status, gin.H{
		"code": strconv.Itoa(status),
		"message": message,
	})
}

func SuccessfulResponse(status int, data interface{}, ctx *gin.Context)  {
	ctx.JSON(status, data)
}