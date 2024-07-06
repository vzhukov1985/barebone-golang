package middleware

import "github.com/gofiber/fiber/v2"

const showInternalServerErrorParam = "showInternalServerErrorInfo"

func ShowInternalServerErrorInfo(ctx *fiber.Ctx) error {
	ctx.Locals(showInternalServerErrorParam, true)
	return ctx.Next()
}
