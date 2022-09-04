package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pmoule/go2hal/hal"
)

var validate = validator.New()

var encoder = hal.NewEncoder()

func WriteHAL(ctx *gin.Context, statusCode int, resource hal.Resource) {
	_, err := encoder.WriteTo(ctx.Writer, statusCode, resource)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func WriteNoContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}
