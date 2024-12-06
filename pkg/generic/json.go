package generic

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/goccy/go-json"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

var _ interface {
	json.UnmarshalerContext
	json.MarshalerContext
} = &Versioned{}

func (w *Versioned) UnmarshalJSON(ctx context.Context, data []byte) error {
	// buffer the data so we can read it twice
	reader := io.Reader(bytes.NewReader(data))
	var buf bytes.Buffer
	reader = io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(reader)

	// Extract the type to decide the concrete implementation
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
	decoder = json.NewDecoder(reader)

	// Decode the YAML node directly into the specific Typed instance
	if err := decoder.DecodeContext(ctx, instance); err != nil {
		return err
	}

	w.Typed = instance
	return nil
}

func (w *Versioned) MarshalJSON(ctx context.Context) ([]byte, error) {
	return json.MarshalContext(ctx, w.Typed)
}
