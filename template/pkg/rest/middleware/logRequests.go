package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"time"
)

func LogRequests(c *fiber.Ctx) error {
	start := time.Now()
	log.Debugw("Incoming request", "method", c.Method(), "uri", c.OriginalURL())

	c.Next()
	log.Debugw("Request time execution", "time ms", time.Since(start).Milliseconds(), "method", c.Method(), "uri", c.OriginalURL())
	return nil
}
