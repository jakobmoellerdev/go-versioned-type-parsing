package generic

import (
	"bytes"
	"context"
	"io"

	"github.com/goccy/go-json"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
)

var _ interface {
	json.UnmarshalerContext
	json.MarshalerContext
} = &Versioned{}

func (w *Versioned) UnmarshalJSON(ctx context.Context, data []byte) error {
	reg, err := registry.FromContext(ctx)
	if err != nil {
		return err
	}
	w.Typed, err = NewTypedJSONDecoder(bytes.NewReader(data), reg).Decode()
	return err
}

func (w *Versioned) MarshalJSON(ctx context.Context) ([]byte, error) {
	return json.MarshalContext(ctx, w.Typed)
}

func NewJSONDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func NewTypedJSONDecoder(data io.Reader, reg registry.Registry) TypedDecoder {
	return NewTypedDecoder(data, reg, NewJSONDecoder)
}
