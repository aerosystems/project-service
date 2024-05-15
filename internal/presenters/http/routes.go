package HttpServer

import (
	"github.com/aerosystems/project-service/internal/models"
)

func (s *Server) setupRoutes() {
	// Private routes X-Api-Key
	s.echo.GET("/v1/token/validate", s.tokenHandler.ValidateToken)

	// Private routes OAuth 2.0: check roles [customer, staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/projects", s.projectHandler.GetProjectList, s.firebaseAuthMiddleware.RoleBased(models.CustomerRole, models.StaffRole))
	s.echo.GET("/v1/projects/:projectId", s.projectHandler.GetProject, s.firebaseAuthMiddleware.RoleBased(models.CustomerRole, models.StaffRole))

	// Private routes OAuth 2.0: check roles [staff]. Auth implemented on API Gateway
	s.echo.POST("/v1/projects", s.projectHandler.ProjectCreate, s.firebaseAuthMiddleware.RoleBased(models.CustomerRole, models.StaffRole))
	s.echo.PATCH("/v1/projects/:projectId", s.projectHandler.UpdateProject, s.firebaseAuthMiddleware.RoleBased(models.StaffRole))
	s.echo.DELETE("/v1/projects/:projectId", s.projectHandler.ProjectDelete, s.firebaseAuthMiddleware.RoleBased(models.StaffRole))
}
