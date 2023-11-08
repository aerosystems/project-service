package services

type Strategy interface {
	IsAccessible(userId int) bool
}

type StartupStrategy struct {
	strategyOwnerId int
}

func (ss *StartupStrategy) IsAccessible(userId int) bool {
	if ss.strategyOwnerId != userId {
		return false
	}
	return true
}

type BusinessStrategy struct {
	strategyOwnerId int
}

func (bs *BusinessStrategy) IsAccessible(userId int) bool {
	if bs.strategyOwnerId != userId {
		return false
	}
	return true
}

type StaffStrategy struct {
	strategyOwnerId int
}

func (sf *StaffStrategy) IsAccessible(userId int) bool {
	return true
}
