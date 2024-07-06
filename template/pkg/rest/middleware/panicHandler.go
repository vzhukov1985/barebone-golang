package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"runtime/debug"
)

func PanicHandler(c *fiber.Ctx) error {
	defer func() {
		r := recover()
		if r != nil {
			var err error
			switch t := r.(type) {
			case string:
				err = errors.Undefined.WithInfo(t)
			case error:
				err = t
			default:
				err = errors.Undefined.WithInfo("unknown")
			}

			c.Status(fiber.StatusInternalServerError)
			log.Errorw("PANIC!!!!", "error", err, "method", c.Method(), "uri", c.OriginalURL(), "stack", string(debug.Stack()))
		}
	}()

	c.Next()
	return nil
}
