package RPCServer

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"time"
)

type ProjectServer struct {
	projectRepo models.ProjectRepository
}

func NewProjectServer(projectRepo models.ProjectRepository) *ProjectServer {
	return &ProjectServer{projectRepo: projectRepo}
}

type CreateProjectRPCPayload struct {
	UserID     int
	UserRole   string
	Name       string
	AccessTime time.Time
}

type ProjectRPCPayload struct {
	ID         int
	UserID     int
	Name       string
	Token      string
	AccessTime time.Time
}

func (r *ProjectServer) CreateProject(payload CreateProjectRPCPayload, resp *string) error {
	projectList, err := r.projectRepo.FindByUserID(payload.UserID)
	if err != nil {
		return err
	}

	if len(projectList) > 0 && payload.UserRole == "startup" {
		return errors.New("user with Startup plan already has project, for create more projects you should switch into Business plan")
	}

	for _, project := range projectList {
		if project.Name == payload.Name {
			err := fmt.Errorf("project with Name %s already exists", payload.Name)
			return err
		}
	}

	var newProject = models.Project{
		UserID: payload.UserID,
		Name:   payload.Name,
	}

	if err = r.projectRepo.Create(&newProject); err != nil {
		return err
	}

	*resp = fmt.Sprintf("project %s successfully created", payload.Name)
	return nil
}

func (r *ProjectServer) GetProject(projectToken string, resp *ProjectRPCPayload) error {
	project, err := r.projectRepo.FindByToken(projectToken)
	if err != nil {
		return err
	}

	*resp = ProjectRPCPayload{
		ID:     project.ID,
		UserID: project.UserID,
		Name:   project.Name,
		Token:  project.Token,
	}
	return nil
}

func (r *ProjectServer) GetProjectList(userID int, resp *[]ProjectRPCPayload) error {
	projectList, err := r.projectRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	for _, project := range projectList {
		*resp = append(*resp, ProjectRPCPayload{
			ID:     project.ID,
			UserID: project.UserID,
			Name:   project.Name,
			Token:  project.Token,
		})
	}
	return nil
}
