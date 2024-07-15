package main

import (
	"errors"
	"fmt"
	"math"
)

type VisitedVertex struct {
	Weight    int32
	Path      string
	IsVisited bool
}

type Graph struct {
	Vertexes map[string]map[string]int32
	Visited  map[string]*VisitedVertex
}

func (g *Graph) AddVertex(vertexNameA, vertexNameB string, weight int32) {
	if _, ok := g.Vertexes[vertexNameA]; !ok {
		g.Vertexes[vertexNameA] = make(map[string]int32)
	}
	vertexA := g.Vertexes[vertexNameA]
	vertexA[vertexNameB] = weight
	g.Vertexes[vertexNameA] = vertexA

	if _, ok := g.Vertexes[vertexNameB]; !ok {
		g.Vertexes[vertexNameB] = make(map[string]int32)
	}
	vertexB := g.Vertexes[vertexNameB]
	vertexB[vertexNameA] = weight
	g.Vertexes[vertexNameB] = vertexB
}

func (g *Graph) getLowestVisitedVertexName() string {
	var vertexName string
	var min int32 = math.MaxInt32
	for visitedVertexName, visitedVertex := range g.Visited {
		if !visitedVertex.IsVisited && min > visitedVertex.Weight {
			vertexName = visitedVertexName
			min = visitedVertex.Weight
		}
	}

	return vertexName
}

func (g *Graph) Calculate(from string) (weight int32, path string, err error) {
	if _, ok := g.Vertexes[from]; !ok {
		return 0, "", errors.New(from + " does not exist in Graph")
	}

	if _, ok := g.Visited[from]; !ok {
		g.Visited[from] = &VisitedVertex{Path: from}
	}

	g.Visited[from].IsVisited = true

	for neighborVertex, neighborWeight := range g.Vertexes[from] {
		if _, ok := g.Visited[neighborVertex]; ok && g.Visited[from].Weight+neighborWeight >= g.Visited[neighborVertex].Weight {
			continue
		}

		if _, ok := g.Visited[neighborVertex]; !ok {
			g.Visited[neighborVertex] = &VisitedVertex{}
		}

		g.Visited[neighborVertex].Weight = g.Visited[from].Weight + neighborWeight
		g.Visited[neighborVertex].Path = g.Visited[from].Path + neighborVertex

		// if _, ok := g.Visited[neighborVertex]; !ok {
		// 	g.Visited[neighborVertex] = &VisitedVertex{
		// 		Weight: g.Visited[from].Weight + neighborWeight,
		// 		Path:   g.Visited[from].Path + neighborVertex,
		// 	}
		// }

		// if g.Visited[from].Weight+neighborWeight < g.Visited[neighborVertex].Weight {
		// 	g.Visited[neighborVertex].Weight = g.Visited[from].Weight + neighborWeight
		// 	g.Visited[neighborVertex].Path = g.Visited[from].Path + neighborVertex
		// }
	}

	newVertexSource := g.getLowestVisitedVertexName()

	if newVertexSource == "" {
		return g.Visited[from].Weight, g.Visited[from].Path, nil
	}

	return g.Calculate(newVertexSource)
}

func main() {
	graph := Graph{
		Vertexes: make(map[string]map[string]int32),
		Visited:  make(map[string]*VisitedVertex),
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
		fmt.Println(err)
		return
	}

	fmt.Println(weight, path)
}
