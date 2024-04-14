package pg

import (
	"context"
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

type Project struct {
	Id        int       `gorm:"primaryKey;unique;autoIncrement"`
	UserUuid  uuid.UUID `gorm:"index:idx_user_id_name,unique"`
	Name      string    `gorm:"index:idx_user_id_name,unique"`
	Token     string    `gorm:"<-"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (p *Project) ToModel() *models.Project {
	return &models.Project{
		Id:        p.Id,
		UserUuid:  p.UserUuid,
		Name:      p.Name,
		Token:     p.Token,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ModelToProjectPg(project *models.Project) *Project {
	return &Project{
		Id:        project.Id,
		UserUuid:  project.UserUuid,
		Name:      project.Name,
		Token:     project.Token,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}

func (r *ProjectRepo) GetById(ctx context.Context, Id int) (*models.Project, error) {
	var project Project
	result := r.db.First(&project, Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByToken(ctx context.Context, token string) (*models.Project, error) {
	var project Project
	result := r.db.First(&project, "token = ?", token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByUserUuid(ctx context.Context, userUuid uuid.UUID) ([]models.Project, error) {
	var pgProjects []Project
	result := r.db.Find(&pgProjects, "user_uuid = ?", userUuid.String())
	if result.Error != nil {
		return nil, result.Error
	}
	projects := make([]models.Project, 0, len(pgProjects))
	for _, project := range pgProjects {
		projects = append(projects, *project.ToModel())
	}
	return projects, nil
}

func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	result := r.db.Create(ModelToProjectPg(project))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepo) Update(ctx context.Context, project *models.Project) error {
	result := r.db.Save(ModelToProjectPg(project))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepo) Delete(ctx context.Context, project *models.Project) error {
	result := r.db.Delete(ModelToProjectPg(project))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
