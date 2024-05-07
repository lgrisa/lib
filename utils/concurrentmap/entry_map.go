package concurrentmap

import (
	"sync"
)

type mapKey interface {
	~int64 | ~int | ~int32 | ~uint64 | ~uint | ~uint32 | ~int16 | ~uint16
}

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (32) map shards.

type Map[K mapKey, O any] [32]*entryMapShard[K, O]

// A "thread" safe string to anything map.
type entryMapShard[K mapKey, O any] struct {
	items        map[K]O
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// New Creates a new concurrent map.
func New[K mapKey, O any]() *Map[K, O] {
	var m Map[K, O]
	for i := 0; i < 32; i++ {
		m[i] = &entryMapShard[K, O]{items: make(map[K]O)}
	}
	return &m
}

// Returns shard under given key
func (m *Map[K, O]) getShard(key K) *entryMapShard[K, O] {
	return m[key&31]
}

// Set Sets the given value under the specified key.
func (m *Map[K, O]) Set(key K, value O) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Upsert Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *Map[K, O]) Upsert(key K, value O, cb func(exist bool, valueInMap O, newValue O) O) (res O) {
	shard := m.getShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

// SetIfAbsent Sets the given value under the specified key if no value was associated with it.
// Returns old value if present. if ok, old is nil, new value is set. if !ok, old is the previous value, new value is not set.
func (m *Map[K, O]) SetIfAbsent(key K, value O) (O, bool) {
	// Get map shard.
	shard := m.getShard(key)
	shard.Lock()
	old, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return old, !ok
}

// Get Retrieves an element from map under given key.
func (m *Map[K, O]) Get(key K) (O, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Count Returns the number of elements within the map.
func (m *Map[K, O]) Count() int {
	result := 0

	for _, shard := range m {
		shard.RLock()
		result += len(shard.items)
		shard.RUnlock()
	}
	return result
}

// Has Looks up an item under specified key
func (m *Map[K, O]) Has(key K) bool {
	// Get shard
	shard := m.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove Removes an element from the map.
func (m *Map[K, O]) Remove(key K) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Pop Removes an element from the map and returns it
func (m *Map[K, O]) Pop(key K) (v O, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// IsEmpty Checks if map is empty.
func (m *Map[K, O]) IsEmpty() bool {
	return m.Count() == 0
}

// MapTupleTryMap Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type MapTupleTryMap[K mapKey, O any] struct {
	Key K
	Val O
}

// Iter Returns a buffered iterator which could be used in a for range loop.
func (m *Map[K, O]) Iter() chan []MapTupleTryMap[K, O] {
	ch := make(chan []MapTupleTryMap[K, O], 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *entryMapShard[K, O]) {
				// Foreach key, value pair.
				shard.RLock()
				tuples := make([]MapTupleTryMap[K, O], len(shard.items))
				idx := 0
				for key, val := range shard.items {
					tuples[idx] = MapTupleTryMap[K, O]{key, val}
					idx++
				}
				shard.RUnlock()
				ch <- tuples
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}

// Items Returns all items as map[string]interface{}
func (m *Map[K, O]) Items() map[K]O {
	tmp := make(map[K]O, 32)

	// Insert items to temporary map.
	for tp := range m.Iter() {
		for _, item := range tp {
			tmp[item.Key] = item.Val
		}
	}

	return tmp
}

// IterCb Callback based iterator, cheapest way to read
// all elements in a map.
func (m *Map[K, O]) IterCb(fn func(key K, v O)) {
	for tp := range m.Iter() {
		for _, item := range tp {
			fn(item.Key, item.Val)
		}
	}
}

// Keys Return all keys
func (m *Map[K, O]) Keys() chan []K {
	ch := make(chan []K, 32)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(32)
		// Foreach shard.
		for _, shard := range m {
			go func(shard *entryMapShard[K, O]) {
				// Foreach key, value pair.
				shard.RLock()
				array := make([]K, len(shard.items))
				idx := 0
				for key := range shard.items {
					array[idx] = key
					idx++
				}
				shard.RUnlock()
				ch <- array
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}
