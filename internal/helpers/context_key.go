package helpers

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}
