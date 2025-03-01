package arkserde

import (
	"reflect"
	"testing"

	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

type testComp struct{}

func TestOptions(t *testing.T) {
	opt := newSerdeOptions(
		Opts.SkipEntities(),
		Opts.SkipAllComponents(),
		Opts.SkipAllResources(),
		Opts.SkipComponents(ecs.C[testComp]()),
		Opts.SkipResources(ecs.C[testComp]()),
	)

	assert.True(t, opt.skipEntities)
	assert.True(t, opt.skipAllComponents)
	assert.True(t, opt.skipAllResources)
	assert.Equal(t, []reflect.Type{ecs.C[testComp]().Type()}, opt.skipComponents)
	assert.Equal(t, []reflect.Type{ecs.C[testComp]().Type()}, opt.skipResources)
}
