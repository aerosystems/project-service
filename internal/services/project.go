package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	RPCServices "github.com/aerosystems/project-service/internal/rpc_services"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type ProjectService interface {
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectById(projectId int) (*models.Project, error)
	GetProjectByToken(token string) (*models.Project, error)
	GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(userUuid uuid.UUID, name string) error
	CreateDefaultProject(userUuid uuid.UUID) error
	UpdateProject(project *models.Project) error
	DeleteProjectById(projectId int) error
	IsProjectExistByToken(projectToken string) bool
	SetStrategy(strategy Strategy)
	isProjectNameExist(name string, projectList []models.Project) bool
}

type ProjectServiceImpl struct {
	projectRepo models.ProjectRepository
	subsRPC     *RPCServices.SubscriptionRPC
	strategy    Strategy
}

func NewProjectServiceImpl(projectRepo models.ProjectRepository, subsRPC *RPCServices.SubscriptionRPC) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
		subsRPC:     subsRPC,
	}
}

func (ps *ProjectServiceImpl) DetermineStrategy(userUuidStr string, role string) error {
	userUuid, err := uuid.Parse(userUuidStr)
	if err != nil {
		return err
	}
	if role == "staff" {
		ps.SetStrategy(&StaffStrategy{userUuid})
		return nil
	}
	kind, err := ps.subsRPC.GetSubscriptionKind(userUuid)
	if err != nil {
		return errors.New("failed to get subscription kind")
	}
	switch kind {
	case "startup":
		ps.SetStrategy(&StartupStrategy{userUuid})
	case "business":
		ps.SetStrategy(&BusinessStrategy{userUuid})
	default:
		return errors.New("unknown subscription kind")
	}
	return nil
}

func (ps *ProjectServiceImpl) GetProjectById(projectId int) (*models.Project, error) {
	project, err := ps.projectRepo.GetById(projectId)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectServiceImpl) GetProjectByToken(token string) (*models.Project, error) {
	project, err := ps.projectRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return nil, errors.New("user is not allowed to access the project")
	}
	return project, nil
}

func (ps *ProjectServiceImpl) GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error) {
	if filterUserUuid != uuid.Nil {
		if !ps.strategy.IsAccessibleByUserUuid(filterUserUuid) {
			return []models.Project{}, nil
		}
		projectList, err = ps.projectRepo.GetByUserUuid(filterUserUuid)
	}
	projectList, err = ps.projectRepo.GetByUserUuid(userUuid)
	return projectList, nil
}

func (ps *ProjectServiceImpl) CreateDefaultProject(userUuid uuid.UUID) error {
	if err := ps.CreateProject(userUuid, "default"); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServiceImpl) CreateProject(userUuid uuid.UUID, name string) error {
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

func (ps *ProjectServiceImpl) UpdateProject(project *models.Project) error {
	if !ps.strategy.IsAccessibleByUserUuid(project.UserUuid) {
		return errors.New("user is not allowed to update the project")
	}
	if err := ps.projectRepo.Update(project); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServiceImpl) DeleteProjectById(projectId int) error {
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

func (ps *ProjectServiceImpl) IsProjectExistByToken(projectToken string) bool {
	project, err := ps.projectRepo.GetByToken(projectToken)
	if err != nil {
		return false
	}
	if project == nil {
		return false
	}
	return true
}

func (ps *ProjectServiceImpl) SetStrategy(strategy Strategy) {
	ps.strategy = strategy
}

func (ps *ProjectServiceImpl) isProjectNameExist(name string, projectList []models.Project) bool {
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
