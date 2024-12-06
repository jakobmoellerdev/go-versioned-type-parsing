package generic

import (
	"bytes"
	"context"
	"io"

	"github.com/goccy/go-yaml"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
)

var _ interface {
	yaml.BytesUnmarshalerContext
	yaml.BytesMarshalerContext
} = &Versioned{}

func (w *Versioned) UnmarshalYAML(ctx context.Context, data []byte) error {
	reg, err := registry.FromContext(ctx)
	if err != nil {
		return err
	}
	w.Typed, err = NewTypedYAMLDecoder(bytes.NewReader(data), reg).Decode()
	return err
}

func (w *Versioned) MarshalYAML(ctx context.Context) ([]byte, error) {
	return yaml.MarshalContext(ctx, w.Typed)
}

func NewYAMLDecoder(reader io.Reader) Decoder {
	return yaml.NewDecoder(reader)
}

func NewTypedYAMLDecoder(data io.Reader, reg registry.Registry) TypedDecoder {
	return NewTypedDecoder(data, reg, NewYAMLDecoder)
}
