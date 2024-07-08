package HttpServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/project"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/token"
	"github.com/aerosystems/project-service/internal/presenters/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log                    *logrus.Logger
	echo                   *echo.Echo
	firebaseAuthMiddleware *middleware.FirebaseAuth
	projectHandler         *project.Handler
	tokenHandler           *token.Handler
}

func NewServer(
	log *logrus.Logger,
	errorHandler *echo.HTTPErrorHandler,
	firebaseAuthMiddleware *middleware.FirebaseAuth,
	projectHandler *project.Handler,
	tokenHandler *token.Handler,
) *Server {
	server := &Server{
		log:                    log,
		echo:                   echo.New(),
		firebaseAuthMiddleware: firebaseAuthMiddleware,
		projectHandler:         projectHandler,
		tokenHandler:           tokenHandler,
	}
	if errorHandler != nil {
		server.echo.HTTPErrorHandler = *errorHandler
	}
	return server
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server project-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}
