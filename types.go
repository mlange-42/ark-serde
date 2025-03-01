package arkserde

import "github.com/mlange-42/ark/ecs"

type deserializer struct {
	World      ecs.EntityDump
	Types      []string
	Components []entry
	Resources  map[string]entry
}

type entry struct {
	Bytes []byte
}

func (e *entry) UnmarshalJSON(jsonData []byte) error {
	e.Bytes = jsonData
	return nil
}
