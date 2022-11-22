package utils

import (
	"encoding/json"
	"enterprise.sidooh/pkg/cache"
	"enterprise.sidooh/pkg/entities"
	"fmt"
	"github.com/Permify/permify-gorm/collections"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

func ConvertStruct(from interface{}, to interface{}) {
	record, _ := json.Marshal(from)
	_ = json.Unmarshal(record, &to)
}

func HasRole(ctx *fiber.Ctx, role string) bool {
	return slices.Contains(ctx.Locals("roles").(collections.Role).GuardNames(), strings.ToLower(role))
}

func IsSuperAdmin(ctx *fiber.Ctx) bool {
	return HasRole(ctx, SUPERADMIN)
}

func IsAdmin(ctx *fiber.Ctx) bool {
	return HasRole(ctx, ADMIN)
}

func GetEnterpriseId(ctx *fiber.Ctx) int {
	return int(ctx.Locals("user").(*entities.UserWithEnterprise).EnterpriseId)
}

func CheckOTP(key string, otp int) bool {
	savedOtp := cache.Cache.Get(fmt.Sprintf("otp_%s", key))
	return savedOtp != nil && (*savedOtp).(int) == otp
}

func SetOTP(key string, otp int) {
	cache.Cache.Set(fmt.Sprintf("otp_%s", key), otp, 5*time.Minute)
}
