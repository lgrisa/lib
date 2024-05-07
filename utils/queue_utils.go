package utils

import "github.com/eapache/queue"

func NewQueue() *Queue {
	return &Queue{
		Queue: queue.New(),
	}
}

type Queue struct {
	*queue.Queue
}

func (q *Queue) Range(f func(v interface{}) (toContinue bool)) {
	for i := 0; i < q.Length(); i++ {
		if !f(q.Get(i)) {
			break
		}
	}
}

func (q *Queue) RangeWithStartIndex(startIndex int, f func(v interface{}) (toContinue bool)) {
	for i := startIndex; i < q.Length(); i++ {
		if !f(q.Get(i)) {
			break
		}
	}
}

func (q *Queue) ReverseRange(f func(v interface{}) (toContinue bool)) {
	for i := 1; i <= q.Length(); i++ {
		if !f(q.Get(-i)) {
			break
		}
	}
}

func (q *Queue) ReverseRangeWithStartIndex(startIndex int, f func(v interface{}) (toContinue bool)) {
	for i := startIndex + 1; i <= q.Length(); i++ {
		if !f(q.Get(-i)) {
			break
		}
	}
}

// NewRingList creates a new RingList with the given capacity.

func NewRingList(capacity int) *RingList {
	return &RingList{
		Queue:    NewQueue(),
		capacity: capacity,
	}
}

type RingList struct {
	*Queue

	capacity int
}

func (r *RingList) Capacity() int {
	return r.capacity
}

func (r *RingList) SetCapacity(toSet int) {
	r.capacity = toSet
}

// Add puts an element on the end of the queue.
func (r *RingList) Add(elem interface{}) {
	if r.capacity <= 0 {
		return
	}

	if r.Length() >= r.capacity {
		r.Queue.Remove()
	}

	r.Queue.Add(elem)
}
