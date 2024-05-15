package handlers

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProjectHandler struct {
	*BaseHandler
	projectUsecase ProjectUsecase
}

func NewProjectHandler(baseHandler *BaseHandler, projectUsecase ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler:    baseHandler,
		projectUsecase: projectUsecase,
	}
}

type CreateProjectRequest struct {
	UserUuid string `json:"userUuid" validate:"required,number" example:"66"`
	Name     string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body CreateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 201 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [post]
func (ph ProjectHandler) ProjectCreate(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	var requestPayload CreateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return ph.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	userUuid, err := uuid.Parse(requestPayload.UserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	if err := ph.projectUsecase.CreateProject(userUuid, requestPayload.Name); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not create default project", err)
	}
	return ph.SuccessResponse(c, http.StatusCreated, "project successfully created", nil)
}

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param	userUuid	query	string	false "UserUuid"
// @Success 200 {object} Response{data=[]models.Project}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (ph ProjectHandler) GetProjectList(c echo.Context) (err error) {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	uuidStr := c.QueryParam("userUuid")
	if len(uuidStr) == 0 {
		uuidStr = accessTokenClaims.UserUuid
	}
	filterUserUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "getting project is forbidden", err)
	}
	userUuid, err := uuid.Parse(accessTokenClaims.UserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	projectList, err := ph.projectUsecase.GetProjectListByUserUuid(userUuid, filterUserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not get projects", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "projects successfully found", projectList)
}

// GetProject godoc
// @Summary get project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectId} [get]
func (ph ProjectHandler) GetProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "getting project is forbidden", err)
	}
	project, err := ph.projectUsecase.GetProjectById(projectId)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully found", project)
}

// UpdateProject godoc
// @Summary update project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Param comment body UpdateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectId} [patch]
func (ph ProjectHandler) UpdateProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	var requestPayload UpdateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return ph.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	project, err := ph.projectUsecase.GetProjectById(projectId)
	if err != nil && project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	project.Name = requestPayload.Name
	if err := ph.projectUsecase.UpdateProject(project); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not update project", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully updated", project)
}

// ProjectDelete godoc
// @Summary delete project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectId} [delete]
func (ph ProjectHandler) ProjectDelete(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	if project, err := ph.projectUsecase.GetProjectById(projectId); err != nil && project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if err := ph.projectUsecase.DeleteProjectById(projectId); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not delete project", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully deleted", nil)
}
