package models

type KindSubscription string

const (
	Startup  KindSubscription = "startup"
	Business KindSubscription = "business"
)

func (k KindSubscription) String() string {
	return string(k)
}
