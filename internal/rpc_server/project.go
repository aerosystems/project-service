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

type RPCProjectPayload struct {
	UserID     int
	UserRole   string
	Name       string
	AccessTime time.Time
}

func (r *ProjectServer) CreateProject(payload RPCProjectPayload, resp *string) error {
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
		UserID:     payload.UserID,
		Name:       payload.Name,
		AccessTime: payload.AccessTime,
	}

	if err = r.projectRepo.Create(&newProject); err != nil {
		return err
	}

	*resp = fmt.Sprintf("project %s successfully created", payload.Name)
	return nil
}
