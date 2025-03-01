package arkserde_test

import (
	"fmt"
	"math/rand"

	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
)

const (
	width  = 40
	height = 12
)

type Coords struct {
	X int
	Y int
}

func Example() {
	rng := rand.New(rand.NewSource(42))

	// Create a world.
	world := ecs.NewWorld(1024)

	// Populate the world with entities, components and resources.
	builder := ecs.NewMap1[Coords](&world)
	builder.NewBatchFn(60, func(entity ecs.Entity, coord *Coords) {
		coord.X = rng.Intn(width)
		coord.Y = rng.Intn(height)
	})

	// Print the original world
	fmt.Println("====== Original world ========")
	printWorld(&world)

	// Serialize the world.
	jsonData, err := arkserde.Serialize(&world)
	if err != nil {
		fmt.Printf("could not serialize: %s\n", err)
		return
	}

	// Print the resulting JSON.
	//fmt.Println(string(jsonData))

	// Create a new, empty world.
	newWorld := ecs.NewWorld(1024)

	// Register required components and resources
	_ = ecs.ComponentID[Coords](&newWorld)

	// Deserialize into the new world.
	err = arkserde.Deserialize(jsonData, &newWorld)
	if err != nil {
		fmt.Printf("could not deserialize: %s\n", err)
		return
	}

	// Print the deserialized world
	fmt.Println("====== Deserialized world ========")
	printWorld(&newWorld)
	// Output: ====== Original world ========
	// --------------------------------O-O---O-
	// -----------------------O----------------
	// -O-------------O------OO--------------O-
	// ----O------------------------OOO--------
	// O--------------OO-O---------------------
	// ------------O-----------O---------------
	// --O-------------O-------O---O------O----
	// O-O-----O----OOO-O--O--------------OO---
	// -----------O---OO----O--O------------O--
	// ------------O-----O---------------------
	// ---O---------------O------O--O----------
	// ------O-OO--O---------OO-OOO-----------O
	// ====== Deserialized world ========
	// --------------------------------O-O---O-
	// -----------------------O----------------
	// -O-------------O------OO--------------O-
	// ----O------------------------OOO--------
	// O--------------OO-O---------------------
	// ------------O-----------O---------------
	// --O-------------O-------O---O------O----
	// O-O-----O----OOO-O--O--------------OO---
	// -----------O---OO----O--O------------O--
	// ------------O-----O---------------------
	// ---O---------------O------O--O----------
	// ------O-OO--O---------OO-OOO-----------O
}

func printWorld(world *ecs.World) {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = '-'
		}
	}

	filter := ecs.NewFilter1[Coords](world)
	query := filter.Query()

	for query.Next() {
		coords := query.Get()
		grid[coords.Y][coords.X] = 'O'
	}

	for i := 0; i < len(grid); i++ {
		fmt.Println(string(grid[i]))
	}
}
