package usecases

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type ProjectUsecase struct {
	projectRepo ProjectRepository
	subsRepo    SubsRepository
	strategy    Strategy
}

func NewProjectUsecase(projectRepo ProjectRepository, subsRPC SubsRepository) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo: projectRepo,
		subsRepo:    subsRPC,
	}
}

func (ps *ProjectUsecase) DetermineStrategy(userUuidStr string, role string) error {
	userUuid, err := uuid.Parse(userUuidStr)
	if err != nil {
		return err
	}
	if role == models.StaffRole.String() {
		ps.SetStrategy(&StaffStrategy{userUuid})
		return nil
	}
	kind, _, err := ps.subsRepo.GetSubscription(userUuid)
	if err != nil {
		return err
	}
	switch kind {
	case models.TrialSubscription:
		fallthrough
	case models.StartupSubscription:
		ps.SetStrategy(&StartupStrategy{userUuid})
	case models.BusinessSubscription:
		ps.SetStrategy(&BusinessStrategy{userUuid})
	default:
		return errors.New("unknown subscription kind")
	}
	return nil
}

func (ps *ProjectUsecase) GetProjectById(projectId int) (*models.Project, error) {
	ctx := context.Background()
	project, err := ps.projectRepo.GetById(ctx, projectId)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectByToken(token string) (*models.Project, error) {
	ctx := context.Background()
	project, err := ps.projectRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	// TODO: if it statement is needed, we should determine strategy before
	//if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
	//	return nil, errors.New("user is not allowed to access the project")
	//}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error) {
	if filterUserUuid != uuid.Nil {
		if !ps.strategy.IsAccessibleByUserUuid(filterUserUuid) {
			return []models.Project{}, nil
		}
		ctx := context.Background()
		projectList, err = ps.projectRepo.GetByUserUuid(ctx, filterUserUuid)
	}
	ctx := context.Background()
	projectList, err = ps.projectRepo.GetByUserUuid(ctx, userUuid)
	return projectList, nil
}

func (ps *ProjectUsecase) CreateDefaultProject(userUuid uuid.UUID) error {
	if err := ps.CreateProject(userUuid, "default"); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) CreateProject(userUuid uuid.UUID, name string) error {
	ctx := context.Background()
	projectList, err := ps.projectRepo.GetByUserUuid(ctx, userUuid)
	if err != nil {
		return err
	}
	if ps.isProjectNameExist(name, projectList) {
		return errors.New("project name already exists")
	}
	if !ps.strategy.IsAccessibleByUserUuid(userUuid) {
		return errors.New("user is not allowed to create a project")
	}
	if !ps.strategy.IsAccessibleByCountProjects(len(projectList)) {
		return errors.New("out of projects limit")
	}
	newProject := NewProject(userUuid, name)
	ctx = context.Background()
	if err := ps.projectRepo.Create(ctx, newProject); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) UpdateProject(project *models.Project) error {
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return errors.New("user is not allowed to update the project")
	}
	ctx := context.Background()
	if err := ps.projectRepo.Update(ctx, project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) DeleteProjectById(projectId int) error {
	ctx := context.Background()
	project, err := ps.projectRepo.GetById(ctx, projectId)
	if err != nil {
		return err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return errors.New("user is not allowed to delete the project")
	}
	ctx = context.Background()
	if err := ps.projectRepo.Delete(ctx, project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) IsProjectExistByToken(projectToken string) bool {
	ctx := context.Background()
	project, err := ps.projectRepo.GetByToken(ctx, projectToken)
	if err != nil {
		return false
	}
	if project == nil {
		return false
	}
	return true
}

func (ps *ProjectUsecase) SetStrategy(strategy Strategy) {
	ps.strategy = strategy
}

func (ps *ProjectUsecase) isProjectNameExist(name string, projectList []models.Project) bool {
	for _, project := range projectList {
		if project.Name == name {
			return true
		}
	}
	return false
}

func NewProject(userUuid uuid.UUID, name string) *models.Project {
	return &models.Project{
		Token:    generateToken(),
		UserUuid: userUuid,
		Name:     name,
	}
}

func generateToken() string {
	rand.Seed(time.Now().Unix())
	sum := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	return fmt.Sprintf("%x", sum)
}
