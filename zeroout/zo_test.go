package zeroout

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
)

func TestZeroOutQC(t *testing.T) {
	algos := []func([]byte, int, int){
		ZeroOut,
		ZeroOutLN,
		ZeroOutNaive,
	}

	for _, F := range algos {
		name := runtime.FuncForPC(reflect.ValueOf(F).Pointer()).Name()
		name = strings.TrimPrefix(name, "github.com/dgraph-io/ristretto-benchmarks/zeroout.")
		t.Run(name, func(t *testing.T) {
			f := func(bs []byte, start, end int) bool {
				if start < 0 || start >= len(bs) {
					return true // bad generation
				}
				cs := make([]byte, len(bs))
				copy(cs, bs)
				F(bs, start, end)

				if end >= len(bs) {
					end = len(bs)
				}
				if end-start < 0 {
					return true // noop
				}
				for i := 0; i < len(cs); i++ {
					if i >= start && i <= end {
						if bs[i] != 0 {
							return false
						}
						continue
					}
					if bs[i] != cs[i] {
						return false
					}
				}
				return true
			}
			if err := quick.Check(f, nil); err != nil {
				t.Error(err)
			}
		})
	}

}

func BenchmarkZeroOut(b *testing.B) {
	for i := 8; i <= 1024; i *= 2 {
		bs := make([]byte, i)
		for i := range bs {
			bs[i] = 255
		}
		b.Run(fmt.Sprintf("ZeroOut_%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(bs); j += 8 {
					for k := j; k < len(bs); k += 8 {
						ZeroOut(bs, j, k)
					}
				}
			}
		})
		b.Run(fmt.Sprintf("ZeroOutLN_%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(bs); j += 8 {
					for k := j; k < len(bs); k += 8 {
						ZeroOutLN(bs, j, k)
					}
				}
			}
		})
		b.Run(fmt.Sprintf("ZeroOutNaive_%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(bs); j += 8 {
					for k := j; k < len(bs); k += 8 {
						ZeroOutNaive(bs, j, k)
					}
				}
			}
		})
	}
}
