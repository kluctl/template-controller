package generators

type Generator interface {
	BuildContexts() ([]*GeneratedContext, error)
}

type GeneratedContext struct {
	Vars map[string]any
}
