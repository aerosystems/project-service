package adapters

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"time"
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
	Uuid         string    `firestore:"uuid"`
	CustomerUuid string    `firestore:"customer_uuid"`
	Name         string    `firestore:"name"`
	Token        string    `firestore:"token"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

func (p *Project) ToModel() *models.Project {
	return &models.Project{
		Uuid:         uuid.MustParse(p.Uuid),
		CustomerUUID: uuid.MustParse(p.CustomerUuid),
		Name:         p.Name,
		Token:        p.Token,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func ModelToProject(project *models.Project) *Project {
	return &Project{
		Uuid:         project.Uuid.String(),
		CustomerUuid: project.CustomerUUID.String(),
		Name:         project.Name,
		Token:        project.Token,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
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

func (r *ProjectRepo) GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Project, error) {
	docRef := r.client.Collection("projects").Doc(uuid.String())
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !doc.Exists() {
		return nil, CustomErrors.ErrProjectNotFound
	}
	var project Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}
	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByToken(ctx context.Context, token string) (*models.Project, error) {
	iter := r.client.Collection("projects").Where("token", "==", token).Limit(1).Documents(ctx)
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

func (r *ProjectRepo) GetByCustomerUuidAndName(ctx context.Context, customerUuid uuid.UUID, name string) (*models.Project, error) {
	iter := r.client.Collection("projects").Where("customer_uuid", "==", customerUuid.String()).Where("name", "==", name).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		return nil, err
	}

	var project Project
	if err = doc.DataTo(&project); err != nil {
		return nil, err
	}

	return project.ToModel(), nil
}

func (r *ProjectRepo) GetByCustomerUuid(ctx context.Context, customerUuid uuid.UUID) ([]models.Project, error) {
	var fireProjects []Project
	iter := r.client.Collection("projects").Where("customer_uuid", "==", customerUuid.String()).Documents(ctx)
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
	_, err := r.client.Collection("projects").Doc(project.Uuid.String()).Set(ctx, ModelToProject(project))
	return err
}

func (r *ProjectRepo) Update(ctx context.Context, project *models.Project) error {
	fireProject := ModelToProject(project)
	_, err := r.client.Collection("projects").Doc(project.Uuid.String()).Set(ctx, fireProject)
	return err
}

func (r *ProjectRepo) Delete(ctx context.Context, project *models.Project) error {
	_, err := r.client.Collection("projects").Doc(project.Uuid.String()).Delete(ctx)
	return err
}
