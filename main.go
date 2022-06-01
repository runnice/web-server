package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Id        int     `json:"id"`
	Name      string  `json:"name" binding:"required"`
	Lastname  string  `json:"lastname" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Age       int     `json:"age" binding:"required"`
	Height    float64 `json:"height" binding:"required"`
	IsActive  bool    `json:"isActive" binding:"required"`
	BirthDate string  `json:"birthDate" binding:"required"`
}

type Users []User

var lastID = 2

type ErrorMsg struct {
	Message string `json:"message"`
}

func errorHandler(fe validator.FieldError) string {
	return fmt.Sprintf("field %s is required", fe.Field())
}

var users Users = readJson("./users.json")

func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello Vinicius!",
	})
}

func getAllWithJson(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getAll(c *gin.Context) {

	jsonDataRaw := `[
		{
		"id": 1,
		"name": "Joao",
		"lastmame": "Silva",
		"email": "joaosilva@gemail.com",
		"age": 20,
		"height": 1.80,
		"isActive": true,
		"birthDate": "01-01-2020"},
	{
		"id":2,
		"name":"Joana",
		"lastname": "Guedes",
		"email": "joana@gemail.com",
		"age": 35,
		"height": 1.85,
		"isActive": true,
		"birthDate": "01/01/1990"
	}
	]`

	var u []User

	if err := json.Unmarshal([]byte(jsonDataRaw), &u); err != nil {
		fmt.Println(err)
	}
	c.JSON(200, u)

}

// Exercício 02 - GoWeb-1 Aula 2
func filterByAge(fileName string, age int) Users {
	// arquivo, err := os.Open(fileName)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer arquivo.Close()

	// byteValue, _ := ioutil.ReadAll(arquivo)

	// var u Users

	// if err := json.Unmarshal(byteValue, &u); err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// // m := make(map[string]interface{})

	// u := readJson(fileName)
	var u2 Users

	for _, v := range users {
		if v.Age == age {
			u2 = append(u2, v)
		}

	}
	return u2

}

func readJson(fileName string) Users {
	arquivo, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer arquivo.Close()

	byteValue, _ := ioutil.ReadAll(arquivo)

	var u Users

	if err := json.Unmarshal(byteValue, &u); err != nil {
		fmt.Println("Error: ", err)
	}

	// m := make(map[string]interface{})

	return u

}

// Exercício 2

func GetId(c *gin.Context) {
	id := c.Param("id")
	convertId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "id must be a number"})
		log.Println(err)

	}
	// users := readJson("./users.json")
	for _, users := range users {
		if users.Id == convertId {
			c.JSON(http.StatusOK, gin.H{
				"user": users,
			})
			return

		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user not found",
	})

}

func filter(c *gin.Context) {

	c.JSON(200, filterByAge("./users.json", 20))

}

// func listUsers(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"data": &users})
// }

func validationToken(c *gin.Context) bool {
	token := c.GetHeader("token")
	if token != "12345" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "token invalid",
		})
		return true
	}
	return false

}

func createUsers(c *gin.Context) {
	var usuario User
	if validationToken(c) {
		return
	}
	if err := c.ShouldBindJSON(&usuario); err != nil {
		var validador validator.ValidationErrors
		if errors.As(err, &validador) {
			out := make([]ErrorMsg, len(validador))
			for i, fe := range validador {
				out[i] = ErrorMsg{errorHandler(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": out,
			})
		}
		return
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": err.Error(),
		// })
	}

	lastID++
	usuario.Id = lastID

	users = append(users, usuario)

	c.JSON(http.StatusOK, gin.H{
		"data": usuario,
	})

}
func main() {

	gin.SetMode("debug")
	router := gin.Default()
	router.GET("/", helloHandler)
	router.GET("/users", getAllWithJson)
	router.GET("/user", getAll)
	router.GET("/filter", filter)
	router.GET("/:id", GetId)

	group := router.Group("/usuarios")
	{
		group.GET("/:id", GetId)
		group.GET("/", getAllWithJson)
		group.POST("/", createUsers)
	}
	router.Run()
}
