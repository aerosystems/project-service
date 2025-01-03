package usecases

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

const (
	defaultProjectName = "Default Project"
)

type ProjectUsecase struct {
	projectRepo            ProjectRepository
	subsRepo               SubsRepository
	checkmailEventsAdapter CheckmailEventsAdapter
	strategy               Strategy
}

func NewProjectUsecase(projectRepo ProjectRepository, subsRPC SubsRepository, checkmailEventsAdapter CheckmailEventsAdapter) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo:            projectRepo,
		subsRepo:               subsRPC,
		checkmailEventsAdapter: checkmailEventsAdapter,
	}
}

func (ps *ProjectUsecase) DetermineStrategy(customerUuidStr string, role string) error {
	customerUuid, err := uuid.Parse(customerUuidStr)
	if err != nil {
		return err
	}
	if role == models.StaffRole.String() {
		ps.SetStrategy(&StaffStrategy{customerUuid})
		return nil
	}
	kind, _, err := ps.subsRepo.GetSubscription(customerUuid)
	if err != nil {
		return err
	}
	switch kind {
	case models.TrialSubscription:
		fallthrough
	case models.StartupSubscription:
		ps.SetStrategy(&StartupStrategy{customerUuid})
	case models.BusinessSubscription:
		ps.SetStrategy(&BusinessStrategy{customerUuid})
	default:
		return errors.New("unknown subscription kind")
	}
	return nil
}

func (ps *ProjectUsecase) GetProjectByUuid(projectUuidStr string) (*models.Project, error) {
	projectUuid, err := uuid.Parse(projectUuidStr)
	if err != nil {
		return nil, CustomErrors.ErrProjectUuidInvalid
	}
	ctx := context.Background()
	project, err := ps.projectRepo.GetByUuid(ctx, projectUuid)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
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
	//if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
	//	return nil, errors.New("user is not allowed to access the project")
	//}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectListByCustomerUuid(customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error) {
	if filterUserUuid != uuid.Nil {
		if !ps.strategy.IsAccessibleByCustomerUuid(filterUserUuid) {
			return []models.Project{}, nil
		}
		ctx := context.Background()
		projectList, err = ps.projectRepo.GetByCustomerUuid(ctx, filterUserUuid)
	}
	ctx := context.Background()
	projectList, err = ps.projectRepo.GetByCustomerUuid(ctx, customerUuid)
	return projectList, nil
}

func (ps *ProjectUsecase) CreateDefaultProject(customerUuid uuid.UUID) (*models.Project, error) {
	return ps.CreateProject(customerUuid, "default")
}

func (ps *ProjectUsecase) InitProject(customerUuidStr string, subscriptionType string, accessTime time.Time) (*models.Project, error) {
	customerUuid, err := uuid.Parse(customerUuidStr)
	if err != nil {
		return nil, CustomErrors.ErrProjectUuidInvalid
	}

	ctx := context.Background()
	if defaultCustomerProject, err := ps.projectRepo.GetByCustomerUuidAndName(ctx, customerUuid, defaultProjectName); err == nil && defaultCustomerProject != nil {
		return nil, CustomErrors.ErrProjectAlreadyExists
	}

	newDefaultProject := NewProject(customerUuid, defaultProjectName)
	if err := ps.projectRepo.Create(ctx, newDefaultProject); err != nil {
		return nil, err
	}

	project, err := ps.projectRepo.GetByCustomerUuidAndName(ctx, customerUuid, defaultProjectName)
	if err != nil {
		return nil, err
	}
	if err := ps.checkmailEventsAdapter.PublishCreateAccessEvent(project.Token, subscriptionType, accessTime); err != nil {
		return nil, err
	}
	return project, nil
}

func (ps *ProjectUsecase) CreateProject(customerUuid uuid.UUID, name string) (*models.Project, error) {
	ctx := context.Background()
	projectList, err := ps.projectRepo.GetByCustomerUuid(ctx, customerUuid)
	if err != nil {
		return nil, err
	}
	if ps.isProjectNameExist(name, projectList) {
		return nil, CustomErrors.ErrProjectNameExists
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(customerUuid) {
		return nil, CustomErrors.ErrForbidden
	}
	if !ps.strategy.IsAccessibleByCountProjects(len(projectList)) {
		return nil, CustomErrors.ErrProjectLimitExceeded
	}
	project := NewProject(customerUuid, name)
	ctx = context.Background()
	if err := ps.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (ps *ProjectUsecase) UpdateProject(projectUuidStr, projectName string) (*models.Project, error) {
	project, err := ps.GetProjectByUuid(projectUuidStr)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
		return nil, CustomErrors.ErrForbidden
	}
	ctx := context.Background()
	project.Name = projectName
	if err := ps.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (ps *ProjectUsecase) DeleteProject(projectUuidStr string) error {
	projectUuid, err := uuid.Parse(projectUuidStr)
	if err != nil {
		return CustomErrors.ErrProjectUuidInvalid
	}
	ctx := context.Background()
	project, err := ps.projectRepo.GetByUuid(ctx, projectUuid)
	if err != nil {
		return err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
		return CustomErrors.ErrForbidden
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

func NewProject(customerUuid uuid.UUID, name string) *models.Project {
	return &models.Project{
		Uuid:         uuid.New(),
		Token:        generateToken(),
		CustomerUUID: customerUuid,
		Name:         name,
	}
}

func generateToken() string {
	rand.Seed(time.Now().Unix())
	sum := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	return fmt.Sprintf("%x", sum)
}
