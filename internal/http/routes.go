package HTTPServer

import (
	"github.com/aerosystems/project-service/internal/models"
)

func (s *Server) setupRoutes() {
	// Private routes X-Api-Key
	s.echo.GET("/v1/token/validate", app.baseHandler.ValidateToken)

	// Private routes OAuth 2.0: check roles [customer, staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/projects", app.baseHandler.GetProjectList, app.oauthMiddleware.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))
	s.echo.GET("/v1/projects/:projectId", app.baseHandler.GetProject, app.oauthMiddleware.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))

	// Private routes OAuth 2.0: check roles [staff]. Auth implemented on API Gateway
	s.echo.POST("/v1/projects", app.baseHandler.ProjectCreate, app.oauthMiddleware.AuthTokenMiddleware(models.StaffRole))
	s.echo.PATCH("/v1/projects/:projectId", app.baseHandler.ProjectUpdate, app.oauthMiddleware.AuthTokenMiddleware(models.StaffRole))
	s.echo.DELETE("/v1/projects/:projectId", app.baseHandler.ProjectDelete, app.oauthMiddleware.AuthTokenMiddleware(models.StaffRole))
}
