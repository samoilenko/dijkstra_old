package main

import (
	"testing"
)

func TestDijkstra(t *testing.T) {
	graph := Graph{
		Vertexes: make(map[string]map[string]int32),
		Visited:  make(map[string]*VisitedVertex),
		Heap:     &HeapMin{tree: make([]*VisitedVertex, 0)},
		// Heap:     &HeapMin{tree: make([]ff, 2)},
	}

	graph.AddVertex("D", "A", 4)
	graph.AddVertex("D", "E", 2)
	graph.AddVertex("A", "E", 4)
	graph.AddVertex("A", "C", 5)
	// graph.AddVertex("A", "C", 4)
	graph.AddVertex("E", "G", 5)
	// graph.AddVertex("E", "G", 1)
	graph.AddVertex("E", "C", 4)
	graph.AddVertex("C", "G", 5)
	graph.AddVertex("C", "F", 5)
	graph.AddVertex("C", "B", 2)
	graph.AddVertex("G", "F", 5)
	graph.AddVertex("B", "F", 2)

	// graph.AddVertex("G", "H", 5)
	// graph.AddVertex("G", "I", 4)
	// graph.AddVertex("I", "J", 2)
	weight, path, err := graph.Calculate("D")
	if err != nil {
		t.Error(err)
		return
	}

	const expectedWeight = 10
	if weight != expectedWeight {
		t.Errorf("expected weight: %d, actual weight: %d", expectedWeight, weight)
		return
	}

	const expectedPath = "DECBF"
	if path != expectedPath {
		t.Errorf("expected path: %s, actual path: %s", expectedPath, path)
		return
	}
}

func BenchmarkDijkstra(b *testing.B) {
	// Initialize your graph here
	graph := &Graph{
		Vertexes: make(map[string]map[string]int32),
		Visited:  make(map[string]*VisitedVertex),
		Heap:     &HeapMin{},
	}

	// Generate fixtures
	for i := 0; i < 100; i++ {
		for j := i + 1; j < 100; j++ {
			graph.AddVertex(
				string(rune('A'+i)),
				string(rune('A'+j)),
				int32(i+j),
			)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := graph.Calculate("A")
		if err != nil {
			b.Fatal(err)
		}
	}
}
