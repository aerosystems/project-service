package HTTPServer

import (
	"github.com/aerosystems/project-service/internal/common/config"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type BaseHandler struct {
	mode      string
	log       *logrus.Logger
	cfg       *config.Config
	validator validator.Validate
}

func NewBaseHandler(
	log *logrus.Logger,
	mode string,
) *BaseHandler {
	return &BaseHandler{
		mode:      mode,
		log:       log,
		validator: validator.Validate{},
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
