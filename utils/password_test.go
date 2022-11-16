package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashing(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isValidPassword := VerifyPassword(hashedPassword, password)

	require.True(t, isValidPassword)
}
