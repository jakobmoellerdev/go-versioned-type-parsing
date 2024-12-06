package generic

import (
	"strings"
	"testing"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

func TestRegistryDecoder_Decode(t *testing.T) {
	reg := registry.NewTypeRegistry()
	typ := types.New("customAccess", "v1alpha1")
	type CustomAccess struct {
		types.VersionedType `json:",inline" yaml:",inline"`
		CustomField         string `json:"customField" yaml:"customField"`
	}
	registry.Register[CustomAccess](reg, typ)
	decoder := NewTypedYAMLDecoder(strings.NewReader(
		`type: customAccess
customField: exampleValue`), reg)

	versioned, err := decoder.Decode()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if versioned == nil {
		t.Fatalf("versioned is nil")
	}
	if versioned.GetType().GetBase() != "customAccess" {
		t.Fatalf("unexpected type: got %q, want %q", versioned.GetType().GetBase(), "customAccess")
	}
}
