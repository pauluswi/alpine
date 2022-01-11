package pin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPinCrypto(t *testing.T) {

	// Generate Salt Unit Test
	salt := GenerateRandomSalt(saltSize)
	assert.Equal(t, saltSize, len(salt))

	// Hash Pin Unit Test
	hashedPassword := HashPassword("123456", salt)
	assert.NotEmpty(t, hashedPassword)

	// Comparing Hash Unit Test
	valid := DoPasswordsMatch(hashedPassword, "123456", salt)
	assert.Equal(t, true, valid)
}
