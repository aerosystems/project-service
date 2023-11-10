package services

type Strategy interface {
	IsAccessibleByUserId(userId int) bool
	IsAccessibleByCountProjects(countProjects int) bool
}

type StartupStrategy struct {
	strategyOwnerId int
}

func (ss *StartupStrategy) IsAccessibleByUserId(userId int) bool {
	if ss.strategyOwnerId != userId {
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
	strategyOwnerId int
}

func (bs *BusinessStrategy) IsAccessibleByUserId(userId int) bool {
	if bs.strategyOwnerId != userId {
		return false
	}
	return true
}

func (bs *BusinessStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}

type StaffStrategy struct {
	strategyOwnerId int
}

func (sf *StaffStrategy) IsAccessibleByUserId(userId int) bool {
	return true
}

func (sf *StaffStrategy) IsAccessibleByCountProjects(countProjects int) bool {
	return true
}
