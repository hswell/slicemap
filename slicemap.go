package slicemap

const (
	AlreadyExists = iota
	Successful
	NotExists
	RejectNilValue
)
const (
	MaxNoUseCount   = 1000
	MaxNoUsePercent = 0.1
)

type KV[K comparable] struct {
	Key   K
	Value any
}

func (kv *KV[K]) GetKey() K {
	return kv.Key
}

type SliceMap[K comparable] struct {
	DataSlice        []*KV[K]
	IndexMap         map[K]int
	maxIndex         int
	maxMaxNoUseCount int
	maxNoUsePercent  float64
}

func NewSliceMap[K comparable]() *SliceMap[K] {
	return &SliceMap[K]{
		DataSlice:        make([]*KV[K], 0, 1000),
		IndexMap:         make(map[K]int, 1000),
		maxMaxNoUseCount: MaxNoUseCount,
		maxNoUsePercent:  MaxNoUsePercent,
	}
}

func (sm *SliceMap[K]) Add(kv *KV[K]) (result int) {
	if kv == nil {
		return RejectNilValue
	}
	key := kv.Key
	if index, exists := sm.IndexMap[key]; exists {
		sm.DataSlice[index] = kv
		return AlreadyExists
	}
	sm.IndexMap[key] = sm.maxIndex
	if sm.maxIndex < len(sm.DataSlice) {
		sm.DataSlice[sm.maxIndex] = kv
	} else {
		sm.DataSlice = append(sm.DataSlice, kv)
	}

	sm.maxIndex++
	return Successful
}

func (sm *SliceMap[K]) Get(key K) (*KV[K], bool) {
	index, exists := sm.IndexMap[key]
	if !exists {
		return nil, false
	}
	value := sm.DataSlice[index]
	if value == nil {
		return nil, false
	}
	return value, true
}

func (sm *SliceMap[K]) Update(kv *KV[K]) (result int) {
	if kv == nil {
		return RejectNilValue
	}
	index, exists := sm.IndexMap[kv.GetKey()]
	if !exists {
		return NotExists
	}
	sm.DataSlice[index] = kv
	return Successful
}

func (sm *SliceMap[K]) Delete(key K) (result int) {
	if index, exists := sm.IndexMap[key]; exists {
		delete(sm.IndexMap, key)
		if index == sm.maxIndex-1 {
			sm.DataSlice[index] = nil
			sm.maxIndex--
		} else {
			sm.maxIndex--
			lastValue := sm.DataSlice[sm.maxIndex]
			sm.DataSlice[index] = lastValue
			if lastValue != nil {
				sm.IndexMap[lastValue.GetKey()] = index
			}
			sm.DataSlice[sm.maxIndex] = nil
		}
		sm.TryDeallocate()
		return Successful
	} else {
		return NotExists
	}
}

func (sm *SliceMap[K]) TryDeallocate() {
	noUse := len(sm.DataSlice) - sm.maxIndex
	percent := float64(noUse) / float64(len(sm.DataSlice))
	if noUse > 0 && noUse > sm.maxMaxNoUseCount && percent > sm.maxNoUsePercent {
		sm.DataSlice = sm.DataSlice[0:sm.maxIndex]
		// init all nil
		for i := range sm.DataSlice[sm.maxIndex:] {
			sm.DataSlice[sm.maxIndex+i] = nil
		}
	}
}

func (sm *SliceMap[K]) Range(f func(*KV[K])) {
	for _, kv := range sm.DataSlice {
		if kv != nil {
			f(kv)
		}
	}
}
