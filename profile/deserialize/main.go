package main

// Profiling:
// go build ./deserialize
// ./deserialize
// go tool pprof -http=":8000" -nodefraction=0.001 ./deserialize cpu.pprof

import (
	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/pkg/profile"
)

type position struct {
	X float64
	Y float64
}

type velocity struct {
	X float64
	Y float64
}

func main() {

	iters := 2500
	entities := 1000

	data := createData(entities)

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(data, iters)
	stop.Stop()
}

func run(data []byte, iters int) {
	world := ecs.NewWorld(1024)
	_ = ecs.ComponentID[position](&world)
	_ = ecs.ComponentID[velocity](&world)

	for j := 0; j < iters; j++ {
		err := arkserde.Deserialize(data, &world)
		if err != nil {
			panic(err.Error())
		}
		world.Reset()
	}
}

func createData(entities int) []byte {
	world := ecs.NewWorld(1024)
	mapper := ecs.NewMap2[position, velocity](&world)
	mapper.NewBatchFn(entities, nil)

	jsonData, err := arkserde.Serialize(&world)
	if err != nil {
		panic(err.Error())
	}
	return jsonData
}
