package handlers

type TokenHandler struct {
	*BaseHandler
	tokenUsecase TokenUsecase
}

func NewTokenHandler(baseHandler *BaseHandler, tokenUsecase TokenUsecase) *TokenHandler {
	return &TokenHandler{
		BaseHandler:  baseHandler,
		tokenUsecase: tokenUsecase,
	}
}
