package sessionHandler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashString(t *testing.T) {
	text := "test123"

	hash, salt := HashString(text)
	assert.NotNil(t, hash, salt)
}

func TestGetHash(t *testing.T) {
	text := "test123"
	salt := "8eE3dChRYBJtIEfE"

	hash := GetHash(text, salt)
	assert.Equal(t, hash, "L_PF-nkmRI_TXi1g8kbiA-XmqjRA1_S99CylqNmtp5g=", "Hashes should be equal.")
}
