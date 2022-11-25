package utils

import (
	"crypto/rand"
	"encoding/hex"
	"enterprise.sidooh/pkg/cache"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	buf := argon2.IDKey([]byte(password), salt, 2, 32*1024, 1, 64)

	return fmt.Sprintf("%v.%v", hex.EncodeToString(buf), hex.EncodeToString(salt)), nil
}

func VerifyPassword(storedPassword string, suppliedPassword string) bool {
	split := strings.Split(storedPassword, ".")
	if len(split) < 2 {
		return false
	}

	salt, _ := hex.DecodeString(split[1])
	buf := argon2.IDKey([]byte(suppliedPassword), salt, 2, 32*1024, 1, 64)

	return hex.EncodeToString(buf) == split[0]
}

func CheckOTP(key string, otp int) bool {
	savedOtp := cache.Cache.Get(fmt.Sprintf("otp_%s", key))
	return savedOtp != nil && (*savedOtp).(int) == otp
}

func SetOTP(key string, otp int) {
	cache.Cache.Set(fmt.Sprintf("otp_%s", key), otp, 5*time.Minute)
}
