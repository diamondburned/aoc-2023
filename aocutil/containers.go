package aocutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"gopkg.in/typ.v4/sync2"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Set is a set of values.
type Set[T comparable] map[T]struct{}

// NewSet returns a new set.
func NewSet[T comparable](cap int) Set[T] {
	return make(Set[T], cap)
}

func NewSetFromSlice[T comparable](v []T) Set[T] {
	set := NewSet[T](len(v))
	for _, v := range v {
		set.Add(v)
	}
	return set
}

// Add adds the given value to the set.
func (s Set[T]) Add(v T) { s[v] = struct{}{} }

// Delete deletes the given value from the set.
func (s Set[T]) Delete(v T) { delete(s, v) }

// Has returns true if the set contains the given value.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Reset resets the set.
func (s *Set[T]) Reset() { *s = make(Set[T], len(*s)) }

// AnyMap is a map with any key type. Internally, keys are converted to strings
// using an opaque encoding. It is not possible to obtain the original key from
// the map. As a result, you cannot iterate over the keys of an AnyMap.
type AnyMap[K any, V any] struct {
	m            map[string]V
	encodeMapKey func(K) string
}

// NewAnyMap returns a new AnyMap.
func NewAnyMap[K any, V any]() AnyMap[K, V] {
	return AnyMap[K, V]{
		m:            map[string]V{},
		encodeMapKey: newKeyEncoder[K](),
	}
}

// Get returns the value for the given key.
func (m AnyMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[m.encodeMapKey(key)]
	return v, ok
}

// GetDefault returns the value for the given key, or the given default value if
// the key is not found.
func (m AnyMap[K, V]) GetDefault(key K, defaultValue V) V {
	v, ok := m.Get(key)
	if ok {
		return v
	}
	return defaultValue
}

// Getz returns the value for the given key or the zero-value if the key is not
// found.
func (m AnyMap[K, V]) Getz(key K) V {
	v, _ := m.Get(key)
	return v
}

// Has returns true if the given key exists in the map.
func (m AnyMap[K, V]) Has(key K) bool {
	_, ok := m.m[m.encodeMapKey(key)]
	return ok
}

// Set sets the given key-value pair into the map.
func (m AnyMap[K, V]) Set(key K, value V) {
	m.m[m.encodeMapKey(key)] = value
}

// Delete deletes the given key from the map.
func (m AnyMap[K, V]) Delete(key K) {
	delete(m.m, m.encodeMapKey(key))
}

// Reset resets the map.
func (m AnyMap[K, V]) Reset() {
	for k := range m.m {
		delete(m.m, k)
	}
}

var keyPool sync2.Pool[strings.Builder]
var keyTypeChecked sync2.Map[reflect.Type, struct{}]

func newKeyEncoder[K any]() func(K) string {
	str := new(strings.Builder)
	enc := cbor.NewEncoder(str)

	return func(key K) string {
		rtype := reflect.TypeOf(key)
		_, checked := keyTypeChecked.LoadOrStore(rtype, struct{}{})
		if !checked {
			rtype = rtypeElem(rtype)
			if rtype.Kind() == reflect.Struct {
				assertHashableStruct(rtype)
			}
		}

		if err := enc.Encode(key); err != nil {
			panic("cannot encode map key as CBOR")
		}

		s := str.String()
		str.Reset()

		return s
	}
}

func rtypeElem(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func assertHashableStruct(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			panic(fmt.Sprintf("field %s.%s is not exported", t, field.Name))
		}

		fieldType := rtypeElem(field.Type)
		if fieldType.Kind() == reflect.Struct {
			assertHashableStruct(fieldType)
		}
	}
}

// AnySet is similar to Set, except it's backed by an AnyMap which allows for
// any key type to be used as the map key. It is API-compatible with Set.
type AnySet[T any] AnyMap[T, struct{}]

// NewAnySet returns a new AnySet.
func NewAnySet[T any]() AnySet[T] {
	return AnySet[T](NewAnyMap[T, struct{}]())
}

// Add adds the given value to the set.
func (s AnySet[T]) Add(v T) { (AnyMap[T, struct{}])(s).Set(v, struct{}{}) }

// Delete deletes the given value from the set.
func (s AnySet[T]) Delete(v T) { (AnyMap[T, struct{}])(s).Delete(v) }

// Has returns true if the set contains the given value.
func (s AnySet[T]) Has(v T) bool { return (AnyMap[T, struct{}])(s).Has(v) }

// Reset resets the set.
func (s AnySet[T]) Reset() { (AnyMap[T, struct{}])(s).Reset() }

// AnyLRU is an LRU cache with any key type. Internally, keys are converted to
// strings using an opaque encoding. It is not possible to obtain the original
// key from the cache. As a result, you cannot iterate over the keys of an
// AnyLRU.
type AnyLRU[K, V any] struct {
	cache     *lru.Cache[string, V]
	encodeKey func(K) string
}

// NewAnyLRU creates a new AnyLRU instance.  If size is invalid, the function
// panics.
func NewAnyLRU[K, V any](size int) *AnyLRU[K, V] {
	cache, err := lru.New[string, V](size)
	if err != nil {
		panic(err)
	}
	return &AnyLRU[K, V]{
		cache:     cache,
		encodeKey: newKeyEncoder[K](),
	}
}

// Get returns the value for the given key.
func (lru *AnyLRU[K, V]) Get(key K) (V, bool) {
	return lru.cache.Get(lru.encodeKey(key))
}

// GetDefault returns the value for the given key, or the given default value if
// the key is not found.
func (lru *AnyLRU[K, V]) GetDefault(key K, defaultValue V) V {
	v, ok := lru.Get(key)
	if ok {
		return v
	}
	return defaultValue
}

// Getz returns the value for the given key or the zero-value if the key is not
// found.
func (lru *AnyLRU[K, V]) Getz(key K) V {
	v, _ := lru.Get(key)
	return v
}

// Has returns true if the given key exists in the map.
func (lru *AnyLRU[K, V]) Has(key K) bool {
	return lru.cache.Contains(lru.encodeKey(key))
}

// Set sets the given key-value pair into the map.
func (lru *AnyLRU[K, V]) Set(key K, value V) {
	lru.cache.Add(lru.encodeKey(key), value)
}

// Delete deletes the given key from the map.
func (lru *AnyLRU[K, V]) Delete(key K) {
	lru.cache.Remove(lru.encodeKey(key))
}
