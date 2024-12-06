package registry

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/jakobmoellerdev/go-versioned-type-parsing/pkg/types"
)

var result func() types.Typed

func BenchmarkContextInjectionAndLoading(b *testing.B) {
	benchmarks := []struct {
		times uint64
	}{
		{1_0},
		{1_00},
		{1_000},
		{10_000},
		{100_000},
	}
	for _, bm := range benchmarks {
		ctx := context.Background()
		reg := NewTypeRegistry()
		for i := range bm.times {
			typ := types.New("customAccess"+strconv.Itoa(int(i)), "v"+strconv.Itoa(int(i)))
			type CustomAccess struct {
				types.VersionedType `json:",inline" yaml:",inline"`
				CustomField         string `json:"customField" yaml:"customField"`
			}
			Register[CustomAccess](reg, typ)
		}
		ctx = Inject(ctx, reg)

		b.Run(fmt.Sprintf("time context access to first of %v types", bm.times), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, ok := ConstructorFromContext(ctx, "customAccess"+strconv.Itoa(0))
				if !ok {
					b.Fatalf("constructor not found")
				}
			}
		})
		b.Run(fmt.Sprintf("time context access to last of %v types", bm.times), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var f func() types.Typed
				f, ok := ConstructorFromContext(ctx, "customAccess"+strconv.Itoa(int(bm.times-1)))
				if !ok {
					b.Fatalf("constructor not found")
				}
				// always store the result to a package level variable
				// so the compiler cannot eliminate the Benchmark itself.
				result = f
			}
		})
	}
}
