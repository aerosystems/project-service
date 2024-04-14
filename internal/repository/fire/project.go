package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"time"
)

const (
	defaultTimeout = 2 * time.Second
)

type ProjectRepo struct {
	client *firestore.Client
}

func NewProjectRepo(client *firestore.Client) *ProjectRepo {
	return &ProjectRepo{
		client: client,
	}
}

type Project struct {
	Id        int       `firestore:"id"`
	UserUuid  string    `firestore:"user_uuid"`
	Name      string    `firestore:"name"`
	Token     string    `firestore:"token"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

func (p *Project) ToModel() *models.Project {
	return &models.Project{
		Id:        p.Id,
		UserUuid:  uuid.MustParse(p.UserUuid),
		Name:      p.Name,
		Token:     p.Token,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ModelToProject(project *models.Project) *Project {
	return &Project{
		Id:        project.Id,
		UserUuid:  project.UserUuid.String(),
		Name:      project.Name,
		Token:     project.Token,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}

func ModelListToProjectList(projects []models.Project) []Project {
	projectList := make([]Project, 0, len(projects))
	for _, project := range projects {
		projectList = append(projectList, *ModelToProject(&project))
	}
	return projectList
}

func ProjectListToModelList(projects []Project) []models.Project {
	projectList := make([]models.Project, 0, len(projects))
	for _, project := range projects {
		projectList = append(projectList, *project.ToModel())
	}
	return projectList
}

func (r *ProjectRepo) GetById(ctx context.Context, id int) (*models.Project, error) {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	docRef := r.client.Collection("projects").Doc(string(id))
	doc, err := docRef.Get(c)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByToken(ctx context.Context, token string) (*models.Project, error) {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	iter := r.client.Collection("projects").Where("token", "==", token).Limit(1).Documents(c)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		return nil, err
	}

	var project Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByUserUuid(ctx context.Context, userUuid uuid.UUID) ([]models.Project, error) {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	var fireProjects []Project
	iter := r.client.Collection("fireProjects").Where("UserUuid", "==", userUuid.String()).Documents(c)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var project Project
		if err := doc.DataTo(&project); err != nil {
			return nil, err
		}
		fireProjects = append(fireProjects, project)
	}
	return ProjectListToModelList(fireProjects), nil
}

func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	fireProject := ModelToProject(project)
	_, err := r.client.Collection("projects").Doc(string(fireProject.Id)).Create(c, fireProject)
	return err
}

func (r *ProjectRepo) Update(ctx context.Context, project *models.Project) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	fireProject := ModelToProject(project)
	_, err := r.client.Collection("projects").Doc(string(fireProject.Id)).Set(c, fireProject)
	return err
}

func (r *ProjectRepo) Delete(ctx context.Context, project *models.Project) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	_, err := r.client.Collection("projects").Doc(string(project.Id)).Delete(c)
	return err
}
