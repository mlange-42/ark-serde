package arkserde

import (
	"testing"

	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

type CompA struct{}
type CompB struct{}

func TestBitMask(t *testing.T) {
	w := ecs.NewWorld(8)

	idA := ecs.ComponentID[CompA](&w)

	mask := bitMask{}

	assert.False(t, mask.Get(idA))
	mask.Set(idA, true)
	assert.True(t, mask.Get(idA))
	mask.Set(idA, false)
	assert.False(t, mask.Get(idA))
}
