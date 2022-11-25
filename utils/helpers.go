package utils

import (
	"encoding/json"
	"enterprise.sidooh/pkg/entities"
	"github.com/Permify/permify-gorm/collections"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"strings"
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

func GetUser(ctx *fiber.Ctx) *entities.UserWithEnterprise {
	return ctx.Locals("user").(*entities.UserWithEnterprise)
}

func GetEnterprise(ctx *fiber.Ctx) entities.Enterprise {
	return GetUser(ctx).Enterprise
}

func GetEnterpriseId(ctx *fiber.Ctx) int {
	return int(GetUser(ctx).EnterpriseId)
}
