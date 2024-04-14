package usecases

import (
	"context"
)

type TokenUsecase struct {
	projectRepo ProjectRepository
}

func NewTokenUsecase(projectRepo ProjectRepository) *TokenUsecase {
	return &TokenUsecase{
		projectRepo: projectRepo,
	}
}

func (tu TokenUsecase) IsProjectExistByToken(projectToken string) bool {
	ctx := context.Background()
	project, err := tu.projectRepo.GetByToken(ctx, projectToken)
	if err != nil {
		return false
	}
	if project == nil {
		return false
	}
	return true
}
