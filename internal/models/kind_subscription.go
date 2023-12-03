package models

type KindSubscription string

const (
	TrialSubscription    KindSubscription = "trial"
	StartupSubscription  KindSubscription = "startup"
	BusinessSubscription KindSubscription = "business"
)

func (k KindSubscription) String() string {
	return string(k)
}
