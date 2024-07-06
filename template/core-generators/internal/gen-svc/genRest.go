package gen_svc

import (
	gen_rest "core-generators/internal/gen-rest"
	"core-generators/internal/utils"
	"fmt"
	"os"
	"path"
)

func genRest(serviceName, svcPath string) error {
	restPath := path.Join(svcPath, "internal", "api", "rest")
	if err := os.MkdirAll(path.Join(restPath, "handlers"), 0777); err != nil {
		return err
	}

	formatFiles := []string{
		path.Join(restPath, "handlers", "handlers.go"),
		path.Join(restPath, "controller.go"),
		path.Join(restPath, "middlewares.go"),
		path.Join(restPath, "endpoints.go"),
	}

	if err := gen_rest.GenSpareHandlersService(path.Join(restPath, "handlers")); err != nil {
		return err
	}

	if err := genRestController(serviceName, restPath); err != nil {
		return err
	}

	if err := genMiddlewares(serviceName, restPath); err != nil {
		return err
	}

	if err := genEndpoints(serviceName, restPath); err != nil {
		return err
	}

	for _, f := range formatFiles {
		utils.FormatFile(f)
	}

	return nil
}

func genRestController(serviceName, restPath string) error {
	outData := fmt.Sprintf(`package rest

import (
	"core-%s/internal/api/rest/handlers"
	"github.com/gofiber/fiber/v2"
)

func CreateController() *fiber.App {
	c := fiber.New()

	h := handlers.Create()

	configureMiddleWares(c)
	configureEndPoints(c, h)

	return c
}`, serviceName)

	filePath := path.Join(restPath, "controller.go")
	if err := os.WriteFile(filePath, []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}

func genMiddlewares(serviceName, restPath string) error {
	outData := `package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/rest/middleware"
)

func configureMiddleWares(c *fiber.App) {
	c.Use(cors.New())
	c.Use(middleware.LogRequests)
	c.Use(middleware.ShowInternalServerErrorInfo)
	c.Use(middleware.ErrorsHandler)
}
`
	filePath := path.Join(restPath, "middlewares.go")
	if err := os.WriteFile(filePath, []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}

func genEndpoints(serviceName, restPath string) error {
	outData := fmt.Sprintf(`package rest

import (
	"core-%s/internal/api/rest/handlers"
	"github.com/gofiber/fiber/v2"
)

func configureEndPoints(c *fiber.App, h *handlers.Handlers) {
}
`, serviceName)
	filePath := path.Join(restPath, "endpoints.go")
	if err := os.WriteFile(filePath, []byte(outData), 0666); err != nil {
		return err
	}

	return nil
}
