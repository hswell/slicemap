package test

import (
	"slicemap"
	"strconv"
	"testing"
)

// cpu: m1 pro
// memory: 32G

//goos: darwin
//goarch: arm64
//pkg: slicemap/test
//BenchmarkSliceMapValue_Add
//BenchmarkSliceMapValue_Add-10       	 3893378	       373.7 ns/op
//BenchmarkMapBased_Add
//BenchmarkMapBased_Add-10            	 4183965	       365.6 ns/op
//BenchmarkSliceMapValue_Get
//BenchmarkSliceMapValue_Get-10       	91218334	        11.57 ns/op
//BenchmarkMapBased_Get
//BenchmarkMapBased_Get-10            	107440154	        12.86 ns/op
//BenchmarkSliceMapValue_Delete
//BenchmarkSliceMapValue_Delete-10    	430066549	         2.755 ns/op
//BenchmarkMapBased_Delete
//BenchmarkMapBased_Delete-10         	504331117	         2.375 ns/op
//BenchmarkSliceMapValue_Range
//BenchmarkSliceMapValue_Range-10     	1000000	         0.0003383 ns/op
//BenchmarkMapBased_Range
//BenchmarkMapBased_Range-10          	1000000	         0.01718 ns/op
//遍历效率提升了一个数量级 map的遍历效率数据越大效率越低

// Create a map based implementation
type MapBased[K comparable] struct {
	data   map[K]*slicemap.KV[K]
	maxKey K
}

func (mb *MapBased[K]) Add(kv *slicemap.KV[K]) {
	if kv == nil {
		return
	}
	mb.data[kv.Key] = kv
}

func (mb *MapBased[K]) Get(key K) (*slicemap.KV[K], bool) {
	value, exists := mb.data[key]
	return value, exists
}

func (mb *MapBased[K]) Delete(key K) {
	delete(mb.data, key)
}

func BenchmarkSliceMapValue_Add(b *testing.B) {
	sm := slicemap.NewSliceMap[string]()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		k := strconv.Itoa(n)
		v := &slicemap.KV[string]{Key: k, Value: "value"}
		sm.Add(v)
	}
}

func BenchmarkMapBased_Add(b *testing.B) {
	mb := &MapBased[string]{data: map[string]*slicemap.KV[string]{}}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		k := strconv.Itoa(n)
		v := &slicemap.KV[string]{Key: k, Value: "value"}
		mb.Add(v)
	}
}

func BenchmarkSliceMapValue_Get(b *testing.B) {
	sm := slicemap.NewSliceMap[string]()
	keys := make([]string, 1000)
	for i := 0; i < len(keys); i++ {
		k := strconv.Itoa(i)
		keys[i] = k
		sm.Add(&slicemap.KV[string]{Key: k, Value: "value"})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := keys[n%len(keys)]
		_, _ = sm.Get(key)
	}
}

func BenchmarkMapBased_Get(b *testing.B) {
	mb := &MapBased[string]{data: map[string]*slicemap.KV[string]{}}
	keys := make([]string, 1000)
	for i := 0; i < len(keys); i++ {
		k := strconv.Itoa(i)
		keys[i] = k
		mb.Add(&slicemap.KV[string]{Key: k, Value: "value"})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := keys[n%len(keys)]
		_, _ = mb.Get(key)
	}
}

func BenchmarkSliceMapValue_Delete(b *testing.B) {
	sm := slicemap.NewSliceMap[string]()
	keys := make([]string, 1000)
	for i := 0; i < len(keys); i++ {
		k := strconv.Itoa(i)
		keys[i] = k
		sm.Add(&slicemap.KV[string]{Key: k, Value: "value"})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := keys[n%len(keys)]
		sm.Delete(key)
	}
}

func BenchmarkMapBased_Delete(b *testing.B) {
	mb := &MapBased[string]{data: map[string]*slicemap.KV[string]{}}
	keys := make([]string, 1000)
	for i := 0; i < len(keys); i++ {
		k := strconv.Itoa(i)
		keys[i] = k
		mb.Add(&slicemap.KV[string]{Key: k, Value: "value"})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		key := keys[n%len(keys)]
		mb.Delete(key)
	}
}

// in my macbook b.n =1000000000, run time out
func BenchmarkSliceMapValue_Range(b *testing.B) {
	sm := slicemap.NewSliceMap[string]()
	for n := 0; n < 100000; n++ {
		k := strconv.Itoa(n)
		v := &slicemap.KV[string]{Key: k, Value: "value"}
		sm.Add(v)
	}
	b.ResetTimer()
	sm.Range(func(s *slicemap.KV[string]) {
		_ = s.GetKey()
	})
}

func BenchmarkMapBased_Range(b *testing.B) {
	mb := &MapBased[string]{data: map[string]*slicemap.KV[string]{}}
	for n := 0; n < 100000; n++ {
		k := strconv.Itoa(n)
		v := &slicemap.KV[string]{Key: k, Value: "value"}
		mb.Add(v)
	}
	b.ResetTimer()
	for _, v := range mb.data {
		_ = v.GetKey()
	}
}
