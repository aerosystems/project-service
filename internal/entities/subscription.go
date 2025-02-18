package entities

type SubscriptionType struct {
	slug string
}

var (
	UnknownSubscription  = SubscriptionType{"unknown"}
	TrialSubscription    = SubscriptionType{"trial"}
	StartupSubscription  = SubscriptionType{"startup"}
	BusinessSubscription = SubscriptionType{"business"}
)

func (k SubscriptionType) String() string {
	return k.slug
}

func NewSubscriptionType(kind string) SubscriptionType {
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
