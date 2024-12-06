package unstructured

import (
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

type Unstructured struct {
	Data map[string]any `json:"-" yaml:"-"`
}

var _ interface {
	json.Marshaler
	json.Unmarshaler
	yaml.BytesMarshaler
	yaml.BytesUnmarshaler
	types.Typed
} = &Unstructured{}

func New() *Unstructured {
	return &Unstructured{
		Data: make(map[string]any),
	}
}

func (u *Unstructured) SetType(v types.VersionedString) {
	u.Data["type"] = v
}

func (u *Unstructured) GetType() types.VersionedString {
	v, _ := Get[types.VersionedString](u, "type")
	return v
}

func Get[T any](u *Unstructured, key string) (T, bool) {
	v, ok := u.Data[key]
	if !ok {
		return *new(T), false
	}
	t, ok := v.(T)
	return t, ok
}

func (u *Unstructured) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Data)
}

func (u *Unstructured) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &u.Data)
}

func (u *Unstructured) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(u.Data)
}

func (u *Unstructured) UnmarshalYAML(data []byte) error {
	return yaml.Unmarshal(data, &u.Data)
}
