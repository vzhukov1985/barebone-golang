package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/env"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"time"
)

type Server struct {
	address string
	app     *fiber.App
}

func CreateServer(c *fiber.App) *Server {
	address := env.GetString("REST_ADDRESS", "", true)
	return &Server{
		address: address,
		app:     c,
	}
}

func (s *Server) Start() {
	log.Infow("Starting REST server...", "address", s.address)
	s.app.Get("/healthCheck", func(ctx *fiber.Ctx) error {
		_, err := ctx.WriteString("OK")
		if err != nil {
			return err
		}
		return nil
	})

	go func() {
		if err := s.app.Listen(s.address); err != nil {
			log.Fatalw("Failed to start REST server", "error", err, "address", s.address)
		}
	}()
}

func (s *Server) Stop() {
	log.Info("Shutting down REST Server...")
	if err := s.app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalw("Failed to shutdown REST server", "error", err, "address", s.address)
	}
	log.Info("REST Server successfully was shut down")
}
