package ticket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewId(t *testing.T) {
	newId := NewId()
	assert.Equal(t, 1, newId)
	assert.Equal(t, 2, id.FreeId)
	Reset()
}

func TestReset(t *testing.T) {
	NewId()
	Reset()
	assert.Equal(t, 1, id.FreeId)
}
