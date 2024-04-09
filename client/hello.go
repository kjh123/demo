package main

import (
	"api"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloService struct {
	connector HelloConnector
}

func NewHelloService(connector HelloConnector) *HelloService {
	return &HelloService{
		connector: connector,
	}
}

func (h *HelloService) Mount(router gin.IRouter) {
	router.GET("/hello", h.Hello)
}

func (h *HelloService) Hello(ctx *gin.Context) {
	r, err := h.connector.CallServer(ctx, &api.HelloRequest{
		Name: ctx.Query("name"),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": r.GetMessage()})
}
