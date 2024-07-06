package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
)

func ErrorsHandler(c *fiber.Ctx) error {
	err := c.Next()

	if err == nil {
		return nil
	}

	var serr *errors.SError

	switch e := err.(type) {
	case *errors.SError:
		serr = e
	case *fiber.Error:
		return err
	default:
		serr = errors.Undefined.WithInfo(err.Error())
	}

	if serr == nil {
		return nil
	}

	if serr.IsBadRequest {
		c.Status(fiber.StatusBadRequest)
		c.JSON(serr.ToResponse())

		log.Debugw("Bad request", "code", serr.Code, "info", serr.Info, "field", serr.Field, "method", c.Method(), "uri", c.OriginalURL())
	} else {

		c.Status(fiber.StatusInternalServerError)

		isShowInfo := c.Locals(showInternalServerErrorParam)
		if isShowInfo != nil {
			c.JSON(serr.ToResponse())
		}

		log.Errorw("Internal server error", "code", serr.Code, "info", serr.Info, "field", serr.Field, "method", c.Method(), "uri", c.OriginalURL())

		return nil
	}

	return nil
}
