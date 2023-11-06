package services

import "github.com/aerosystems/project-service/internal/models"

type ProjectService interface {
	CreateDefaultProject(userId int) error
}

type ProjectServiceImpl struct {
	projectRepo models.ProjectRepository
}

func NewProjectServiceImpl(projectRepo models.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
	}
}

func (ps *ProjectServiceImpl) CreateDefaultProject(userId int) error {
	var newProject = models.Project{
		UserId: userId,
		Name:   "default",
	}

	if err := ps.projectRepo.Create(&newProject); err != nil {
		return err
	}

	return nil
}
