package validators

import (
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"time"
)

func ValidateProject(project models.Project) error {
	if project.UserID == 0 {
		return errors.New("userID does not be empty")
	}

	if project.Name == "" {
		return errors.New("name does not be empty")
	}

	if project.AccessTime.IsZero() {
		return errors.New("accessTime does not be empty")
	}

	if project.AccessTime.Before(time.Now()) {
		return errors.New("accessTime should be more then NOW")
	}

	return nil
}
