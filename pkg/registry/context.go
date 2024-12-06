package registry

import (
	"context"
	"maps"
	"sync"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

type constructorMapType struct {
	mu           sync.RWMutex
	constructors map[string]func() types.Typed
}

// constructorContextKey is a unique key type for the constructor map in the context.
// used to avoid collisions with context values from other packages.
type constructorContextKey string

// ConstructorContextKeyPrefix serves as a unique identifier for the constructor map in the context.
const ConstructorContextKeyPrefix = constructorContextKey("type.registry.constructorMap")

// Inject injects the registry into the context using a map for O(1) lookups.
func Inject(ctx context.Context, registry *TypeRegistry) context.Context {
	if value := ctx.Value(ConstructorContextKeyPrefix); value != nil {
		if mapType, ok := value.(*constructorMapType); ok {
			maps.Copy(mapType.constructors, registry.constructors)
			return ctx
		}
	}

	return context.WithValue(ctx, ConstructorContextKeyPrefix, &constructorMapType{
		constructors: registry.constructors,
	})
}

// ConstructorFromContext retrieves a constructor by alias in O(1) time.
func ConstructorFromContext(ctx context.Context, alias string) (func() types.Typed, bool) {
	value := ctx.Value(ConstructorContextKeyPrefix)
	if value == nil {
		return nil, false
	}

	constructorMap, ok := value.(*constructorMapType)
	if !ok {
		return nil, false
	}

	constructorMap.mu.RLock()
	defer constructorMap.mu.RUnlock()

	constructor, exists := constructorMap.constructors[alias]
	return constructor, exists
}
