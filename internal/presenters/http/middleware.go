package HTTPServer

import (
	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) setupMiddleware() {
	s.addLog(s.log)
	s.addCORS()
	s.echo.Use(s.addLoggerToContext)
}

func (s *Server) addLog(log *logrus.Logger) {
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))
	s.echo.Use(middleware.Recover())
}

func (s *Server) addCORS() {
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}
	s.echo.Use(middleware.CORSWithConfig(DefaultCORSConfig))
}

func (s *Server) addLoggerToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().WithContext(logctx.New(c.Request().Context(), logrus.NewEntry(s.log)))
		return next(c)
	}
}
