package controller

import (
	_ "embed"

	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skhaz/scheduler/model"
	"github.com/skhaz/scheduler/repository"
	"github.com/skhaz/scheduler/workflow"
)

//go:embed manifest.yaml
var manifest string

type query struct {
	After time.Time `form:"after"`
	Limit int       `form:"limit,default=10" binding:"gte=1,lte=100"`
}

type params struct {
	ID string `uri:"uuid" validate:"required,uuid4"`
}

func GetTriggerRepository(ctx *gin.Context) repository.Repository {
	return ctx.MustGet("RepositoryRegistry").(*repository.RepositoryRegistry).MustRepository("TriggerRepository")
}

func GetWorkflow(ctx *gin.Context) workflow.Interface {
	return ctx.MustGet("Workflow").(workflow.Interface)
}

func GetManifest(trigger *model.Trigger) ([]byte, error) {
	tmp, err := template.New("manifest").Parse(manifest)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := tmp.Execute(&buffer, trigger); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GetTriggers(ctx *gin.Context) {
	var q = query{}

	if err := ctx.ShouldBindQuery(&q); err != nil {
		HandleError(ctx, err)
		return
	}

	e, err := GetTriggerRepository(ctx).List(q.After, q.Limit)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	WriteHAL(ctx, http.StatusOK, e.(model.TriggerCollection).ToHAL(ctx.Request.URL.Path, ctx.Request.URL.Query()))
}

func CreateTrigger(ctx *gin.Context) {
	body := model.Trigger{}

	if err := ctx.BindJSON(&body); err != nil {
		HandleError(ctx, err)
		return
	}

	if err := validate.Struct(body); err != nil {
		HandleError(ctx, err)
		return
	}

	e, err := GetTriggerRepository(ctx).Create(&body)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	trigger := e.(*model.Trigger)

	manifest, err := GetManifest(trigger)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	if err := GetWorkflow(ctx).Apply(manifest, workflow.Deploy); err != nil {
		HandleError(ctx, err)
		return
	}

	selfHref, _ := url.JoinPath(ctx.Request.URL.Path, trigger.ID.String())
	WriteHAL(ctx, http.StatusCreated, trigger.ToHAL(selfHref))
}

func GetTrigger(ctx *gin.Context) {
	p := params{}

	if err := ctx.ShouldBindUri(&p); err != nil {
		HandleError(ctx, err)
		return
	}

	if err := validate.Struct(p); err != nil {
		HandleError(ctx, err)
		return
	}

	e, err := GetTriggerRepository(ctx).Get(p.ID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	WriteHAL(ctx, http.StatusOK, e.(*model.Trigger).ToHAL(ctx.Request.URL.Path))
}

func DeleteTrigger(ctx *gin.Context) {
	p := params{}

	if err := ctx.ShouldBindUri(&p); err != nil {
		HandleError(ctx, err)
		return
	}

	if err := validate.Struct(p); err != nil {
		HandleError(ctx, err)
		return
	}

	repository := GetTriggerRepository(ctx)

	e, err := repository.Get(p.ID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	trigger := e.(*model.Trigger)

	manifest, err := GetManifest(trigger)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	if err := GetWorkflow(ctx).Apply(manifest, workflow.Displace); err != nil {
		HandleError(ctx, err)
		return
	}

	if _, err := repository.Delete(p.ID); err != nil {
		HandleError(ctx, err)
		return
	}

	WriteNoContent(ctx)
}
