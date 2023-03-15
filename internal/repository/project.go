package repository

import (
	"github.com/aerosystems/project-service/internal/helpers"
	"github.com/aerosystems/project-service/internal/models"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

func (r *ProjectRepo) FindByID(ID int) (*models.Project, error) {
	var project models.Project
	result := r.db.Find(&project, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepo) FindByToken(token string) (*models.Project, error) {
	var project models.Project
	result := r.db.Find(&project, token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepo) FindByUserID(UserID int) (*models.Project, error) {
	var project models.Project
	result := r.db.Find(&project, UserID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepo) Create(project *models.Project) error {
	project.Token = helpers.GenerateToken()
	result := r.db.Create(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepo) Update(project *models.Project) error {
	result := r.db.Save(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepo) Delete(project *models.Project) error {
	result := r.db.Delete(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
