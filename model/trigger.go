package model

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pmoule/go2hal/hal"
	"gorm.io/gorm"
)

const (
	After        = "after"
	Limit        = "limit"
	NextRelation = "next"
)

type Trigger struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null" json:"id"`
	Name     string    `gorm:"type:varchar(32);not null" json:"name"`
	Schedule string    `gorm:"type:varchar(32);not null" json:"schedule" validate:"cron"`
	Timezone string    `gorm:"type:varchar(64);default:UTC;not null" json:"timezone" validate:"timezone"`
	Url      string    `gorm:"type:varchar(2048);not null" json:"url"`
	Method   string    `gorm:"type:varchar(8);not null" json:"method"`
	Success  int       `gorm:"type:smallint;default:200;not null" json:"success"`
	Timeout  int       `gorm:"type:smallint;default:60;not null" json:"timeout" validate:"gte=1,lte=300"`
	Retry    int       `gorm:"type:smallint;default:3;not null" json:"retry" validate:"gte=1,lte=10"`
	// Enabled  bool      `gorm:"type:bool;default:false" json:"enabled"`
	// Secret    string         `gorm:"type:text;default:null" json:"secret,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,->" json:"-"`
}

type TriggerCollection []*Trigger

func (t *Trigger) ToHAL(selfHref string) (root hal.Resource) {
	root = hal.NewResourceObject()
	root.AddData(t)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: selfHref}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return
}

func (c TriggerCollection) ToHAL(selfHref string, queryString url.Values) (root hal.Resource) {
	type NameOnly struct {
		Name string `json:"name"`
	}

	type Result struct {
		Count   int               `json:"count"`
		Results TriggerCollection `json:"results"`
	}

	root = hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	selfRel.SetLink(&hal.LinkObject{Href: selfHref})
	root.AddLink(selfRel)

	el, hasLast := Last(c)
	if hasLast {
		after, err := el.CreatedAt.MarshalText()
		if NoError(err) {
			queryString.Set(After, string(after))

			nextRel, _ := hal.NewLinkRelation(NextRelation)
			nextLink := &hal.LinkObject{Href: strings.Join([]string{selfHref, queryString.Encode()}, "?")}
			nextRel.SetLink(nextLink)
			root.AddLink(nextRel)
		}
	}

	var embedded []hal.Resource

	for _, i := range c {
		selfLink, _ := hal.NewLinkObject(fmt.Sprintf("%s/%v", selfHref, i.ID))

		selfRel, _ := hal.NewLinkRelation("self")
		selfRel.SetLink(selfLink)

		resource := hal.NewResourceObject()
		resource.AddLink(selfRel)
		resource.AddData(NameOnly{i.Name})

		embedded = append(embedded, resource)
	}

	triggers, _ := hal.NewResourceRelation("triggers")
	triggers.SetResources(embedded)
	root.AddResource(triggers)
	root.AddData(Result{len(c), c})

	return
}
