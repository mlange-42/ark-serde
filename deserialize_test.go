package arkserde_test

import (
	"fmt"
	"testing"

	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeSkipEntities(t *testing.T) {
	jsonData, _, _, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w, arkserde.Opts.SkipEntities())
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 0)
	query.Close()

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestDeserializeSkipAllComponents(t *testing.T) {
	jsonData, parent, child, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)
	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w, arkserde.Opts.SkipAllComponents())
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 3)
	query.Close()

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestDeserializeSkipComponents(t *testing.T) {
	jsonData, parent, child, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	posId := ecs.ComponentID[Position](&w)
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)

	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w, arkserde.Opts.SkipComponents(ecs.C[Position]()))
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 3)

	query.Next()
	assert.False(t, query.Has(posId))
	assert.False(t, query.Has(velId))

	query.Next()
	assert.False(t, query.Has(posId))
	assert.False(t, query.Has(velId))

	query.Next()
	assert.False(t, query.Has(posId))
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*ChildOf)(query.Get(childId)), ChildOf{Entity: parent})

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestDeserializeSkipAllResources(t *testing.T) {
	jsonData, _, _, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)

	err = arkserde.Deserialize(jsonData, &w, arkserde.Opts.SkipAllResources())
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}
}

func TestDeserializeSkipResources(t *testing.T) {
	jsonData, _, _, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)

	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w, arkserde.Opts.SkipResources(ecs.C[Position]()))
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestDeserializeErrors(t *testing.T) {
	world := createWorld(false)

	err := arkserde.Deserialize([]byte("{xxx}"), world)
	assert.Contains(t, err.Error(), "invalid character 'x'")

	err = arkserde.Deserialize([]byte(textOk), world)
	assert.Contains(t, err.Error(), "component type is not registered")

	world = createWorld(true)

	err = arkserde.Deserialize([]byte(textOk), world)
	assert.Contains(t, err.Error(), "resource type is not registered")

	world = createWorld(true)
	_ = ecs.ResourceID[Velocity](world)

	err = arkserde.Deserialize([]byte(textOk), world)
	assert.Contains(t, err.Error(), "resource type registered but nil")

	world = createWorld(true)
	velAccess := ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textOk), world)
	assert.Nil(t, err)

	world = createWorld(true)
	velAccess = ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textErrEntities), world)
	assert.Contains(t, err.Error(), "world has 2 alive entities")

	world = createWorld(true)
	velAccess = ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textErrTypes), world)
	assert.Contains(t, err.Error(), "cannot unmarshal object")

	world = createWorld(true)
	velAccess = ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textErrComponent), world)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world = createWorld(true)
	velAccess = ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textErrComponent2), world)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world = createWorld(true)
	velAccess = ecs.NewResource[Velocity](world)
	velAccess.Add(&Velocity{})
	err = arkserde.Deserialize([]byte(textErrResource), world)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world = createWorld(true)
	err = arkserde.Deserialize([]byte(textErrRelation), world)
	assert.Contains(t, err.Error(), "cannot unmarshal object into Go value of type [2]uint32")
}

func createWorld(vel bool) *ecs.World {
	world := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&world)
	_ = ecs.ComponentID[ChildOf](&world)
	_ = ecs.ComponentID[ChildRelation](&world)
	if vel {
		_ = ecs.ComponentID[Velocity](&world)
	}
	return &world
}

const textOk = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : [
	  "arkserde_test.Velocity",
	  "arkserde_test.ChildOf",
	  "arkserde_test.Position"
	],
	"Components" : [
	  {
		"arkserde_test.Position" : {"X":1,"Y":2}
	  },
	  {
		"arkserde_test.Position" : {"X":3,"Y":4},
		"arkserde_test.Velocity" : {"X":5,"Y":6},
		"arkserde_test.ChildOf" : {"Entity":[1,0]}
	  }
	],
	"Resources" : {
		"arkserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrEntities = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : [
		"arkserde_test.Velocity",
		"arkserde_test.ChildOf",
		"arkserde_test.Position"
	],
	"Components" : [
		{
		"arkserde_test.Position" : {"X":3,"Y":4},
		"arkserde_test.Velocity" : {"X":5,"Y":6},
		"arkserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"arkserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrTypes = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : {"a": "b"},
	"Components" : [
		{
		  "arkserde_test.Position" : {"X":1,"Y":2}
		},
		{
		"arkserde_test.Position" : {"X":3,"Y":4},
		"arkserde_test.Velocity" : {"X":5,"Y":6},
		"arkserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"arkserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrComponent = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : [
		"arkserde_test.Velocity",
		"arkserde_test.ChildOf",
		"arkserde_test.Position"
	],
	"Components" : [
		[],
		{
		"arkserde_test.Position" : {"X":3,"Y":4},
		"arkserde_test.Velocity" : {"X":5,"Y":6},
		"arkserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"arkserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrComponent2 = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : [
		"arkserde_test.Velocity",
		"arkserde_test.ChildOf",
		"arkserde_test.Position"
	],
	"Components" : [
		{
		  "arkserde_test.Position" : []
		},
		{
		"arkserde_test.Position" : {"X":3,"Y":4},
		"arkserde_test.Velocity" : {"X":5,"Y":6},
		"arkserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"arkserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrRelation = `{
	"World" : {"Entities":[[0,4294967295],[1,4294967295],[2,0],[3,0]],"Alive":[2,3],"Next":0,"Available":0},
	"Types" : [
	  "arkserde_test.Position",
	  "arkserde_test.ChildRelation"
	],
	"Components" : [
	  {
		"arkserde_test.Position" : {"X":1,"Y":2}
	  },
	  {
		"arkserde_test.Position" : {"X":5,"Y":6},
		"arkserde_test.ChildRelation.ark.relation.Target" : {},
		"arkserde_test.ChildRelation" : {"Dummy":0}
	  }
	],
	"Resources" : {
	}}`

const textErrResource = `{
	"World" : {"Entities":[[0,4294967295], [1,4294967295]],"Alive":[],"Next":0,"Available":0},
	"Types" : [],
	"Components" : [],
	"Resources" : {
		"arkserde_test.Velocity" : []
	}}`

func benchmarkDeserializeJSON(n int, b *testing.B) {
	w := ecs.NewWorld(1024)

	mapper := ecs.NewMap2[Position, Velocity](&w)
	mapper.NewBatchFn(n, nil)

	jsonData, err := arkserde.Serialize(&w)
	if err != nil {
		panic(err.Error())
	}

	w2 := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w2)
	_ = ecs.ComponentID[Velocity](&w2)

	err = arkserde.Deserialize(jsonData, &w2)
	if err != nil {
		panic(err.Error())
	}
	w2.Reset()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = arkserde.Deserialize(jsonData, &w2)
		if err != nil {
			panic(err.Error())
		}
		b.StopTimer()
		w2.Reset()
		b.StartTimer()
	}
}

func BenchmarkDeserializeJSON_100(b *testing.B) {
	benchmarkDeserializeJSON(100, b)
}

func BenchmarkDeserializeJSON_1000(b *testing.B) {
	benchmarkDeserializeJSON(1000, b)
}

func BenchmarkDeserializeJSON_10000(b *testing.B) {
	benchmarkDeserializeJSON(10000, b)
}

func BenchmarkDeserializeJSON_100000(b *testing.B) {
	benchmarkDeserializeJSON(100000, b)
}
