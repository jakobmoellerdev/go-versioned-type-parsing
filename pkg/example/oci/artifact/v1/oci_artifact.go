package v1

import (
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/example"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/registry"
	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

// Register default example types in the registry.
func init() {
	registry.Register[OCIArtifactAccess](example.DefaultRegistry, OCIArtifactAccessType)
}

var OCIArtifactAccessType = types.New("ociArtifact", "v1")

// OCIArtifactAccess specifies example to an OCI artifact.
type OCIArtifactAccess struct {
	types.VersionedType `json:",inline" yaml:",inline"`
	ImageReference      string `json:"imageReference" yaml:"imageReference"` // OCI image reference
}
