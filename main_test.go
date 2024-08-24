package main

import (
	"bytes"
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
	for i := 0; i < 100; i++ {
		for j := i + 1; j < 100; j++ {
			graph.AddVertex(
				[]byte{byte('A' + i)},
				[]byte{byte('A' + j)},
				int32(i+j),
			)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := graph.Calculate([]byte{'A'})
		if err != nil {
			b.Fatal(err)
		}
	}
}
