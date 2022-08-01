// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

// ChunkQueueIterator is the iterator type of ChunkQueue. Iterating ChunkQueue
// by iterators is not thread-safe. Manipulating invalid iterators may incur
// panics. Don't use an iterator of an element that has already been dequeued.
// Instead, use with checks and in loop. E.g.
//
// for it := someQueue.First(); it.Valid(); it.Next() {
// 		... // operations cannot pop element
// }
// Note: Begin() and First() are interchangeable
// for it := someQueue.Begin(); it.Valid(); { 			// forwards
// 		it.Next()
// 		q.Dequeue() // can pop element
// }
// for it := someQueue.Last(); it.Valid(); it.Next() {	// backwards
// 		...
// }
// for it := someQueue.End(); it.Prev(); {				// backwards
//		...
// }
type ChunkQueueIterator[T any] struct {
	idxInChunk int
	chunk      *chunk[T]
}

// First returns the first valid iterator of the queue, which represents the
// first element (if exists)
func (q *ChunkQueue[T]) First() *ChunkQueueIterator[T] {
	return &ChunkQueueIterator[T]{
		chunk:      q.firstChunk(),
		idxInChunk: q.firstChunk().l,
	}
}

// Last returns the last valid iterator of the queue, which represents the
// last element (if exists)
func (q *ChunkQueue[T]) Last() *ChunkQueueIterator[T] {
	return &ChunkQueueIterator[T]{
		chunk:      q.lastChunk(),
		idxInChunk: q.lastChunk().r - 1,
	}
}

// Begin is an alias of First(), for convenient
func (q *ChunkQueue[T]) Begin() *ChunkQueueIterator[T] {
	return q.First()
}

// End creates a special iterator of the queue representing the end. End()
// iterator is not valid iterator since it's not in the queue. Calling Valid()
// always get false.
func (q *ChunkQueue[T]) End() *ChunkQueueIterator[T] {
	return &ChunkQueueIterator[T]{
		chunk:      q.lastChunk(),
		idxInChunk: q.chunkLength,
	}
}

// GetIterator returns an iterator of a given index and nil for invalid indices
func (q *ChunkQueue[T]) GetIterator(idx int) *ChunkQueueIterator[T] {
	if idx < 0 || idx >= q.size {
		return nil
	}
	idx += q.chunks[q.head].l
	return &ChunkQueueIterator[T]{
		chunk:      q.chunks[q.head+idx/q.chunkLength],
		idxInChunk: idx % q.chunkLength,
	}
}

// Valid indicates if the element of the iterator is in queue
func (it *ChunkQueueIterator[T]) Valid() bool {
	return it.chunk != nil && it.idxInChunk >= it.chunk.l && it.idxInChunk < it.chunk.r
}

// Value returns the element value of a valid iterator which is in queue.
// It's meaningless and may panic otherwise.
func (it *ChunkQueueIterator[T]) Value() T {
	return it.chunk.data[it.idxInChunk]
}

// Replace replaces the element of the valid iterator. It returns true on success
func (it *ChunkQueueIterator[T]) Replace(v T) bool {
	if it.Valid() {
		it.chunk.data[it.idxInChunk] = v
		return true
	}
	return false
}

// Index returns the index of a valid iterator, and -1 otherwise.
// Attention: The time complexity is O(N). Please avoid using this method
func (it *ChunkQueueIterator[T]) Index() int {
	if !it.Valid() {
		return -1
	}
	q := it.chunk.queue
	idx := 0
	for i := q.head; i < q.tail; i++ {
		if q.chunks[i] != it.chunk {
			idx += q.chunks[i].len()
		} else {
			idx += it.idxInChunk - it.chunk.l
			break
		}
	}
	return idx
}

// Next updates the current iterator to its next iterator. It returns true if
// the next iterator is still in queue, and false otherwise. Calling Next for
// an invalid iterator is meaningless, and using invalid iterators may panic.
// Using Next
func (it *ChunkQueueIterator[T]) Next() bool {
	if it.chunk == nil {
		return false
	}

	it.idxInChunk++
	if it.idxInChunk < it.chunk.r {
		return true
	}

	c, q := it.chunk, it.chunk.queue
	if it.idxInChunk == q.chunkLength && c.next != nil && !c.empty() {
		it.idxInChunk, it.chunk = 0, c.next
		return true
	}

	it.idxInChunk = q.chunkLength
	return false
}

// Prev updates the current to its previous iterator. It returns true if the
// next iterator is in queue, and false otherwise. The Prev of an end iterator
// points to the last element of the queue if the queue is not empty.
// The return boolean value is useful for backwards iteration. E.g.
// `for it := someQueue.Last(); it.Valid(); it.Next() {...}`
// `for it := someQueue.End(); it.Prev; {...} `
func (it *ChunkQueueIterator[T]) Prev() bool {
	if it.chunk == nil {
		return false
	}

	c := it.chunk
	if it.idxInChunk < c.l || it.idxInChunk >= c.r {
		// if the iterator is an end iterator and the queue is not empty,
		// then the iterator shall point to the last element.
		if it.idxInChunk == len(c.data) && !c.queue.Empty() {
			lastChunk := c.queue.lastChunk()
			it.chunk, it.idxInChunk = lastChunk, lastChunk.r-1
			return true
		}
		return false
	}

	it.idxInChunk--
	if it.idxInChunk >= it.chunk.l {
		return true
	}

	if c.prev != nil {
		it.chunk, it.idxInChunk = c.prev, c.prev.r-1
		return true
	}
	it.idxInChunk = -1
	return false
}
