package searchkeys

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

var algos = []func([]uint64, uint64) int16{
	Naive,
	Binary,
	Clever,
	Parallel,
	SSESearch,
	AVXSearch,
	ASMSearch,
	Search2,
}

func TestSearch(t *testing.T) {

	keys := make([]uint64, 512)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}

	for _, F := range algos {
		name := runtime.FuncForPC(reflect.ValueOf(F).Pointer()).Name()
		name = strings.TrimPrefix(name, "github.com/dgraph-io/ristretto-benchmarks/")
		t.Run(name, func(t *testing.T) {
			for i := 0; i < len(keys); i++ {
				idx := int(F(keys, uint64(i)))
				require.Equal(t, (i+1)/2, idx, "%v\n%v", i, keys)
			}
			require.Equal(t, 256, int(F(keys, math.MaxUint64>>1)))
			require.Equal(t, 256, int(F(keys, math.MaxInt64)))
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	for s := 2; s < 32768; s *= 2 {
		b.StopTimer()
		keys := make([]uint64, s)
		for i := 0; i < len(keys); i += 2 {
			keys[i] = uint64(i)
			keys[i+1] = 1
		}
		for _, F := range algos {
			name := runtime.FuncForPC(reflect.ValueOf(F).Pointer()).Name()
			name = strings.TrimPrefix(name, "github.com/dgraph-io/ristretto-benchmarks/searchkeys.")
			b.StartTimer()
			b.Run(fmt.Sprintf("%v\t%d", name, s), func(b *testing.B) {
				var idx int16
				for i := 0; i < b.N; i++ {
					for j := 0; j < s; j++ {
						idx = F(keys, uint64(j))
					}
				}
				_ = idx
			})
			b.StopTimer()
		}
	}
}

type kv struct {
	k, v uint64
}

type kvs []kv

func (l kvs) Len() int           { return len(l) }
func (l kvs) Less(i, j int) bool { return l[i].k < l[j].k }
func (l kvs) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func BenchmarkSearchUnsorted(b *testing.B) {
	for s := 2; s < 32768; s *= 2 {
		b.StopTimer()
		keys := make([]uint64, s)
		for i := 0; i < len(keys); i += 2 {
			keys[i] = uint64(i)
			keys[i+1] = 1
		}
		askv := (*(*kvs)(unsafe.Pointer(&keys)))[:s/2]
		rand.Shuffle(len(askv), askv.Swap)
		for _, F := range algos {
			name := runtime.FuncForPC(reflect.ValueOf(F).Pointer()).Name()
			name = strings.TrimPrefix(name, "github.com/dgraph-io/ristretto-benchmarks/searchkeys.")
			b.StartTimer()
			b.Run(fmt.Sprintf("%v\t%d", name, s), func(b *testing.B) {
				var idx int16
				for i := 0; i < b.N; i++ {
					for j := 0; j < s; j++ {
						idx = F(keys, uint64(j))
					}
				}
				_ = idx
			})
			b.StopTimer()
		}
	}
}
