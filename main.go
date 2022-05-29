package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello Vinicius!",
	})
}

func main() {
	// Criando router com Gin
	router := gin.Default()

	// gin.SetMode("release") // Opção para o Gin em produção

	// Criando arquivo de log do servidor
	f, err := os.Create("gin.log")
	if err != nil {
		log.Fatal(err)
	}

	gin.DefaultWriter = io.MultiWriter(f)
	// Capturar o request GET "/hello-world"
	router.GET("/hello", helloHandler)
	router.Run()

}
