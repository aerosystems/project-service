package CustomErrors

type EchoHandlerMode struct {
	slug string
}

func (m EchoHandlerMode) String() string {
	return m.slug
}

var (
	ProductionMode  = EchoHandlerMode{"prod"}
	DevelopmentMode = EchoHandlerMode{"dev"}
)

func NewEchoHandlerMode(slug string) EchoHandlerMode {
	switch slug {
	case ProductionMode.slug:
		return ProductionMode
	default:
		return DevelopmentMode
	}
}
