package generic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

// TypedDecoder is a decoder that converts into a types.Typed
// based on underlying typing rules. As such it is more concrete
// than Decoder in that it decides on the type to parse and the instanciation
// itself
type TypedDecoder interface {
	Decode() (types.Typed, error)
}

// Decoder is any arbitrary object that can decode objects into a given
// struct through go tags. An example is encoding/json
type Decoder interface {
	Decode(v any) error
}

// NewTypedDecoder creates a new TypedDecoder backed by a registry.
func NewTypedDecoder(data io.Reader, reg registry.Registry, decoder func(io.Reader) Decoder) TypedDecoder {
	return &RegistryDecoder{
		decoder:  decoder,
		reader:   data,
		registry: reg,
	}
}

// RegistryDecoder is a TypedDecoder backed by a registry.Registry
// to do typing decisions
type RegistryDecoder struct {
	reader   io.Reader
	decoder  func(io.Reader) Decoder
	registry registry.Registry
}

func (d *RegistryDecoder) Decode() (types.Typed, error) {
	// buffer the data so we can read it twice
	reader := d.reader
	var buf bytes.Buffer
	reader = io.TeeReader(reader, &buf)
	decoder := d.decoder(reader)

	// Extract the type to decide the concrete implementation
	var raw types.VersionedType
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}
	if raw.Type.GetBase() == "" {
		return nil, fmt.Errorf("missing or invalid 'type' in object that was expected to be typed")
	}

	instance, err := d.registry.New(raw.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to create new type %s: %w", raw.Type, err)
	}

	// reset the reader to the beginning of the data
	reader = io.MultiReader(&buf, reader)
	decoder = d.decoder(reader)

	// Decode the YAML node directly into the specific Typed instance
	if err := decoder.Decode(instance); err != nil {
		return nil, err
	}

	return instance, nil
}
