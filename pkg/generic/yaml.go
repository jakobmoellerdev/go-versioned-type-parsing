package generic

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

var _ interface {
	yaml.BytesUnmarshalerContext
	yaml.BytesMarshalerContext
} = &Versioned{}

func (w *Versioned) UnmarshalYAML(ctx context.Context, data []byte) error {
	// buffer the data so we can read it twice
	reader := io.Reader(bytes.NewReader(data))
	var buf bytes.Buffer
	reader = io.TeeReader(reader, &buf)
	decoder := yaml.NewDecoder(reader)

	// Extract the "types" field to decide the concrete implementation
	var raw types.VersionedType
	if err := decoder.DecodeContext(ctx, &raw); err != nil {
		return err
	}
	if raw.Type.GetBase() == "" {
		return fmt.Errorf("missing or invalid 'type' in object that was expected to be typed")
	}

	newFn, ok := registry.ConstructorFromContext(ctx, raw.Type.String())
	if newFn == nil {
		return fmt.Errorf("type %s is not registered within context", raw.Type)
	}
	if !ok {
		return fmt.Errorf("unexpected function registered in type registry for %s: %T", raw.Type, newFn)
	}
	instance := newFn()

	// reset the reader to the beginning of the data
	reader = io.MultiReader(&buf, reader)
	decoder = yaml.NewDecoder(reader)

	// Decode the YAML node directly into the specific Typed instance
	if err := decoder.DecodeContext(ctx, instance); err != nil {
		return err
	}

	w.Typed = instance
	return nil
}

func (w *Versioned) MarshalYAML(ctx context.Context) ([]byte, error) {
	return yaml.Marshal(w.Typed)
}
