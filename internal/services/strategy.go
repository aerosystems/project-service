package services

import "github.com/aerosystems/project-service/internal/models"

type Strategy interface {
	IsAccessible(userId int, projectList []models.Project) bool
}

type StartupStrategy struct {
}

func (ss *StartupStrategy) IsAccessible(userId int, projectList []models.Project) bool {
	if len(projectList) > 1 {
		return false
	}
	return true
}

type BusinessStrategy struct {
}

func (bs *BusinessStrategy) IsAccessible(userId int, projectList []models.Project) bool {
	return true
}

type StaffStrategy struct {
}

func (sf *StaffStrategy) IsAccessible(userId int, projectList []models.Project) bool {
	return true
}
