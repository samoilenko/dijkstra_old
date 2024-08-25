package main

import (
	"bytes"
	"strconv"
	"testing"
)

func TestDijkstra(t *testing.T) {
	graph := NewGraph()

	graph.AddVertex([]byte{'D'}, []byte{'A'}, 4)
	graph.AddVertex([]byte{'D'}, []byte{'E'}, 2)
	graph.AddVertex([]byte{'A'}, []byte{'E'}, 4)
	graph.AddVertex([]byte{'A'}, []byte{'C'}, 5)
	// graph.AddVertex("A", "C", 4)
	graph.AddVertex([]byte{'E'}, []byte{'G'}, 5)
	// graph.AddVertex("E", "G", 1)
	graph.AddVertex([]byte{'E'}, []byte{'C'}, 4)
	graph.AddVertex([]byte{'C'}, []byte{'G'}, 5)
	graph.AddVertex([]byte{'C'}, []byte{'F'}, 5)
	graph.AddVertex([]byte{'C'}, []byte{'B'}, 2)
	graph.AddVertex([]byte{'G'}, []byte{'F'}, 5)
	graph.AddVertex([]byte{'B'}, []byte{'F'}, 2)

	// graph.AddVertex("G", "H", 5)
	// graph.AddVertex("G", "I", 4)
	// graph.AddVertex("I", "J", 2)
	weight, path, err := graph.Calculate([]byte{'D'})
	if err != nil {
		t.Error(err)
		return
	}

	const expectedWeight = 10
	if weight != expectedWeight {
		t.Errorf("expected weight: %d, actual weight: %d", expectedWeight, weight)
		return
	}

	expectedPath := []byte{'D', 'E', 'C', 'B', 'F'}
	if !bytes.Equal(expectedPath, path) {
		t.Errorf("expected path: %s, actual path: %s", string(expectedPath), string(path))
		return
	}
}

func BenchmarkDijkstra(b *testing.B) {
	graph := NewGraph()

	// Generate fixtures
	bufferA := make([]byte, 0, 10) // Preallocate buffer size based on expected length
	bufferB := make([]byte, 0, 10)

	for i := 0; i < 10000; i++ {
		bufferA = bufferA[:0] // Reset buffer without reallocating
		bufferA = append(bufferA, 'A')
		bufferA = strconv.AppendInt(bufferA, int64(i), 10)

		for j := i + 1; j < 10000; j++ {
			bufferB = bufferB[:0]
			bufferB = append(bufferB, 'A')
			bufferB = strconv.AppendInt(bufferB, int64(j), 10)

			graph.AddVertex(
				bufferA,
				bufferB,
				int32(i+j),
			)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := graph.Calculate([]byte("A0"))
		if err != nil {
			b.Fatal(err)
		}
	}
}
