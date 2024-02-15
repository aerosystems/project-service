package models

type KindRole string

const (
	CustomerRole KindRole = "customer"
	StaffRole    KindRole = "staff"
)

func (k KindRole) String() string {
	return string(k)
}
