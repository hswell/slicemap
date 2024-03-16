package testcode

import (
	"slicemap"
	"testing"
)

func TestSliceMapValue_Add(t *testing.T) {
	sm := slicemap.NewSliceMap[string]()

	// Test case 1: Add a non-nil KV
	kv1 := &slicemap.KV[string]{Key: "key1", Value: "value1"}
	result1 := sm.Add(kv1)
	if result1 != slicemap.Successful {
		t.Errorf("Expected %d, but got %d", slicemap.Successful, result1)
	}

	// Test case 2: Add a nil KV
	result2 := sm.Add(nil)
	if result2 != slicemap.RejectNilValue {
		t.Errorf("Expected %d, but got %d", slicemap.RejectNilValue, result2)
	}

	// Test case 3: Add a KV with existing key
	kv3 := &slicemap.KV[string]{Key: "key1", Value: "value2"}
	result3 := sm.Add(kv3)
	if result3 != slicemap.AlreadyExists {
		t.Errorf("Expected %d, but got %d", slicemap.AlreadyExists, result3)
	}
}

func TestSliceMapValue_Get(t *testing.T) {
	sm := slicemap.NewSliceMap[string]()

	// Test case 1: Get a KV that exists
	kv1 := &slicemap.KV[string]{Key: "key1", Value: "value1"}
	sm.Add(kv1)
	value1, exists1 := sm.Get("key1")
	if !exists1 || value1.Key != "key1" || value1.Value != "value1" {
		t.Errorf("Expected (%s, %s), but got (%v, %v)", "key1", "value1", value1.Key, value1.Value)
	}

	// Test case 2: Get a KV that doesn't exist
	value2, exists2 := sm.Get("key2")
	if exists2 {
		t.Errorf("Expected (%v, %v), but got (%v, %v)", nil, false, value2, exists2)
	}
}

func TestSliceMapValue_Update(t *testing.T) {
	sm := slicemap.NewSliceMap[string]()

	// Test case 1: Update a KV that exists
	kv1 := &slicemap.KV[string]{Key: "key1", Value: "value1"}
	sm.Add(kv1)
	kv2 := &slicemap.KV[string]{Key: "key1", Value: "value2"}
	result1 := sm.Update(kv2)
	if result1 != slicemap.Successful {
		t.Errorf("Expected %d, but got %d", slicemap.Successful, result1)
	}

	// Test case 2: Update a KV that doesn't exist
	kv3 := &slicemap.KV[string]{Key: "key2", Value: "value2"}
	result2 := sm.Update(kv3)
	if result2 != slicemap.NotExists {
		t.Errorf("Expected %d, but got %d", slicemap.NotExists, result2)
	}

	// Test case 3: Update a nil KV
	result3 := sm.Update(nil)
	if result3 != slicemap.RejectNilValue {
		t.Errorf("Expected %d, but got %d", slicemap.RejectNilValue, result3)
	}
}

func TestSliceMapValue_Delete(t *testing.T) {
	sm := slicemap.NewSliceMap[string]()

	// Test case 1: Delete a KV that exists
	kv1 := &slicemap.KV[string]{Key: "key1", Value: "value1"}
	sm.Add(kv1)
	result1 := sm.Delete("key1")
	if result1 != slicemap.Successful {
		t.Errorf("Expected %d, but got %d", slicemap.Successful, result1)
	}

	// Test case 2: Delete a KV that doesn't exist
	result2 := sm.Delete("key2")
	if result2 != slicemap.NotExists {
		t.Errorf("Expected %d, but got %d", slicemap.NotExists, result2)
	}
}

func TestSliceMapValue_TryDeallocate(t *testing.T) {
	_ = &slicemap.SliceMap[string]{DataSlice: []*slicemap.KV[string]{nil, nil, nil}, IndexMap: map[string]int{"key1": 0, "key2": 1, "key3": 2}}
	//
	//// Test case 1: TryDeallocate with no unused elements
	//sm.TryDeallocate()
	//if len(sm.DataSlice) != 3 {
	//	t.Errorf("Expected %d, but got %d", 3, len(sm.DataSlice))
	//}
	//
	//// Test case 2: TryDeallocate with unused elements
	//sm.DataSlice[1] = &slicemap.KV[string]{Key: "key4", Value: "value4"}
	//sm.DataSlice[2] = nil
	//sm.TryDeallocate()
}

func TestSliceMapValue_Range(t *testing.T) {
	sm := &slicemap.SliceMap[string]{DataSlice: []*slicemap.KV[string]{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}, IndexMap: map[string]int{"key1": 0, "key2": 1}}

	// Test case: Range over non-nil SliceMaps
	var values []string
	sm.Range(func(kv *slicemap.KV[string]) {
		values = append(values, kv.Value.(string))
	})
	expectedValues := []string{"value1", "value2"}
	if len(values) != len(expectedValues) || values[0] != expectedValues[0] || values[1] != expectedValues[1] {
		t.Errorf("Expected [%v], but got [%v]", expectedValues, values)
	}
}
