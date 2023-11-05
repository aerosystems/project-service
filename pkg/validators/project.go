package validators

import (
	"errors"
	"github.com/aerosystems/project-service/internal/models"
)

func ValidateProject(project models.Project) error {
	if project.UserId == 0 {
		return errors.New("userID does not be empty")
	}

	if project.Name == "" {
		return errors.New("name does not be empty")
	}

	return nil
}
