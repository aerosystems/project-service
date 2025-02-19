package entities

type Role struct {
	slug string
}

var (
	UnknownRole  = Role{"unknown"}
	CustomerRole = Role{"customer"}
	StaffRole    = Role{"staff"}
)

func (k Role) String() string {
	return k.slug
}

func RoleFromString(kind string) Role {
	switch kind {
	case CustomerRole.String():
		return CustomerRole
	case StaffRole.String():
		return StaffRole
	default:
		return UnknownRole
	}
}
