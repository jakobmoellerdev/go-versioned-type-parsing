package generic

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

func TestAccessSpec(t *testing.T) {
	reg := registry.NewTypeRegistry()
	typ := types.New("customAccess", "v1alpha1")
	type CustomAccess struct {
		types.VersionedType `json:",inline" yaml:",inline"`
		CustomField         string `json:"customField" yaml:"customField"`
	}
	registry.Register[CustomAccess](reg, typ)

	// Example YAML input
	yamlInput := `
type: customAccess
customField: exampleValue
`

	ctx := registry.Inject(context.Background(), reg)

	// Unmarshal into Versioned backed by registry.
	versioned := &Versioned{}
	if err := yaml.NewDecoder(strings.NewReader(yamlInput)).DecodeContext(ctx, versioned); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Assert the type and retrieve the specific implementation.
	customAccess, ok := versioned.Typed.(*CustomAccess)
	if !ok {
		t.Fatalf("type assertion failed: expected *CustomAccess, got %T", versioned.Typed)
	}

	// Verify the unmarshaled data.
	fmt.Printf("Custom Access: %+v\n", customAccess)
	if customAccess.CustomField != "exampleValue" {
		t.Errorf("unexpected CustomField value: got %q, want %q", customAccess.CustomField, "exampleValue")
	}
}
