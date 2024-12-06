package registry

import (
	"context"
	"fmt"
	"maps"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

// constructorContextKey is a unique key type for the constructor map in the context.
// used to avoid collisions with context values from other packages.
type constructorContextKey string

// ConstructorContextKey serves as a unique identifier for the constructor map in the context.
const ConstructorContextKey = constructorContextKey("type.registry")

// Inject injects the registry into the context using a map for O(1) lookups.
func Inject(ctx context.Context, registry *TypeRegistry) context.Context {
	if value := ctx.Value(ConstructorContextKey); value != nil {
		if mapType, ok := value.(*TypeRegistry); ok {
			maps.Copy(mapType.constructors, registry.constructors)
			return ctx
		}
	}

	return context.WithValue(ctx, ConstructorContextKey, registry)
}

func FromContext(ctx context.Context) (*TypeRegistry, error) {
	value := ctx.Value(ConstructorContextKey)
	if value == nil {
		return nil, fmt.Errorf("no registry available in ctx under %s", ConstructorContextKey)
	}

	reg, ok := value.(*TypeRegistry)
	if !ok {
		return nil, fmt.Errorf("registry in ctx at %s is not a valid TypeRegistry", ConstructorContextKey)
	}
	return reg, nil
}

// ConstructorFromContext retrieves a constructor by alias in O(1) time.
func ConstructorFromContext(ctx context.Context, alias string) (func() types.Typed, error) {
	reg, err := FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get constructor from context: %w", err)
	}

	reg.mu.RLock()
	defer reg.mu.RUnlock()
	constructor, exists := reg.constructors[alias]
	if !exists {
		return nil, fmt.Errorf("there is constructor available for alias %s", alias)
	}

	return constructor, nil
}
