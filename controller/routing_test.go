package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	server := InitServer()

	server.registerRoutes()

	assert.NotEmpty(t, server.router.Routes())
}
