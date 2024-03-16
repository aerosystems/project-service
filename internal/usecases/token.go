package usecases

type TokenUsecase struct {
	projectRepo ProjectRepository
}

func NewTokenUsecase(projectRepo ProjectRepository) *TokenUsecase {
	return &TokenUsecase{
		projectRepo: projectRepo,
	}
}

func (tu TokenUsecase) IsProjectExistByToken(projectToken string) bool {
	project, err := tu.projectRepo.GetByToken(projectToken)
	if err != nil {
		return false
	}
	if project == nil {
		return false
	}
	return true
}
