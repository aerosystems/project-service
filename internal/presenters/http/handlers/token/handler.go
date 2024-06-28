package token

import (
	"github.com/aerosystems/project-service/internal/presenters/http/handlers"
)

type Handler struct {
	*handlers.BaseHandler
	tokenUsecase handlers.TokenUsecase
}

func NewTokenHandler(baseHandler *handlers.BaseHandler, tokenUsecase handlers.TokenUsecase) *Handler {
	return &Handler{
		BaseHandler:  baseHandler,
		tokenUsecase: tokenUsecase,
	}
}
