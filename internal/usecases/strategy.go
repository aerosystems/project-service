package usecases

import "github.com/google/uuid"

type Strategy interface {
	IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool
	IsAccessibleByCountProjects(countProjects int) bool
}

type StartupStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (ss *StartupStrategy) IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool {
	if ss.strategyOwnerUuid != customerUuid {
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

func (bs *BusinessStrategy) IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool {
	if bs.strategyOwnerUuid != customerUuid {
		return false
	}
	return true
}

func (bs *BusinessStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}

type FreeStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (fs *FreeStrategy) IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool {
	if fs.strategyOwnerUuid != customerUuid {
		return false
	}
	return true
}

func (fs *FreeStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	if countProjects > 1 {
		return false
	}
	return true
}

type ServiceStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (ss *ServiceStrategy) IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool {
	return true
}

func (ss *ServiceStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}

type StaffStrategy struct {
	strategyOwnerUuid uuid.UUID
}

func (sf *StaffStrategy) IsAccessibleByCustomerUuid(customerUuid uuid.UUID) bool {
	return true
}

func (sf *StaffStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}
