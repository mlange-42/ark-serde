package main

// Profiling:
// go build ./serialize
// ./serialize
// go tool pprof -http=":8000" -nodefraction=0.001 ./serialize cpu.pprof

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

	world := ecs.NewWorld(1024)
	mapper := ecs.NewMap2[position, velocity](&world)
	mapper.NewBatchFn(entities, nil)

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(&world, iters)
	stop.Stop()
}

func run(world *ecs.World, iters int) {
	var jsonData []byte
	var err error
	for j := 0; j < iters; j++ {
		jsonData, err = arkserde.Serialize(world)
		if err != nil {
			panic(err.Error())
		}
	}
	if len(jsonData) == 0 {
		panic("check failed")
	}
}
