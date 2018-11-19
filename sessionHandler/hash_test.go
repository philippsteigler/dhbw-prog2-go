package sessionHandler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt := generateSalt()
	assert.Len(t, salt, 16, "Length should be 16 Byte.")

}

func TestHashString(t *testing.T) {
	text := "test123"

	hash, salt := HashString(text)
	assert.NotNil(t, hash, salt)
}

func TestGetHash(t *testing.T) {
	text := "test123"
	salt := "8eE3dChRYBJtIEfE"

	hash := GetHash(text, salt)
	assert.Equal(t, "L_PF-nkmRI_TXi1g8kbiA-XmqjRA1_S99CylqNmtp5g=", hash, "Hashes should be equal.")
}
