package main

import (
	"fmt"
)

type VisitedVertex struct {
	Name   string
	Path   string
	Weight int32
	Index  int
}

func (v *VisitedVertex) GetValue() int32 {
	return v.Weight
}

func (v *VisitedVertex) SetIndex(index int) {
	v.Index = index
}

func NewGraph() *Graph {
	return &Graph{
		Vertexes: NewMap[*Map[int32]](10000),
		Visited:  NewMap[*VisitedVertex](100),
		Heap:     &HeapMin{tree: make([]*VisitedVertex, 0)},
	}
}

type Graph struct {
	Vertexes *Map[*Map[int32]]
	Visited  *Map[*VisitedVertex]
	Heap     *HeapMin
}

func (g *Graph) AddVertex(vertexNameA, vertexNameB string, weight int32) {
	if _, ok := g.Vertexes.Get(vertexNameA); !ok {
		g.Vertexes.Set(vertexNameA, NewMap[int32](0))
	}
	vertexA, _ := g.Vertexes.Get(vertexNameA)
	vertexA.Set(vertexNameB, weight)

	if _, ok := g.Vertexes.Get(vertexNameB); !ok {
		g.Vertexes.Set(vertexNameB, NewMap[int32](0))
	}
	vertexB, _ := g.Vertexes.Get(vertexNameB)
	vertexB.Set(vertexNameA, weight)
}

func (g *Graph) Calculate(from string) (weight int32, path string, err error) {
	if _, ok := g.Vertexes.Get(from); !ok {
		return 0, "", fmt.Errorf("%s does not exist in Graph", string(from))
	}

	queue := make(chan string, 1)
	defer close(queue)
	queue <- from

	for currentVertexName := range queue {
		if _, ok := g.Visited.Get(currentVertexName); !ok {
			g.Visited.Set(currentVertexName, &VisitedVertex{Path: currentVertexName, Name: currentVertexName})
		}

		neighbors, _ := g.Vertexes.Get(currentVertexName)
		currentVertex, _ := g.Visited.Get(currentVertexName)
		for neighborVertexName, neighborWeight := range neighbors.Iterator() {
			destinationWeight := currentVertex.Weight + neighborWeight

			// if vertex has been visited and new weight is bigger than current weight then go to the next neighbor
			if visitedNeighbor, ok := g.Visited.Get(neighborVertexName); ok {
				if destinationWeight >= visitedNeighbor.Weight {
					continue
				}
				g.Heap.Delete(visitedNeighbor.Index)
			} else {
				g.Visited.Set(neighborVertexName, &VisitedVertex{Name: neighborVertexName})
			}

			visitedNeighbor, _ := g.Visited.Get(neighborVertexName)
			visitedNeighbor.Weight = destinationWeight
			visitedNeighbor.Path = currentVertex.Path + neighborVertexName
			g.Heap.Add(visitedNeighbor)
		}

		newVertexSource := g.Heap.GetRoot()
		if newVertexSource == nil {
			return currentVertex.Weight, currentVertex.Path, nil
		} else {
			queue <- newVertexSource.Name
		}
	}

	return 0, "", nil
}

func main() {
	fmt.Println("run tests with the following command: go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench .")
	fmt.Println("go tool pprof -http=:8084 mem.prof")
	fmt.Println("go tool pprof -http=:8084 cpu.prof")
}
