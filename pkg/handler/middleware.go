package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h Handler) middleware(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("recovered: %v", err)
		}
	}()
}
