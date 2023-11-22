package services

import "github.com/google/uuid"

type Strategy interface {
	IsAccessibleByUserUuid(userUuid uuid.UUID) bool
	IsAccessibleByCountProjects(countProjects int) bool
}

type StartupStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (ss *StartupStrategy) IsAccessibleByUserUuid(userUuid uuid.UUID) bool {
	if ss.strategyOwnerUuid != userUuid {
		return false
	}
	return true
}

func (ss *StartupStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	if countProjects > 1 {
		return false
	}
	return true
}

type BusinessStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (bs *BusinessStrategy) IsAccessibleByUserUuid(userUuid uuid.UUID) bool {
	if bs.strategyOwnerUuid != userUuid {
		return false
	}
	return true
}

func (bs *BusinessStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}

type StaffStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (sf *StaffStrategy) IsAccessibleByUserUuid(userUuid uuid.UUID) bool {
	return true
}

func (sf *StaffStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}
