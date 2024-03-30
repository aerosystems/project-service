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
	ctx    context.Context
}

func NewProjectRepo(client *firestore.Client, ctx context.Context) *ProjectRepo {
	return &ProjectRepo{
		client: client,
		ctx:    ctx,
	}
}

func (r *ProjectRepo) GetById(Id int) (*models.Project, error) {
	docRef := r.client.Collection("projects").Doc(string(Id))
	doc, err := docRef.Get(r.ctx)
	if err != nil {
		return nil, err
	}

	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepo) GetByToken(Token string) (*models.Project, error) {
	iter := r.client.Collection("projects").Where("Token", "==", Token).Limit(1).Documents(r.ctx)
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

func (r *ProjectRepo) GetByUserUuid(userUuid uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	iter := r.client.Collection("projects").Where("UserUuid", "==", userUuid.String()).Documents(r.ctx)
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

func (r *ProjectRepo) Create(project *models.Project) error {
	_, _, err := r.client.Collection("projects").Add(r.ctx, project)
	return err
}

func (r *ProjectRepo) Update(project *models.Project) error {
	_, err := r.client.Collection("projects").Doc(string(project.Id)).Set(r.ctx, project)
	return err
}

func (r *ProjectRepo) Delete(project *models.Project) error {
	_, err := r.client.Collection("projects").Doc(string(project.Id)).Delete(r.ctx)
	return err
}
