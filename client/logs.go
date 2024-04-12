package main

import (
	"data/logs"

	"github.com/gin-gonic/gin"
)

type LogService struct {
	logWrite logs.LogWriter
}

func NewLogService(influxWrite logs.LogWriter) *LogService {
	return &LogService{logWrite: influxWrite}
}

func (i *LogService) Mount(router gin.IRouter) {
	router.POST("/logs", i.Write)
}

type BehaviorLogRequest logs.BehaviorLog

func (i *LogService) Write(ctx *gin.Context) {
	var req BehaviorLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := i.logWrite.Writer(ctx, logs.BehaviorLog{
		UID:  req.UID,
		IP:   req.IP,
		Tags: req.Tags,
		UA:   req.UA,
	}); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": true})
}
