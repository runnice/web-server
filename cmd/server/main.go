package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/runnice/web-server/cmd/server/controllers"
	"github.com/runnice/web-server/internal/users"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo := users.NewRepository()
	service := users.NewService(repo)
	u := controllers.NewUser(service)

	r := gin.Default()

	ur := r.Group("/users")
	ur.GET("/", u.GetAll)
	ur.POST("/", u.Store)
	ur.PUT("/:id", u.Update)
	ur.PATCH("/:id", u.UpdateName)
	ur.DELETE("/:id", u.Delete)
	r.Run()
}
