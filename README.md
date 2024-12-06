# Versioned Type Parsing in Go

This is a project essay on versioned type parsing in golang with `yaml` and `json` as the input data format.

## Problem Statement

Given a `yaml` or `json` input data, we want to parse it into a versioned 
type. The versioned type is a struct that is only available at runtime.
Additionally, one should be able to dynamically register new types.
Any field combination should be parseable.

An example of the input data is as follows:

```yaml
version: 1
type: type1
data:
  field1: value1
  field2: value2
---
version: 2
type: type1
data:
  field3: value3
---
version: 1
type: type2
data:
  field4: value4
  field5: value5
```

## Approach

To achieve strong casteable typing in go, we need

1. A versioned type that is only available at runtime.
2. A way to dynamically register new types.
3. A way to parse any field combination into the correct versioned type
4. A way to parse `yaml` and `json` input data in 2 steps:
    1. Parse the input data into a generic type that allows to determine the 
       actual typing (e.g. just type and version)
    2. Parse the generic type into the versioned type after determining the 
       actual typing by looking into a constructor in the registry

There are 3 core problems with the standard library `encoding/json`, which 
makes it hard to implement this:

1. It is impossible to parse a generic type into a versioned type without
   using a custom unmarshaler.
2. It is impossible to pass a registry dynamically into the encoding, forcing 
   us to use global variables. We would rather want a dynamic context.
3. Context parsing is not bound to `context.Context`


## Solution

1. We will use `github.com/goccy/go-json` and `github.com/goccy/go-yaml` as 
   the underlying parsing libraries. They are faster and more flexible than 
   the standard library and allow us to pass a `context.Context` into the 
   unmarshalling process.
2. We will use a custom unmarshaler to parse the generic type into the 
   versioned type. We will need one custom unmarshal function for each 
   supported protocol / language. This custom function will not actually 
   determine which versioned type to use, but will only introspect the 
   correct type and then delegate the decision to the context.
3. We will use a context to pass a function that creates a correct type into 
   the unmarshaler. This allows us to dynamically register new types and to have a dynamic context.
4. We will use a `generic` wrapper to hook our custom functions into the 
   unmarshaler. This wrapper will be the only thing that is exposed to the 
   user. The user will not have to deal with the custom unmarshaler directly.

## Implementation

The implementation is split into 4 parts:

1. [`types/versioned.go`](pkg/types/versioned.go): The versioned type that 
   contains the only shared information between types and is used to 
   differentiate them.
2. [`registry/typed.go`](pkg/registry/typed.go): The registry that contains 
   the constructors for the versioned types and also binds them into the 
   `context.Context`
3. [`generic.go`](pkg/generic/generic.go): The generic wrapper that hooks the 
   custom 
   unmarshaler into the `encoding/json` and `encoding/yaml` libraries.
   - [`json.go`](pkg/generic/json.go): The custom unmarshaler for `json`
   - [`yaml.go`](pkg/generic/yaml.go): The custom unmarshaler for `yaml`
4. [`unstructured/unstructured.go`](pkg/unstructured/unstructured.go): The 
   unstructured type that is used to parse any unknown type if allowed by 
   the registry. Allows coding that is not bound to the versioned types.

## Usage & Example

See the following [example test](./README.md) on how to use the implementation:
```go
package generic

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/api/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/api/types"
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
```

For a more realistic example see this [folder](pkg/example).