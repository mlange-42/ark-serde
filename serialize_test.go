package arkserde_test

import (
	"fmt"
	"testing"

	arkserde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type ChildOf struct {
	Entity ecs.Entity
}

type ChildRelation struct {
	ecs.RelationMarker
	Dummy int
}

type Generic[T any] struct {
	Value T
}

func serialize(opts ...arkserde.Option) ([]byte, ecs.Entity, ecs.Entity, error) {
	w := ecs.NewWorld(1024)
	u := w.Unsafe()

	posId := ecs.ComponentID[Position](&w)
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)

	parent := u.NewEntity(posId)
	*(*Position)(u.Get(parent, posId)) = Position{X: 1, Y: 2}

	child := u.NewEntity(posId, velId, childId)
	*(*Position)(u.Get(child, posId)) = Position{X: 3, Y: 4}
	*(*Velocity)(u.Get(child, velId)) = Velocity{X: 5, Y: 6}
	*(*ChildOf)(u.Get(child, childId)) = ChildOf{Entity: parent}

	u.NewEntity()

	resId := ecs.ResourceID[Velocity](&w)
	resId2 := ecs.ResourceID[Position](&w)
	w.Resources().Add(resId, &Velocity{X: 1000, Y: 0})
	w.Resources().Add(resId2, &Position{X: 1000, Y: 0})

	js, err := arkserde.Serialize(&w, opts...)
	return js, parent, child, err
}

func TestSerialize(t *testing.T) {
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

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 3)

	query.Next()
	assert.False(t, query.Has(posId))
	assert.False(t, query.Has(velId))

	query.Next()
	assert.True(t, query.Has(posId))
	assert.False(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 1, Y: 2})

	query.Next()
	assert.True(t, query.Has(posId))
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 3, Y: 4})
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*ChildOf)(query.Get(childId)), ChildOf{Entity: parent})

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestSerializeSkipEntities(t *testing.T) {
	jsonData, _, _, err := serialize(arkserde.Opts.SkipEntities())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)

	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 0)
	query.Close()

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestSerializeSkipAllComponents(t *testing.T) {
	jsonData, parent, child, err := serialize(arkserde.Opts.SkipAllComponents())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w)
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

func TestSerializeSkipComponents(t *testing.T) {
	jsonData, parent, child, err := serialize(arkserde.Opts.SkipComponents(ecs.C[Position]()))

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)

	ecs.AddResource(&w, &Position{})
	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := ecs.NewUnsafeFilter(&w).Query()

	assert.Equal(t, query.Count(), 3)

	query.Next()
	assert.False(t, query.Has(velId))

	query.Next()
	assert.False(t, query.Has(velId))

	query.Next()
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*ChildOf)(query.Get(childId)), ChildOf{Entity: parent})

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestSerializeSkipAllResources(t *testing.T) {
	jsonData, _, _, err := serialize(arkserde.Opts.SkipAllResources())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}
}

func TestSerializeSkipResources(t *testing.T) {
	jsonData, _, _, err := serialize(arkserde.Opts.SkipResources(ecs.C[Position]()))

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)

	ecs.AddResource(&w, &Velocity{})

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	res := ecs.GetResource[Velocity](&w)
	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestSerializeRelation(t *testing.T) {
	w := ecs.NewWorld(1024)
	u := w.Unsafe()

	posId := ecs.ComponentID[Position](&w)
	relId := ecs.ComponentID[ChildRelation](&w)

	parent := u.NewEntity(posId)
	*(*Position)(u.Get(parent, posId)) = Position{X: 1, Y: 2}

	child1 := u.NewEntityRel([]ecs.ID{posId, relId}, ecs.RelID(relId, ecs.Entity{}))
	*(*Position)(u.Get(child1, posId)) = Position{X: 3, Y: 4}
	*(*ChildRelation)(u.Get(child1, relId)) = ChildRelation{}

	child2 := u.NewEntityRel([]ecs.ID{posId, relId}, ecs.RelID(relId, ecs.Entity{}))
	*(*Position)(u.Get(child2, posId)) = Position{X: 5, Y: 6}
	*(*ChildRelation)(u.Get(child2, relId)) = ChildRelation{}

	u.SetRelations(child2, ecs.RelID(relId, parent))

	jsonData, err := arkserde.Serialize(&w)
	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}
	fmt.Println(string(jsonData))

	w = ecs.NewWorld(1024)
	_ = ecs.ComponentID[Position](&w)
	relId = ecs.ComponentID[ChildRelation](&w)

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	assert.Equal(t, u.GetRelation(child1, relId), ecs.Entity{})
	assert.Equal(t, u.GetRelation(child2, relId), parent)
}

func TestSerializeGeneric(t *testing.T) {
	w := ecs.NewWorld(1024)
	u := w.Unsafe()

	gen1Id := ecs.ComponentID[Generic[int32]](&w)
	gen2Id := ecs.ComponentID[Generic[float32]](&w)

	e1 := u.NewEntity(gen1Id)
	*(*Generic[int32])(u.Get(e1, gen1Id)) = Generic[int32]{Value: 1}

	e2 := u.NewEntity(gen2Id)
	*(*Generic[float32])(u.Get(e2, gen2Id)) = Generic[float32]{Value: 2.0}

	e3 := u.NewEntity(gen1Id, gen2Id)
	*(*Generic[int32])(u.Get(e3, gen1Id)) = Generic[int32]{Value: 3}
	*(*Generic[float32])(u.Get(e3, gen2Id)) = Generic[float32]{Value: 4.0}

	jsonData, err := arkserde.Serialize(&w)
	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}
	fmt.Println(string(jsonData))

	w = ecs.NewWorld(1024)
	_ = ecs.ComponentID[Generic[int32]](&w)
	_ = ecs.ComponentID[Generic[float32]](&w)

	err = arkserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	intMap := ecs.NewMap[Generic[int32]](&w)
	floatMap := ecs.NewMap[Generic[float32]](&w)

	assert.Equal(t, Generic[int32]{Value: 1}, *intMap.Get(e1))
	assert.False(t, floatMap.Has(e1))

	assert.False(t, intMap.Has(e2))
	assert.Equal(t, Generic[float32]{Value: 2.0}, *floatMap.Get(e2))

	assert.Equal(t, Generic[int32]{Value: 3}, *intMap.Get(e3))
	assert.Equal(t, Generic[float32]{Value: 4.0}, *floatMap.Get(e3))

	_, _, _ = e1, e2, e3
}

func benchmarkSerializeJSON(n int, b *testing.B) {
	w := ecs.NewWorld(1024)

	mapper := ecs.NewMap2[Position, Velocity](&w)
	mapper.NewBatchFn(n, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := arkserde.Serialize(&w)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkSerializeJSON_100(b *testing.B) {
	benchmarkSerializeJSON(100, b)
}

func BenchmarkSerializeJSON_1000(b *testing.B) {
	benchmarkSerializeJSON(1000, b)
}

func BenchmarkSerializeJSON_10000(b *testing.B) {
	benchmarkSerializeJSON(10000, b)
}

func BenchmarkSerializeJSON_100000(b *testing.B) {
	benchmarkSerializeJSON(100000, b)
}
