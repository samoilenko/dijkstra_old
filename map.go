package main

import (
	"iter"
)

const (
	defaultLen = 100
	offset64   = 14695981039346656037
	prime64    = 1099511628211
)

type item[T any] struct {
	key  string
	data T
}

type Map[T any] struct {
	size  int
	items []*item[T]
}

func (m *Map[T]) reHash() {
	newLen := len(m.items)*2 + 1
	items := m.items
	m.items = make([]*item[T], newLen)
	m.size = 0

	for _, item := range items {
		if item != nil {
			m.Set(item.key, item.data)
		}
	}
}

func (m *Map[T]) Set(rawKey string, data T) {
	hashIndex := m.generateHashIndex(rawKey)
	for hashIndex >= len(m.items)-1 || m.size == len(m.items) {
		m.reHash()
		hashIndex = m.generateHashIndex(rawKey)
	}

	key := rawKey
	if m.items[hashIndex] == nil || m.items[hashIndex].key == "" {
		m.items[hashIndex] = &item[T]{
			key:  key,
			data: data,
		}
	} else if m.items[hashIndex].key == rawKey {
		m.items[hashIndex].data = data
	} else {
		// find next free slot
		for i := hashIndex + 1; i < len(m.items); i++ {
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

func (m *Map[T]) Get(rawKey string) (T, bool) {
	hashIndex := m.generateHashIndex(rawKey)
	if hashIndex >= len(m.items) || m.items[hashIndex] == nil {
		var zero T // zero value of T
		return zero, false
	}

	if m.items[hashIndex].key == rawKey {
		return m.items[hashIndex].data, true
	}

	for i := hashIndex + 1; i < len(m.items); i++ {
		if m.items[i] != nil && len(m.items[i].key) == len(rawKey) && m.items[i].key == rawKey {
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

func (m *Map[T]) Iterator() iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for _, item := range m.items {
			if item != nil && !yield(item.key, item.data) {
				return
			}
		}
	}
}

func (m *Map[T]) generateHashIndex(hashData string) int {
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

func NewMap[T any](initLen int) *Map[T] {
	var mapLen int
	if initLen > 0 {
		mapLen = initLen
	} else {
		mapLen = 1
	}
	return &Map[T]{
		size:  0,
		items: make([]*item[T], mapLen),
	}
}
