package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skhaz/scheduler/model"
	"github.com/skhaz/scheduler/repository"
	"github.com/skhaz/scheduler/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type TriggerRepository struct {
	err      error
	trigger  *model.Trigger
	triggers model.TriggerCollection
	success  bool
}

func (r *TriggerRepository) Configure(db *gorm.DB) {
}

func (r *TriggerRepository) List(after time.Time, limit int) (any, error) {
	return r.triggers, r.err
}

func (r *TriggerRepository) Get(id any) (any, error) {
	return r.trigger, r.err
}

func (r *TriggerRepository) Create(entity any) (any, error) {
	return r.trigger, r.err
}

func (r *TriggerRepository) Update(id any, entity any) (bool, error) {
	return r.success, r.err
}

func (r *TriggerRepository) Delete(id any) (bool, error) {
	return r.success, r.err
}

type Workflow struct {
}

func (wf *Workflow) Apply(manifest []byte, op workflow.Operation) error { return nil }

func TestGetTriggers(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	triggers := model.TriggerCollection{{Name: randstr.String(16)}}
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &TriggerRepository{triggers: triggers}))

	GetTriggers(ctx)

	// namedMap := BuildMultipleResources("Triggers", Triggers).ToMap()
	// payload, _ := json.Marshal(namedMap.Content)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	// assert.Equal(t, payload, r.Body.Bytes())
}

func TestCreateTrigger(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)

	/*
			ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();not null" json:"id"`
		Name      string         `gorm:"type:varchar(32);not null" json:"name"`
		Schedule  string         `gorm:"type:varchar(32);not null" json:"schedule" validate:"cron"`
		Timezone  string         `gorm:"type:varchar(64);default:UTC;not null" json:"timezone" validate:"timezone"`
		Url       string         `gorm:"type:varchar(2048);not null" json:"url"`
		Method    string         `gorm:"type:varchar(8);not null" json:"method"`
		Success   int            `gorm:"type:smallint;default:200;not null" json:"success"`
		Timeout   int            `gorm:"type:smallint;default:60;not null" json:"timeout" validate:"gte=1,lte=300"`
		Retry     int            `gorm:"type:smallint;default:3;not null" json:"retry" validate:"gte=1,lte=10"`
		CreatedAt time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
		UpdatedAt time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index,->" json:"-"`*/
	trigger := model.Trigger{
		Name:     randstr.String(16),
		Schedule: "* * * * *",
		Timezone: "UTC",
		Url:      "https://httpbin.org/status/200",
		Timeout:  60,
		Retry:    3,
	}

	b, err := json.Marshal(trigger)
	assert.NoError(t, err)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/triggers", bytes.NewBuffer(b))

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &TriggerRepository{trigger: &trigger}))
	ctx.Set("Workflow", &Workflow{})

	CreateTrigger(ctx)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), fmt.Sprintf(`"name":"%v"`, trigger.Name))
}

func TestGetTrigger(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	trigger := model.Trigger{Name: randstr.String(16)}
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "f3edb291-a99d-4a43-8de0-1d6acd00c64d"}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &TriggerRepository{trigger: &trigger}))

	GetTrigger(ctx)

	// namedMap := BuildSingleResource(fmt.Sprintf("/Triggers/%v", uid), Trigger).ToMap()
	// payload, _ := json.Marshal(namedMap.Content)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	// assert.Equal(t, payload, r.Body.Bytes())
}

func TestDeleteTrigger(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	id := uuid.New()
	trigger := model.Trigger{ID: id}
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: id.String()}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &TriggerRepository{trigger: &trigger}))
	ctx.Set("Workflow", &Workflow{})

	DeleteTrigger(ctx)

	assert.Equal(t, http.StatusNoContent, r.Code)
	assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
}
