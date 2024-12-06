package generic

import (
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

// Versioned is a types.Typed that can be decoded
// with the help of a registry.Registry, thus making it generic.
type Versioned struct {
	types.Typed `json:"-" yaml:"-"`
}
