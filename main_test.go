package main

import (
	"strconv"
	"testing"
)

func TestDijkstra(t *testing.T) {
	graph := NewGraph()

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

	expectedPath := "DECBF"
	if expectedPath != path {
		t.Errorf("expected path: %s, actual path: %s", string(expectedPath), string(path))
		return
	}
}

func FuzzCalculation(f *testing.F) {
	graph := NewGraph()

	graph.AddVertex("D", "A", 4)
	graph.AddVertex("D", "E", 2)
	graph.AddVertex("A", "E", 4)
	graph.AddVertex("A", "C", 4)
	graph.AddVertex("E", "G", 5)
	graph.AddVertex("E", "G", 1)
	graph.AddVertex("E", "C", 4)
	graph.AddVertex("C", "G", 5)
	graph.AddVertex("C", "F", 5)
	graph.AddVertex("C", "B", 2)
	graph.AddVertex("G", "F", 5)
	graph.AddVertex("B", "F", 2)

	// const expectedWeight = 10

	f.Add("D")
	f.Fuzz(func(t *testing.T, vertexName string) {
		graph.Calculate(vertexName)
	})
}

func FuzzAdd(f *testing.F) {
	graph := NewGraph()

	f.Add("D", "A", 4)
	f.Add("D", "E", 2)
	f.Add("A", "E", 4)
	f.Add("A", "C", 4)
	f.Add("E", "G", 5)
	f.Add("E", "G", 1)
	f.Add("E", "C", 4)
	f.Add("C", "G", 5)
	f.Add("C", "F", 5)
	f.Add("C", "B", 2)
	f.Add("G", "F", 5)
	f.Add("B", "F", 2)

	f.Fuzz(func(t *testing.T, vertexAName, vertexBName string, weight int) {
		graph.AddVertex(vertexAName, vertexBName, int32(weight))
	})
}

func BenchmarkDijkstra(b *testing.B) {
	graph := NewGraph()

	// Generate fixtures
	for i := 0; i < 10000; i++ {
		for j := i + 1; j < 10000; j++ {
			graph.AddVertex(
				"A"+strconv.Itoa(i),
				"A"+strconv.Itoa(j),
				int32(i+j),
			)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := graph.Calculate("A0")
		if err != nil {
			b.Fatal(err)
		}
	}
}
