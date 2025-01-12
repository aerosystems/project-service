package usecases

import (
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"time"
)

const (
	defaultProjectName = "Default Project"
)

type ProjectUsecase struct {
	projectRepo            ProjectRepository
	subscriptionAdapter    SubscriptionAdapter
	checkmailEventsAdapter CheckmailEventsAdapter
	strategy               Strategy
}

func NewProjectUsecase(projectRepo ProjectRepository, subscriptionAdapter SubscriptionAdapter) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo:         projectRepo,
		subscriptionAdapter: subscriptionAdapter,
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
	kind, _, err := ps.subscriptionAdapter.GetSubscription(context.TODO(), customerUuid)
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
	// TODO: if its statement is needed, we should determine strategy before
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
	ps.SetStrategy(&ServiceStrategy{})
	project, err := ps.projectRepo.GetByCustomerUuidAndName(context.Background(), customerUuid, defaultProjectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return ps.CreateProject(customerUuid, defaultProjectName)
	}
	return project, nil
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

	newDefaultProject := models.NewProject(customerUuid, defaultProjectName)
	if err = ps.projectRepo.Create(ctx, newDefaultProject); err != nil {
		return nil, err
	}

	project, err := ps.projectRepo.GetByCustomerUuidAndName(ctx, customerUuid, defaultProjectName)
	if err != nil {
		return nil, err
	}
	if err = ps.checkmailEventsAdapter.PublishCreateAccessEvent(project.Token, subscriptionType, accessTime); err != nil {
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
	project := models.NewProject(customerUuid, name)
	ctx = context.Background()
	if err = ps.projectRepo.Create(ctx, project); err != nil {
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
	if err = ps.projectRepo.Update(ctx, project); err != nil {
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
	if err = ps.projectRepo.Delete(ctx, project); err != nil {
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
