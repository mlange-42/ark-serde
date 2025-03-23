package arkserde

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/goccy/go-json"
	"github.com/mlange-42/ark/ecs"
)

// Serialize an Ark [ecs.World] to JSON.
//
// Serializes the following:
//   - Entities and the entity pool
//   - All components of all entities
//   - All resources
//
// All components and resources must be "JSON-able" with [encoding/json].
//
// The options can be used to skip some or all components,
// entities entirely, and/or some or all resources.
func Serialize(world *ecs.World, options ...Option) ([]byte, error) {
	opts := newSerdeOptions(options...)

	builder := strings.Builder{}

	builder.WriteString("{\n")

	if err := serializeWorld(world, &builder, &opts); err != nil {
		return nil, err
	}
	if !opts.skipEntities {
		builder.WriteString(",\n")
	}

	serializeTypes(world, &builder, &opts)
	builder.WriteString(",\n")

	if err := serializeComponents(world, &builder, &opts); err != nil {
		return nil, err
	}
	builder.WriteString(",\n")

	if err := serializeResources(world, &builder, &opts); err != nil {
		return nil, err
	}
	builder.WriteString("}\n")

	data := []byte(builder.String())
	if opts.compressed {
		var err error
		data, err = compressGZip(data, opts.compressionLevel)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func serializeWorld(world *ecs.World, builder *strings.Builder, opts *serdeOptions) error {
	if opts.skipEntities {
		return nil
	}

	entities := world.Unsafe().DumpEntities()

	jsonData, err := json.Marshal(entities)
	if err != nil {
		return err
	}
	builder.WriteString(fmt.Sprintf("\"World\" : %s", string(jsonData)))
	return nil
}

func serializeTypes(world *ecs.World, builder *strings.Builder, opts *serdeOptions) {
	if opts.skipEntities || opts.skipAllComponents {
		builder.WriteString("\"Types\" : []")
		return
	}

	builder.WriteString("\"Types\" : [\n")

	types := map[ecs.ID]reflect.Type{}

	allComps := ecs.ComponentIDs(world)
	for _, id := range allComps {
		if info, ok := ecs.ComponentInfo(world, id); ok {
			if !slices.Contains(opts.skipComponents, info.Type) {
				types[id] = info.Type
			}
		}
	}
	maxComp := len(types) - 1
	counter := 0
	for _, tp := range types {
		builder.WriteString(fmt.Sprintf("  \"%s\"", tp.String()))
		if counter < maxComp {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
		counter++
	}

	builder.WriteString("]")
}

func serializeComponents(world *ecs.World, builder *strings.Builder, opts *serdeOptions) error {
	if opts.skipEntities {
		builder.WriteString("\"Components\" : []")
		return nil
	}

	skipComponents := bitMask{}
	for _, tp := range opts.skipComponents {
		id := ecs.TypeID(world, tp)
		skipComponents.Set(id, true)
	}

	builder.WriteString("\"Components\" : [\n")

	query := ecs.NewUnsafeFilter(world).Query()
	lastEntity := query.Count() - 1
	counter := 0
	tempIDs := []ecs.ID{}
	for query.Next() {
		if opts.skipAllComponents {
			builder.WriteString("  {")
		} else {
			builder.WriteString("  {\n")

			ids := query.IDs()

			tempIDs = tempIDs[:0]
			for i := range ids.Len() {
				id := ids.Get(i)
				if !skipComponents.Get(id) {
					tempIDs = append(tempIDs, id)
				}
			}
			last := len(tempIDs) - 1

			for i, id := range tempIDs {
				info, _ := ecs.ComponentInfo(world, id)

				if info.IsRelation {
					target := query.GetRelation(id)
					eJSON, err := target.MarshalJSON()
					if err != nil {
						return err
					}
					builder.WriteString(fmt.Sprintf("    \"%s%s\" : %s,\n", info.Type.String(), targetTag, eJSON))
				}

				comp := query.Get(id)
				value := reflect.NewAt(info.Type, comp).Interface()
				jsonData, err := json.Marshal(value)
				if err != nil {
					return err
				}
				builder.WriteString(fmt.Sprintf("    \"%s\" : ", info.Type.String()))
				builder.WriteString(string(jsonData))
				if i < last {
					builder.WriteString(",")
				}
				builder.WriteString("\n")
			}
		}
		builder.WriteString("  }")
		if counter < lastEntity {
			builder.WriteString(",")
		}
		builder.WriteString("\n")

		counter++
	}
	builder.WriteString("]")

	return nil
}

func serializeResources(world *ecs.World, builder *strings.Builder, opts *serdeOptions) error {
	if opts.skipAllResources {
		builder.WriteString("\"Resources\" : {}")
		return nil
	}

	builder.WriteString("\"Resources\" : {\n")

	resTypes := map[ecs.ResID]reflect.Type{}
	allRes := ecs.ResourceIDs(world)
	for _, id := range allRes {
		if tp, ok := ecs.ResourceType(world, id); ok {
			if !slices.Contains(opts.skipResources, tp) {
				resTypes[id] = tp
			}
		}
	}

	last := len(resTypes) - 1
	counter := 0
	for id, tp := range resTypes {
		res := world.Resources().Get(id)
		rValue := reflect.ValueOf(res)
		ptr := rValue.UnsafePointer()

		value := reflect.NewAt(tp, ptr).Interface()
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}

		builder.WriteString("    ")
		builder.WriteString(fmt.Sprintf("\"%s\" : ", tp.String()))
		builder.WriteString(string(jsonData))

		if counter < last {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
		counter++
	}

	builder.WriteString("}")

	return nil
}
