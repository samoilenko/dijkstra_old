package main

import (
	"bytes"
	"iter"
)

const (
	defaultLen = 100
	offset64   = 14695981039346656037
	prime64    = 1099511628211
)

type item[T any] struct {
	key  []byte
	data T
}

type Map[T any] struct {
	size  int
	items []*item[T]
}

func (m *Map[T]) Set(rawKey []byte, data T) {
	hashIndex := m.generateHashIndex(rawKey)
	if hashIndex >= len(m.items)-1 {
		newLen := len(m.items)*2 + 1
		tmp := &Map[T]{
			size:  0,
			items: make([]*item[T], newLen),
		}

		for _, item := range m.items {
			if item != nil {
				tmp.Set(item.key, item.data)
			}
		}
		m.items = tmp.items
		hashIndex = m.generateHashIndex(rawKey)
	}

	key := make([]byte, len(rawKey))
	copy(key, rawKey)

	if m.items[hashIndex] == nil || m.items[hashIndex].key == nil {
		m.items[hashIndex] = &item[T]{
			key:  key,
			data: data,
		}
	} else if bytes.Equal(m.items[hashIndex].key, rawKey) {
		m.items[hashIndex].data = data
	} else {
		// find next free slot
		for i := hashIndex + 1; i < len(m.items)-1; i++ {
			if m.items[i] == nil {
				m.items[i] = &item[T]{
					key:  key,
					data: data,
				}
				break
			}
			if i == len(m.items)-1 {
				i = 0
			}
			if i == hashIndex {
				panic("Can't find a free slot")
			}
		}
	}
	m.size += 1
}

func (m *Map[T]) Get(rawKey []byte) (T, bool) {
	hashIndex := m.generateHashIndex(rawKey)
	if hashIndex >= len(m.items) || m.items[hashIndex] == nil {
		var zero T // zero value of T
		return zero, false
	}

	if bytes.Equal(m.items[hashIndex].key, rawKey) {
		return m.items[hashIndex].data, true
	}

	for i := hashIndex + 1; i < len(m.items)-1; i++ {
		if m.items[i] != nil && bytes.Equal(m.items[i].key, rawKey) {
			return m.items[i].data, true
		}
		if i == len(m.items)-1 {
			i = 0
		}
		if i == hashIndex {
			var zero T // zero value of T
			return zero, false
		}
	}

	var zero T // zero value of T
	return zero, false
}

func (m *Map[T]) Iterator() iter.Seq2[[]byte, T] {
	return func(yield func([]byte, T) bool) {
		for _, item := range m.items {
			if item != nil && !yield(item.key, item.data) {
				return
			}
		}
	}
}

func (m *Map[T]) generateHashIndex(hashData []byte) int {
	hash := uint64(offset64)
	for _, data := range hashData {
		hash ^= uint64(data) // FNV-1a is XOR then *
		hash *= prime64
	}

	// if m.size == 0 {
	// 	return int(hash & uint64(0))
	// }

	return int(hash & uint64(len(m.items)-1))
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		size:  0,
		items: make([]*item[T], 1),
	}
}
