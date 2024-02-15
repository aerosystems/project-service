package usecases

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	RPCServices "github.com/aerosystems/project-service/internal/repository/rpc"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type ProjectUsecase struct {
	projectRepo ProjectRepository
	subsRPC     *RPCServices.SubscriptionRPC
	strategy    Strategy
}

func NewProjectUsecase(projectRepo ProjectRepository, subsRPC *RPCServices.SubscriptionRPC) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo: projectRepo,
		subsRPC:     subsRPC,
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
	kind, accessTime, err := ps.subsRPC.GetSubscription(userUuid)
	if err != nil {
		return err
	}
	// TODO: is it necessary?
	_ = accessTime
	//if accessTime.Before(time.Now()) {
	//	return errors.New("subscription is expired")
	//}
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
	project, err := ps.projectRepo.GetById(projectId)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectByToken(token string) (*models.Project, error) {
	project, err := ps.projectRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectUsecase) GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error) {
	if filterUserUuid != uuid.Nil {
		if !ps.strategy.IsAccessibleByUserUuid(filterUserUuid) {
			return []models.Project{}, nil
		}
		projectList, err = ps.projectRepo.GetByUserUuid(filterUserUuid)
	}
	projectList, err = ps.projectRepo.GetByUserUuid(userUuid)
	return projectList, nil
}

func (ps *ProjectUsecase) CreateDefaultProject(userUuid uuid.UUID) error {
	if err := ps.CreateProject(userUuid, "default"); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) CreateProject(userUuid uuid.UUID, name string) error {
	projectList, err := ps.projectRepo.GetByUserUuid(userUuid)
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
	if err := ps.projectRepo.Create(newProject); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) UpdateProject(project *models.Project) error {
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return errors.New("user is not allowed to update the project")
	}
	if err := ps.projectRepo.Update(project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) DeleteProjectById(projectId int) error {
	project, err := ps.projectRepo.GetById(projectId)
	if err != nil {
		return err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return errors.New("user is not allowed to delete the project")
	}
	if err := ps.projectRepo.Delete(project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectUsecase) IsProjectExistByToken(projectToken string) bool {
	project, err := ps.projectRepo.GetByToken(projectToken)
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
