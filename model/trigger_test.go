package model

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func TestWorkspaceJsonMarshal(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	name := randstr.String(16)

	trigger1 := struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Schedule  string    `json:"schedule"`
		Timezone  string    `json:"timezone"`
		Url       string    `json:"url"`
		Method    string    `json:"method"`
		Success   int       `json:"success"`
		Timeout   int       `json:"timeout"`
		Retry     int       `json:"retry"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	trigger2 := Trigger{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	json1, _ := json.Marshal(trigger1)
	json2, _ := json.Marshal(trigger2)
	json1, _ = JSONRemarshal(json1)
	json2, _ = JSONRemarshal(json2)
	assert.Equal(t, string(json1), string(json2))
}

func TestSingleWorkspaceHAL(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	name := randstr.String(16)
	url := "/" + randstr.String(16)

	trigger := Trigger{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type Self struct {
		Href string `json:"href"`
	}

	type Links struct {
		Self Self `json:"self"`
	}

	type HAL struct {
		Links     Links     `json:"_links"`
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Method    string    `json:"method"`
		Retry     int       `json:"retry"`
		Schedule  string    `json:"schedule"`
		Success   int       `json:"success"`
		Timeout   int       `json:"timeout"`
		Timezone  string    `json:"timezone"`
		UpdatedAt time.Time `json:"updated_at"`
		Url       string    `json:"url"`
	}

	expected, _ := json.Marshal(HAL{Links: Links{
		Self: Self{
			Href: url,
		},
	},
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})

	resource := trigger.ToHAL(url)
	namedMap := resource.ToMap()
	actual, _ := json.Marshal(namedMap.Content)

	expected, _ = JSONRemarshal(expected)
	actual, _ = JSONRemarshal(actual)
	assert.Equal(t, string(expected), string(actual))
}

func TestMultipleWorkspacesHAL(t *testing.T) {
	now := time.Now()
	uid := uuid.New()
	name := randstr.String(16)
	path := "/" + randstr.String(16)
	queryString := url.Values{}

	e := Trigger{
		ID:        uid,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ec := TriggerCollection{}
	ec = append(ec, &e)

	type Href struct {
		Href string `json:"href"`
	}

	type Link struct {
		Self Href `json:"self"`
	}

	type LinkWithNext struct {
		Next Href `json:"next"`
		Self Href `json:"self"`
	}

	type Links struct {
		Link Link   `json:"_links"`
		Name string `json:"name,omitempty"`
	}

	type Embedded struct {
		Workspaces []Links `json:"triggers"`
	}

	type HAL struct {
		Embedded Embedded          `json:"_embedded"`
		Link     LinkWithNext      `json:"_links"`
		Count    int               `json:"count"`
		Results  TriggerCollection `json:"results"`
	}

	after, _ := e.CreatedAt.MarshalText()
	query := url.Values{}
	query.Add(After, string(after))

	hal, _ := json.Marshal(HAL{
		Embedded: Embedded{
			Workspaces: []Links{
				{
					Link{
						Self: Href{
							path + "/" + uid.String(),
						},
					},
					name,
				},
			}},
		Link: LinkWithNext{
			Next: Href{strings.Join([]string{path, query.Encode()}, "?")},
			Self: Href{path},
		},
		Count:   1,
		Results: ec,
	})

	resource := ec.ToHAL(path, queryString)
	namedMap := resource.ToMap()
	actual, _ := json.Marshal(namedMap.Content)
	expected, _ := JSONRemarshal(hal) // Sort keys to match with HAL's marshaling
	assert.Equal(t, string(expected), string(actual))
}
