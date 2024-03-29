package main

import (
	"api"
	"data"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserService struct {
	connector   UserConnector
	userQueryer data.UserQueryer
}

func NewUserService(connector UserConnector, userQueryer data.UserQueryer) *UserService {
	return &UserService{connector: connector, userQueryer: userQueryer}
}

func (u *UserService) Mount(router gin.IRouter) {
	router.POST("/users", u.Create)
	router.GET("/users", u.List)
}

type CreateUserForm struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (u *UserService) Create(ctx *gin.Context) {
	var validate CreateUserForm
	if err := ctx.ShouldBind(&validate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := u.connector.CallServeUser(ctx, &api.UserRequest{Name: validate.Name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (u *UserService) List(ctx *gin.Context) {
	users, err := u.userQueryer.Select(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}
