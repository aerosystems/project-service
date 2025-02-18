package usecases

import (
	"context"
	"errors"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
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

func (ps *ProjectUsecase) SetStrategy(strategy Strategy) {
	ps.strategy = strategy
}

func (ps *ProjectUsecase) DetermineStrategy(ctx context.Context, userUUID uuid.UUID, role entities.Role) error {
	if role == entities.StaffRole {
		ps.SetStrategy(&StaffStrategy{userUUID})
		return nil
	}
	kind, _, err := ps.subscriptionAdapter.GetSubscription(context.TODO(), userUUID)
	if err != nil {
		return err
	}
	switch kind {
	case entities.TrialSubscription:
		fallthrough
	case entities.StartupSubscription:
		ps.SetStrategy(&StartupStrategy{userUUID})
	case entities.BusinessSubscription:
		ps.SetStrategy(&BusinessStrategy{userUUID})
	default:
		return entities.ErrForbidden
	}
	return nil
}

func (ps *ProjectUsecase) GetProjectByUuid(ctx context.Context, projectUuidStr string) (*entities.Project, error) {
	projectUuid, err := uuid.Parse(projectUuidStr)
	if err != nil {
		return nil, entities.ErrProjectUuidInvalid
	}
	project, err := ps.projectRepo.GetByUuid(ctx, projectUuid)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectByToken(ctx context.Context, token string) (*entities.Project, error) {
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

func (ps *ProjectUsecase) GetProjectListByCustomerUuid(ctx context.Context, customerUuid, filterUserUuid uuid.UUID) (projectList []entities.Project, err error) {
	if filterUserUuid != uuid.Nil {
		if !ps.strategy.IsAccessibleByCustomerUuid(filterUserUuid) {
			return []entities.Project{}, nil
		}
		projectList, err = ps.projectRepo.GetByCustomerUuid(ctx, filterUserUuid)
	}
	projectList, err = ps.projectRepo.GetByCustomerUuid(ctx, customerUuid)
	return projectList, nil
}

func (ps *ProjectUsecase) CreateDefaultProject(ctx context.Context, customerUuid uuid.UUID) (*entities.Project, error) {
	ps.SetStrategy(&ServiceStrategy{})
	project, err := ps.projectRepo.GetByCustomerUuidAndName(ctx, customerUuid, defaultProjectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return ps.CreateProject(ctx, customerUuid, defaultProjectName)
	}
	return project, nil
}

func (ps *ProjectUsecase) CreateProject(ctx context.Context, customerUuid uuid.UUID, name string) (*entities.Project, error) {
	projectList, err := ps.projectRepo.GetByCustomerUuid(ctx, customerUuid)
	if err != nil {
		return nil, err
	}
	if ps.isProjectNameExist(name, projectList) {
		return nil, entities.ErrProjectNameExists
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(customerUuid) {
		return nil, entities.ErrForbidden
	}
	if !ps.strategy.IsAccessibleByCountProjects(len(projectList)) {
		return nil, entities.ErrProjectLimitExceeded
	}
	project := entities.NewProject(customerUuid, name)
	if err = ps.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (ps *ProjectUsecase) UpdateProject(ctx context.Context, projectUuidStr, projectName string) (*entities.Project, error) {
	project, err := ps.GetProjectByUuid(ctx, projectUuidStr)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
		return nil, entities.ErrForbidden
	}
	project.Name = projectName
	if err = ps.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (ps *ProjectUsecase) DeleteProject(ctx context.Context, projectUuidStr string) error {
	projectUuid, err := uuid.Parse(projectUuidStr)
	if err != nil {
		return entities.ErrProjectUuidInvalid
	}
	project, err := ps.projectRepo.GetByUuid(ctx, projectUuid)
	if err != nil {
		return err
	}
	if !ps.strategy.IsAccessibleByCustomerUuid(project.CustomerUUID) {
		return entities.ErrForbidden
	}
	if err = ps.projectRepo.Delete(ctx, project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) IsProjectExistByToken(ctx context.Context, projectToken string) bool {
	project, err := ps.projectRepo.GetByToken(ctx, projectToken)
	if err != nil {
		return false
	}
	if project == nil {
		return false
	}
	return true
}

func (ps *ProjectUsecase) isProjectNameExist(name string, projectList []entities.Project) bool {
	for _, project := range projectList {
		if project.Name == name {
			return true
		}
	}
	return false
}
