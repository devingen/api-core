package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInt(t *testing.T) {
	var field Field = map[string]interface{}{
		"limit": 1,
	}
	assert.Equal(t, field.GetInt("limit"), 1)
}

func TestGetIntFromFloat(t *testing.T) {
	var field Field = map[string]interface{}{
		"limit": float64(1),
	}
	assert.Equal(t, field.GetInt("limit"), 1)
}
