package HTTPServer

import (
	"context"
	"github.com/aerosystems/common-service/presenters/httpserver"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	srv *httpserver.Server
}

func NewHTTPServer(
	cfg *Config,
	log *logrus.Logger,
	firebaseAuth *FirebaseAuth,
	handler *Handler,
) *Server {
	return &Server{
		srv: httpserver.NewHTTPServer(
			&httpserver.Config{
				Host: cfg.Host,
				Port: cfg.Port,
			},

			httpserver.WithCustomErrorHandler(httpserver.NewCustomErrorHandler(cfg.Mode)),

			httpserver.WithValidator(httpserver.NewCustomValidator()),

			httpserver.WithMiddleware(middleware.CORSWithConfig(middleware.CORSConfig{
				Skipper:      middleware.DefaultSkipper,
				AllowOrigins: []string{"*"},
				AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
			})),
			httpserver.WithMiddleware(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
				LogURI:    true,
				LogStatus: true,
				LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
					log.WithFields(logrus.Fields{
						"URI":    values.URI,
						"status": values.Status,
					}).Info("request")

					return nil
				},
			})),
			httpserver.WithMiddleware(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					c.Request().WithContext(logctx.New(c.Request().Context(), logrus.NewEntry(log)))
					return next(c)
				}
			}),
			httpserver.WithMiddleware(middleware.Recover()),

			httpserver.WithRouter(http.MethodPost, "/v1/token/validate", handler.ValidateToken),
			httpserver.WithRouter(http.MethodGet, "/v1/projects", handler.GetProjectList,
				firebaseAuth.RoleBasedAuth(entities.CustomerRole, entities.StaffRole)),
			httpserver.WithRouter(http.MethodPost, "/v1/projects/:project_uuid", handler.GetProject,
				firebaseAuth.RoleBasedAuth(entities.CustomerRole, entities.StaffRole)),
			httpserver.WithRouter(http.MethodPost, "/v1/projects", handler.ProjectCreate,
				firebaseAuth.RoleBasedAuth(entities.CustomerRole, entities.StaffRole)),
			httpserver.WithRouter(http.MethodPatch, "/v1/projects/:project_uuid", handler.UpdateProject,
				firebaseAuth.RoleBasedAuth(entities.StaffRole)),
			httpserver.WithRouter(http.MethodDelete, "/v1/projects/:project_uuid", handler.ProjectDelete,
				firebaseAuth.RoleBasedAuth(entities.StaffRole)),
		),
	}
}

func (s *Server) Run() error {
	return s.srv.Run()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
