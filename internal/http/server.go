package HttpServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/infrastructure/rest"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log            *logrus.Logger
	echo           *echo.Echo
	projectHandler *rest.ProjectHandler
	tokenHandler   *rest.TokenHandler
	tokenService   TokenService
}

func NewServer(
	log *logrus.Logger,
	projectHandler *rest.ProjectHandler,
	tokenHandler *rest.TokenHandler,
	tokenService TokenService,
) *Server {
	return &Server{
		log:            log,
		echo:           echo.New(),
		projectHandler: projectHandler,
		tokenHandler:   tokenHandler,
		tokenService:   tokenService,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server project-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}
