package main

import (
	"event/event"
	"event/kafka"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Pusher struct {
	Sender event.Sender
}

type PushParams struct {
	fx.In

	Sender event.Sender
}

func NewPusher(params PushParams) *Pusher {
	return &Pusher{Sender: params.Sender}
}

func (p *Pusher) Mount(router gin.IRouter) {
	router.POST("/push", p.Push)
}

type PushRequest struct {
	Key string `json:"key" form:"key" binding:"required"`
	MSG string `json:"msg" form:"msg" binding:"required"`
}

func (p *Pusher) Push(ctx *gin.Context) {
	var validate PushRequest
	if err := ctx.ShouldBind(&validate); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := p.Sender.Send(ctx, kafka.NewMessage(validate.Key, []byte(validate.MSG)))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": validate.MSG})
}
