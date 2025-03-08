package arkserde

import "github.com/mlange-42/ark/ecs"

const targetTag = ".ark.relation.Target"

type deserializer struct {
	World      ecs.EntityDump
	Types      []string
	Components []entry
	Resources  map[string]entry
}

type entry struct {
	Bytes []byte
}

func (e *entry) UnmarshalJSON(jsonData []byte) error {
	e.Bytes = jsonData
	return nil
}

type component struct {
	ID     ecs.ID
	Comp   interface{}
	Target ecs.Entity
}

// bitMask is a 256 bit bit-mask.
// It is also a [Filter] for including certain components.
type bitMask struct {
	bits [4]uint64 // 4x 64 bits of the mask
}

// Get reports whether the bit at the given index [ID] is set.
func (b *bitMask) Get(bit ecs.ID) bool {
	id := bit.Index()
	idx := id / 64
	offset := id - (64 * idx)
	mask := uint64(1 << offset)
	return b.bits[idx]&mask == mask
}

// Set sets the state of the bit at the given index.
func (b *bitMask) Set(bit ecs.ID, value bool) {
	id := bit.Index()
	idx := id / 64
	offset := id - (64 * idx)
	if value {
		b.bits[idx] |= (1 << offset)
	} else {
		b.bits[idx] &= ^(1 << offset)
	}
}
