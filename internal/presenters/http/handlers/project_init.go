package handlers

import (
	"encoding/json"
	"fmt"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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
	CustomerUuid     string    `json:"customerUuid"`
	SubscriptionType string    `json:"subscriptionType"`
	AccessTime       time.Time `json:"accessTime"`
}

// InitProject godoc
// @Summary Init project
// @Description Init project
// @Tags project
// @Accept json
// @Produce json
// @Param customerUuid body string true "Customer UUID"
// @Success 200 {object} Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project/init [post]
func (ph ProjectHandler) InitProject(c echo.Context) error {
	var req InitProjectRequest
	if err := c.Bind(&req); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	var event CreateProjectEvent
	if err := json.Unmarshal(req.Message.Data, &event); err != nil {
		return CustomErrors.ErrInvalidRequestPayload
	}
	project, err := ph.projectUsecase.InitProject(event.CustomerUuid, event.SubscriptionType, event.AccessTime)
	if err != nil {
		fmt.Printf("Error: %v\n. Type: %T\n", err, err)
		return err
	}

	return c.JSON(http.StatusCreated, ModelToProject(project))
}
