package project

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InitProjectRequest struct {
	InitProjectRequestBody
}

type InitProjectRequestBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type CreateProjectEvent struct {
	CustomerUuid string `json:"customerUuid"`
}

// InitProject godoc
// @Summary Init project
// @Description Init project
// @Tags project
// @Accept json
// @Produce json
// @Param customerUuid body string true "Customer UUID"
// @Success 201 {object} Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project/init [post]
func (ph Handler) InitProject(c echo.Context) error {
	var req InitProjectRequest
	if err := c.Bind(&req); err != nil {
		return ph.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	var event CreateProjectEvent
	if err := json.Unmarshal(req.Message.Data, &event); err != nil {
		return ph.ErrorResponse(c, http.StatusUnprocessableEntity, "invalid request body", err)
	}
	project, err := ph.projectUsecase.InitProject(event.CustomerUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not init project", err)
	}

	return ph.SuccessResponse(c, http.StatusCreated, "project successfully created", project)
}
