package generators

import "github.com/kluctl/kluctl/v2/pkg/utils/uo"

type Generator interface {
	BuildContexts() ([]*GeneratedContext, error)
}

type GeneratedContext struct {
	Vars *uo.UnstructuredObject
}
