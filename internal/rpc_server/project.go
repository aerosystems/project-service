package RPCServer

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"gorm.io/gorm"
	"log"
	"time"
)

type ProjectServer struct {
	projectRepo models.ProjectRepository
}

func NewProjectServer(projectRepo models.ProjectRepository) *ProjectServer {
	return &ProjectServer{projectRepo: projectRepo}
}

type ProjectPayload struct {
	UserID     int
	Name       string
	AccessTime time.Time
}

func (r *ProjectServer) CreateProject(payload ProjectPayload, resp *string) error {
	log.Printf("received request to create project %s from %d", payload.Name, payload.UserID)

	project, err := r.projectRepo.FindByUserID(payload.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if project != nil {
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

	*resp = "OK"
	return nil
}
