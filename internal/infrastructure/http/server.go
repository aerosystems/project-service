package HttpServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/infrastructure/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log            *logrus.Logger
	echo           *echo.Echo
	accessSecret   string
	projectHandler *handlers.ProjectHandler
	tokenHandler   *handlers.TokenHandler
}

func NewServer(
	log *logrus.Logger,
	accessSecret string,
	projectHandler *handlers.ProjectHandler,
	tokenHandler *handlers.TokenHandler,
) *Server {
	return &Server{
		log:            log,
		echo:           echo.New(),
		accessSecret:   accessSecret,
		projectHandler: projectHandler,
		tokenHandler:   tokenHandler,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server project-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}