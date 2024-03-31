package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type ProjectRepo struct {
	client *firestore.Client
}

func NewProjectRepo(client *firestore.Client) *ProjectRepo {
	return &ProjectRepo{
		client: client,
	}
}

func (r *ProjectRepo) GetById(ctx context.Context, id int) (*models.Project, error) {
	docRef := r.client.Collection("projects").Doc(string(id))
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return &project, nil
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

	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepo) GetByUserUuid(ctx context.Context, userUuid uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	iter := r.client.Collection("projects").Where("UserUuid", "==", userUuid.String()).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	_, _, err := r.client.Collection("projects").Add(ctx, project)
	return err
}

func (r *ProjectRepo) Update(ctx context.Context, project *models.Project) error {
	_, err := r.client.Collection("projects").Doc(string(project.Id)).Set(ctx, project)
	return err
}

func (r *ProjectRepo) Delete(ctx context.Context, project *models.Project) error {
	_, err := r.client.Collection("projects").Doc(string(project.Id)).Delete(ctx)
	return err
}
