package arkserde

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/mlange-42/ark/ecs"
)

type component struct {
	ID     ecs.ID
	Comp   interface{}
	Target ecs.Entity
}

// Deserialize an Ark [ecs.World] from JSON.
//
// The world must be prepared the following way:
//   - The world must not contain any alive or dead entities (i.e. a new or [ecs.World.Reset] world)
//   - All required component types must be registered using [ecs.ComponentID]
//   - All required resources must be added as dummies using [ecs.AddResource]
//
// The options can be used to skip some or all components,
// entities entirely, and/or some or all resources.
// It only some components or resources are skipped,
// they still need to be registered to the world.
//
// # Query iteration order
//
// After deserialization, it is not guaranteed that entity iteration order in queries is the same as before.
// More precisely, it should at first be the same as before, but will likely deviate over time from what would
// happen when continuing the original, serialized run. Multiple worlds deserialized from the same source should,
// however, behave exactly the same.
func Deserialize(jsonData []byte, world *ecs.World, options ...Option) error {
	opts := newSerdeOptions(options...)

	deserial := deserializer{}
	if err := json.Unmarshal(jsonData, &deserial); err != nil {
		return err
	}

	if !opts.skipEntities {
		world.LoadEntities(&deserial.World)
	}

	if err := deserializeComponents(world, &deserial, &opts); err != nil {
		return err
	}
	if err := deserializeResources(world, &deserial, &opts); err != nil {
		return err
	}

	return nil
}

func deserializeComponents(world *ecs.World, deserial *deserializer, opts *serdeOptions) error {
	u := world.Unsafe()

	if opts.skipEntities {
		return nil
	}

	infos := map[ecs.ID]ecs.CompInfo{}
	ids := map[string]ecs.ID{}
	allComps := ecs.ComponentIDs(world)
	for _, id := range allComps {
		if info, ok := ecs.ComponentInfo(world, id); ok {
			infos[id] = info
			ids[info.Type.String()] = id
		}
	}

	for _, tp := range deserial.Types {
		if _, ok := ids[tp]; !ok {
			return fmt.Errorf("component type is not registered: %s", tp)
		}
	}

	if len(deserial.Components) != len(deserial.World.Alive) {
		return fmt.Errorf("found components for %d entities, but world has %d alive entities", len(deserial.Components), len(deserial.World.Alive))
	}

	skipComponents := ecs.Mask{}
	for _, tp := range opts.skipComponents {
		id := ecs.TypeID(world, tp)
		skipComponents.Set(id, true)
	}

	for i, comps := range deserial.Components {
		entity := deserial.World.Entities[deserial.World.Alive[i]]

		mp := map[string]entry{}

		if err := json.Unmarshal(comps.Bytes, &mp); err != nil {
			return err
		}

		target := ecs.Entity{}
		var targetComp ecs.ID
		hasRelation := false
		components := make([]component, 0, len(mp))
		compIDs := make([]ecs.ID, 0, len(mp))
		for tpName, value := range mp {
			if tpName == targetTag {
				if err := json.Unmarshal(value.Bytes, &target); err != nil {
					return err
				}
				continue
			}

			id := ids[tpName]
			if skipComponents.Get(id) {
				continue
			}

			info := infos[id]

			if info.IsRelation {
				targetComp = id
				hasRelation = true
			}

			comp := reflect.New(info.Type).Interface()
			if err := json.Unmarshal(value.Bytes, &comp); err != nil {
				return err
			}
			compIDs = append(compIDs, id)
			components = append(components, component{
				ID:   id,
				Comp: comp,
			})
		}

		if len(components) == 0 {
			continue
		}

		if !hasRelation {
			target = ecs.Entity{}
		}

		u.Add(entity, compIDs...)
		for _, comp := range components {
			assign(world, entity, comp.ID, comp.Comp)
		}
		if !target.IsZero() {
			u.SetRelations(entity, ecs.RelID(targetComp, target))
		}
	}
	return nil
}

func assign(world *ecs.World, entity ecs.Entity, id ecs.ID, comp interface{}) {
	dst := world.Unsafe().Get(entity, id)
	rValue := reflect.ValueOf(comp).Elem()

	valueType := rValue.Type()
	valuePtr := reflect.NewAt(valueType, dst)
	valuePtr.Elem().Set(rValue)
}

func deserializeResources(world *ecs.World, deserial *deserializer, opts *serdeOptions) error {
	if opts.skipAllResources {
		return nil
	}

	resTypes := map[ecs.ResID]reflect.Type{}
	resIds := map[string]ecs.ResID{}
	allRes := ecs.ResourceIDs(world)
	skipResources := ecs.Mask{}
	for _, id := range allRes {
		if tp, ok := ecs.ResourceType(world, id); ok {
			resTypes[id] = tp
			resIds[tp.String()] = id

			if slices.Contains(opts.skipResources, tp) {
				skipResources.Set(ecs.ID(id), true)
			}
		}
	}

	for tpName, res := range deserial.Resources {
		resID, ok := resIds[tpName]
		if !ok {
			return fmt.Errorf("resource type is not registered: %s", tpName)
		}
		if skipResources.Get(ecs.ID(resID)) {
			continue
		}

		tp := resTypes[resID]

		resLoc := world.Resources().Get(resID)
		if resLoc == nil {
			return fmt.Errorf("resource type registered but nil: %s", tpName)
		}

		ptr := reflect.ValueOf(resLoc).UnsafePointer()
		value := reflect.NewAt(tp, ptr).Interface()

		if err := json.Unmarshal(res.Bytes, &value); err != nil {
			return err
		}
	}
	return nil
}
