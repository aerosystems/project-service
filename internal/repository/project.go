package repository

import (
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
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

func (r *ProjectRepo) GetById(Id int) (*models.Project, error) {
	var project models.Project
	result := r.db.First(&project, Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepo) GetByToken(Token string) (*models.Project, error) {
	var project models.Project
	result := r.db.First(&project, "token = ?", Token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepo) GetByUserUuid(userUuid uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Find(&projects, "user_uuid = ?", userUuid.String())
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *ProjectRepo) Create(project *models.Project) error {
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
