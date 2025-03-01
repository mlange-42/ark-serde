package arkserde_test

import (
	"fmt"

	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
)

func Example_options() {
	world := ecs.NewWorld(1024)
	builder := ecs.NewMap2[Position, Velocity](&world)
	builder.NewBatch(10, &Position{}, &Velocity{})

	// Serialize the world, skipping Velocity.
	jsonData, err := arkserde.Serialize(
		&world,
		arkserde.Opts.SkipComponents(ecs.C[Velocity]()),
	)
	if err != nil {
		fmt.Printf("could not serialize: %s\n", err)
		return
	}

	newWorld := ecs.NewWorld(1024)

	// Register required components and resources
	_ = ecs.ComponentID[Position](&newWorld)

	err = arkserde.Deserialize(jsonData, &newWorld)
	if err != nil {
		fmt.Printf("could not deserialize: %s\n", err)
		return
	}
}
