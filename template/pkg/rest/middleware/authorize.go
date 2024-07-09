package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/models"
)

func Authorize(roles ...models.UserRole) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
