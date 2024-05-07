package sortkeys

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// i64

func NewI64KV(k int64, v interface{}) *I64KV {
	return &I64KV{
		K: k,
		V: v,
	}
}

type I64KV struct {
	K int64
	V interface{}
}

func (kv *I64KV) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *I64KV) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *I64KV) IntValue() int {
	return kv.V.(int)
}

func (kv *I64KV) StringValue() string {
	return kv.V.(string)
}

type I64KVSlice []*I64KV

func (a I64KVSlice) Len() int           { return len(a) }
func (a I64KVSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a I64KVSlice) Less(i, j int) bool { return a[i].K < a[j].K }

// u64

func NewU64KV(k uint64, v interface{}) *U64KV {
	return &U64KV{
		K: k,
		V: v,
	}
}

type U64KV struct {
	K uint64
	V interface{}
}

func (kv *U64KV) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *U64KV) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *U64KV) IntValue() int {
	return kv.V.(int)
}

func (kv *U64KV) StringValue() string {
	return kv.V.(string)
}

type U64KVSlice []*U64KV

func (a U64KVSlice) Len() int           { return len(a) }
func (a U64KVSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a U64KVSlice) Less(i, j int) bool { return a[i].K < a[j].K }

type U64K2V struct {
	K1 uint64
	K2 uint64
	V  interface{}
}

func (kv *U64K2V) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *U64K2V) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *U64K2V) IntValue() int {
	return kv.V.(int)
}

func (kv *U64K2V) StringValue() string {
	return kv.V.(string)
}

type U64K2VSlice []*U64K2V

func (a U64K2VSlice) Len() int      { return len(a) }
func (a U64K2VSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a U64K2VSlice) Less(i, j int) bool {
	if a[i].K1 == a[j].K1 {
		return a[i].K2 < a[j].K2
	}
	return a[i].K1 < a[j].K1
}

type U64K3V struct {
	K1 uint64
	K2 uint64
	K3 uint64
	V  interface{}
}

func (kv *U64K3V) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *U64K3V) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *U64K3V) IntValue() int {
	return kv.V.(int)
}

func (kv *U64K3V) StringValue() string {
	return kv.V.(string)
}

type U64K3VSlice []*U64K3V

func (a U64K3VSlice) Len() int      { return len(a) }
func (a U64K3VSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a U64K3VSlice) Less(i, j int) bool {
	ai, aj := a[i], a[j]
	if ai.K1 == aj.K1 {
		if ai.K2 == aj.K2 {
			return ai.K3 < aj.K3
		}
		return ai.K2 < aj.K2
	}
	return ai.K1 < aj.K1
}

// int

func NewIntKV(k int, v interface{}) *IntKV {
	return &IntKV{
		K: k,
		V: v,
	}
}

type IntKV struct {
	K int
	V interface{}
}

func (kv *IntKV) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *IntKV) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *IntKV) IntValue() int {
	return kv.V.(int)
}

func (kv *IntKV) StringValue() string {
	return kv.V.(string)
}

type IntKVSlice []*IntKV

func (a IntKVSlice) Len() int           { return len(a) }
func (a IntKVSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IntKVSlice) Less(i, j int) bool { return a[i].K < a[j].K }

// string

type StrKV struct {
	K string
	V interface{}
}

func (kv *StrKV) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *StrKV) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *StrKV) IntValue() int {
	return kv.V.(int)
}

func (kv *StrKV) StringValue() string {
	return kv.V.(string)
}

type StrKVSlice []*StrKV

func (a StrKVSlice) Len() int           { return len(a) }
func (a StrKVSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StrKVSlice) Less(i, j int) bool { return a[i].K < a[j].K }

// kv

type KV struct {
	K interface{}
	V interface{}
}

func (kv *KV) U64Key() uint64 {
	return kv.K.(uint64)
}

func (kv *KV) I64Key() int64 {
	return kv.K.(int64)
}

func (kv *KV) IntKey() int {
	return kv.K.(int)
}

func (kv *KV) StringKey() string {
	return kv.K.(string)
}

func (kv *KV) U64Value() uint64 {
	return kv.V.(uint64)
}

func (kv *KV) I64Value() int64 {
	return kv.V.(int64)
}

func (kv *KV) IntValue() int {
	return kv.V.(int)
}

func (kv *KV) StringValue() string {
	return kv.V.(string)
}

type Pair[K constraints.Ordered, V any] struct {
	K K
	V V
}

func SortPairs[K constraints.Ordered, V any](pairs []*Pair[K, V]) {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].K < pairs[j].K
	})
}

type PairArray[K constraints.Ordered, V any] struct {
	Array []*Pair[K, V]
}

func (a *PairArray[K, V]) Append(k K, v V) *PairArray[K, V] {
	a.Array = append(a.Array, &Pair[K, V]{k, v})
	return a
}

func (a *PairArray[K, V]) Sort() {
	sort.Slice(a.Array, func(i, j int) bool {
		return a.Array[i].K < a.Array[j].K
	})
}
