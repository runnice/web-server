package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/runnice/web-server/internal/users"
	"github.com/runnice/web-server/pkg/web"
)

type request struct {
	Id        int     `json:"id"`
	Name      string  `json:"name" binding:"required"`
	Lastname  string  `json:"lastname" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Age       int     `json:"age" binding:"required"`
	Height    float64 `json:"height" binding:"required"`
	IsActive  bool    `json:"isActive" binding:"required"`
	BirthDate string  `json:"birthDate" binding:"required"`
}

type putRequest struct {
	Name string `json:"name" binding:"required"`
}

type User struct {
	service users.Service
}

// Recebe um Service(Pacote Interno) e retorna o controller instanciado
func NewUser(u users.Service) *User {
	return &User{
		service: u,
	}
}

func isTokenValid(ctx *gin.Context) bool {
	token := ctx.GetHeader("token")
	if token != os.Getenv("TOKEN") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return true
	}
	return false

}

// Vai validar a solicitação, passar a tarefa ao Service e devolver resposta ao cliente
func (c *User) GetAll(ctx *gin.Context) {
	if isTokenValid(ctx) {
		return
	}
	u, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(500, web.NewResponse(500, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, u))
}

func (c *User) Store(ctx *gin.Context) {
	if isTokenValid(ctx) {
		return
	}
	var request request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(404, web.NewResponse(404, err.Error()))
		return
	}
	u, err := c.service.Store(request.Id, request.Age, request.Name, request.Lastname, request.Email, request.BirthDate, request.Height, request.IsActive)
	if err != nil {
		ctx.JSON(404, web.NewResponse(404, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, u))

}

func (c *User) Update(ctx *gin.Context) {
	if isTokenValid(ctx) {
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, web.NewResponse(400, err.Error()))
		return
	}

	var req request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		ctx.JSON(400, gin.H{"error": "O nome do usuário é obrigatório"})
		return
	}
	if req.Lastname == "" {
		ctx.JSON(400, gin.H{"error": "O sobrenome do usuário é obrigatório"})
		return
	}
	if req.Email == "" {
		ctx.JSON(400, gin.H{"error": "O e-mail é obrigatório"})
		return
	}
	if req.Age == 0 {
		ctx.JSON(400, gin.H{"error": "A idade é obrigatória"})
		return
	}
	if req.Height == 0.0 {
		ctx.JSON(400, gin.H{"error": "A altura é obrigatória"})
		return
	}
	if req.BirthDate == "" {
		ctx.JSON(400, gin.H{"error": "A data de nascimento é obrigatória"})
		return
	}
	u, err := c.service.Update(int(id), req.Age, req.Name, req.Lastname, req.Email, req.BirthDate, req.Height, req.IsActive)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, u))
}

func (c *User) UpdateName(ctx *gin.Context) {
	if isTokenValid(ctx) {
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid ID"})
		return
	}
	var req putRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		ctx.JSON(400, gin.H{"error": "O nome do usuário é obrigatório"})
		return
	}
	u, err := c.service.UpdateName(int(id), req.Name)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, u))
}

func (c *User) Delete(ctx *gin.Context) {
	if isTokenValid(ctx) {
		return
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid ID"})
		return
	}
	err = c.service.Delete(int(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Usuário deletado com sucesso"})
}
