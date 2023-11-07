package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	RPCServices "github.com/aerosystems/project-service/internal/rpc_services"
	"math/rand"
	"strconv"
	"time"
)

type ProjectService interface {
	DetermineStrategy(senderId int, role string, userId int) error
	CreateProject(userId int, name string) error
	CreateDefaultProject(userId int) error
	setStrategy(strategy Strategy)
	generateToken() string
	isNameExists(name string, projectList []models.Project) bool
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

func (ps *ProjectServiceImpl) DetermineStrategy(senderId int, role string, userId int) error {
	if role == "staff" {
		ps.setStrategy(&StaffStrategy{})
		return nil
	}
	if senderId != userId {
		return errors.New("can't create a project for another user")
	}
	kind, err := ps.subsRPC.GetSubscriptionKind(userId)
	if err != nil {
		return errors.New("failed to get subscription kind")
	}
	switch kind {
	case "startup":
		ps.setStrategy(&StartupStrategy{})
	case "business":
		ps.setStrategy(&BusinessStrategy{})
	default:
		return errors.New("unknown subscription kind")
	}
	return nil
}

func (ps *ProjectServiceImpl) CreateDefaultProject(userId int) error {
	if err := ps.CreateProject(userId, "default"); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServiceImpl) CreateProject(userId int, name string) error {
	projectList, err := ps.projectRepo.GetByUserId(userId)
	if err != nil {
		return err
	}
	if ps.isNameExists(name, projectList) {
		return errors.New("project name already exists")
	}
	if !ps.strategy.IsAccessible(userId, projectList) {
		return errors.New("user is not allowed to create a project")
	}
	var newProject = models.Project{
		Token:  ps.generateToken(),
		UserId: userId,
		Name:   name,
	}
	if err := ps.projectRepo.Create(&newProject); err != nil {
		return err
	}
	return nil
}

func (ps *ProjectServiceImpl) setStrategy(strategy Strategy) {
	ps.strategy = strategy
}

func (ps *ProjectServiceImpl) generateToken() string {
	rand.Seed(time.Now().Unix())
	sum := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	return fmt.Sprintf("%x", sum)
}

func (ps *ProjectServiceImpl) isNameExists(name string, projectList []models.Project) bool {
	for _, project := range projectList {
		if project.Name == name {
			return true
		}
	}
	return false
}
