package arkserde

import (
	"reflect"

	"github.com/mlange-42/ark/ecs"
)

// Opts is a helper to create Option instances.
var Opts = Options{}

// Option is an option. Modifies o.
// Create them using [Opts].
type Option func(o *serdeOptions)

// Options is a helper to create Option instances.
// Use it via the instance [Opts].
type Options struct{}

// SkipAllResources skips serialization or de-serialization of all resources.
func (o Options) SkipAllResources() Option {
	return func(o *serdeOptions) {
		o.skipAllResources = true
	}
}

// SkipAllComponents skips serialization or de-serialization of all components.
func (o Options) SkipAllComponents() Option {
	return func(o *serdeOptions) {
		o.skipAllComponents = true
	}
}

// SkipEntities skips serialization or de-serialization of all entities and components.
func (o Options) SkipEntities() Option {
	return func(o *serdeOptions) {
		o.skipEntities = true
	}
}

// SkipComponents skips serialization or de-serialization of certain components.
//
// When deserializing, the skipped components must still be registered.
func (o Options) SkipComponents(comps ...ecs.Comp) Option {
	return func(o *serdeOptions) {
		o.skipComponents = make([]reflect.Type, len(comps))
		for i, c := range comps {
			o.skipComponents[i] = c.Type()
		}
	}
}

// SkipResources skips serialization or de-serialization of certain resources.
//
// When deserializing, the skipped resources must still be registered.
func (o Options) SkipResources(comps ...ecs.Comp) Option {
	return func(o *serdeOptions) {
		o.skipResources = make([]reflect.Type, len(comps))
		for i, c := range comps {
			o.skipResources[i] = c.Type()
		}
	}
}

type serdeOptions struct {
	skipAllResources  bool
	skipAllComponents bool
	skipEntities      bool

	skipComponents []reflect.Type
	skipResources  []reflect.Type
}

func newSerdeOptions(opts ...Option) serdeOptions {
	o := serdeOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
