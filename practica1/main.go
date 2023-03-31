package main

import (
	"github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)


//Ejemplo
func HelloWorld(ctxt *gin.Context){
	ctxt.String(200, "Hello World!")
}

//Practica
//Ejercicio 1
func PingPong(ctxt *gin.Context){
	ctxt.String(200, "Pong")
}

//Ejercicio 2 
type Person struct {
    Name    string  `json:"nombre"`
    Surname  string  `json:"apellido"`
}

var p1 = Person{ Name : "Andrea", Surname: "Rivas"}

func Saludar(c *gin.Context){
	var s Person
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	mensaje := fmt.Sprintf("Hola %v %v", s.Name, s.Surname)

	c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})
}


func main(){
	server := gin.Default()
	server.GET("/hello-world", HelloWorld)
	server.GET("/ping", PingPong)
	server.POST("/saludar", Saludar) 
	server.Run()
}