package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"schneider.vip/problem"
)

func HandleError(ctx *gin.Context, err error) {
	var p *problem.Problem

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		p = problem.New(
			problem.Title("Record Not Found"),
			problem.Type("errors:database/record-not-found"),
			problem.Detail(err.Error()),
			problem.Status(http.StatusNotFound),
		)
	default:
		p = problem.New(
			problem.Title("Bad Request"),
			problem.Type("errors:http/bad-request"),
			problem.Detail(err.Error()),
			problem.Status(http.StatusBadRequest),
		)
	}

	if _, err := p.WriteTo(ctx.Writer); err != nil {
		panic(err)
	}
}
