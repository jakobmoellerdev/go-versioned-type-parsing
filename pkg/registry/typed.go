package registry

import (
	"fmt"
	"sync"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/unstructured"
)

type Registry interface {
	New(typ types.VersionedString) (types.Typed, error)
}

// TypeRegistry is a dynamic registry for Typed types.
type TypeRegistry struct {
	mu sync.RWMutex
	// allowUnknown allows unknown types to be created.
	// if the constructors cannot determine a match,
	// this will trigger the creation of an unstructured.Unstructured with New instead of failing.
	allowUnknown bool
	constructors map[string]func() types.Typed // Maps types names to constructor functions
}

type TypeRegistryOption func(*TypeRegistry)

// WithAllowUnknown allows unknown types to be created.
func WithAllowUnknown(allowUnknown bool) TypeRegistryOption {
	return func(registry *TypeRegistry) {
		registry.allowUnknown = allowUnknown
	}
}

// NewTypeRegistry creates a new registry.
func NewTypeRegistry(opts ...TypeRegistryOption) *TypeRegistry {
	reg := &TypeRegistry{
		constructors: make(map[string]func() types.Typed),
	}
	for _, opt := range opts {
		opt(reg)
	}
	return reg
}

// New instantiates types.Typed based on its constructors.
func (registry *TypeRegistry) New(typ types.VersionedString) (types.Typed, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	// construct by full type
	construct, exists := registry.constructors[typ.GetType().String()]
	if exists {
		return construct(), nil
	}

	// construct by base (no version)
	construct, exists = registry.constructors[typ.GetBase()]
	if exists {
		return construct(), nil
	}

	if registry.allowUnknown {
		return unstructured.New(), nil
	}

	return nil, fmt.Errorf("unsupported example types: %s", typ)
}

// Register registers a new Typed type in the registry for a given alias.
// The alias is used to identify the type in the registry.
// The constructor function is used to instantiate the type.
// In case the type was already registered under that alias,
// the new constructor will override the old one.
func Register[T any, P interface {
	*T
	types.Typed
}](registry *TypeRegistry, alias ...types.Typed) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	construct := func() types.Typed {
		return P(new(T))
	}
	for _, alias := range alias {
		for _, alias := range []string{
			alias.GetType().String(),
			alias.GetType().GetBase(),
		} {
			registry.constructors[alias] = construct
		}
	}
}
