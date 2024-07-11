package HttpServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/project"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/token"
	"github.com/aerosystems/project-service/internal/presenters/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	port                   int
	log                    *logrus.Logger
	echo                   *echo.Echo
	firebaseAuthMiddleware *middleware.FirebaseAuth
	projectHandler         *project.Handler
	tokenHandler           *token.Handler
}

func NewServer(
	port int,
	log *logrus.Logger,
	errorHandler *echo.HTTPErrorHandler,
	firebaseAuthMiddleware *middleware.FirebaseAuth,
	projectHandler *project.Handler,
	tokenHandler *token.Handler,
) *Server {
	server := &Server{
		port:                   port,
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
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
