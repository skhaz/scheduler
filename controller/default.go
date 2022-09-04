package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"schneider.vip/problem"
)

func NoRoute(ctx *gin.Context) {
	p := problem.New(
		problem.Title("Not Found"),
		problem.Type("errors:http/not-found"),
		problem.Status(http.StatusNotFound),
	)

	if _, err := p.WriteTo(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
