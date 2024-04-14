package models

type KindSubscription struct {
	slug string
}

var (
	UnknownSubscription  = KindSubscription{"unknown"}
	TrialSubscription    = KindSubscription{"trial"}
	StartupSubscription  = KindSubscription{"startup"}
	BusinessSubscription = KindSubscription{"business"}
)

func (k KindSubscription) String() string {
	return k.slug
}

func NewKindSubscription(kind string) KindSubscription {
	switch kind {
	case TrialSubscription.String():
		return TrialSubscription
	case StartupSubscription.String():
		return StartupSubscription
	case BusinessSubscription.String():
		return BusinessSubscription
	default:
		return UnknownSubscription
	}
}
