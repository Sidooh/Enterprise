package utils

import (
	"enterprise.sidooh/pkg/cache"
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

func TestOTP(t *testing.T) {
	cache.Init()

	otp := int(RandomInt(3))
	SetOTP("OTP", otp)

	validOTP := CheckOTP("OTP", otp)
	require.True(t, validOTP)

	validOTP = CheckOTP("NO_OTP", otp)
	require.False(t, validOTP)

	validOTP = CheckOTP("RANDOM", otp)
	require.False(t, validOTP)
}
