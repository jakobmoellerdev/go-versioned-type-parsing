package types

import (
	"fmt"
	"strings"
)

type VersionedType struct {
	Type VersionedString `json:"type" yaml:"type"`
}

func (t *VersionedType) GetType() VersionedString {
	return t.Type
}

// Typed is any object that is defined by a type that is versioned.
type Typed interface {
	// GetType returns the objects type and version
	GetType() VersionedString
}

type VersionedString string

func (t VersionedString) String() string {
	return string(t)
}

func (t VersionedString) GetBase() string {
	return strings.Split(string(t), "/")[0]
}

func (t VersionedString) GetVersion() string {
	split := strings.Split(string(t), "/")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (t VersionedString) GetType() VersionedString {
	return t
}

func New(base, version string) Typed {
	return VersionedString(fmt.Sprintf("%s/%s", base, version))
}
