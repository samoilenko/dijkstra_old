package main

import (
	"fmt"
	"runtime/debug"
)

type VisitedVertex struct {
	Name   []byte
	Path   []byte
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
		Vertexes: NewMap[*Map[int32]](),
		Visited:  NewMap[*VisitedVertex](),
		Heap:     &HeapMin{tree: make([]*VisitedVertex, 0)},
	}
}

type Graph struct {
	Vertexes *Map[*Map[int32]]
	Visited  *Map[*VisitedVertex]
	Heap     *HeapMin
}

func (g *Graph) AddVertex(vertexNameA, vertexNameB []byte, weight int32) {
	if _, ok := g.Vertexes.Get(vertexNameA); !ok {
		g.Vertexes.Set(vertexNameA, NewMap[int32]())
	}
	vertexA, _ := g.Vertexes.Get(vertexNameA)
	vertexA.Set(vertexNameB, weight)

	if _, ok := g.Vertexes.Get(vertexNameB); !ok {
		g.Vertexes.Set(vertexNameB, NewMap[int32]())
	}
	vertexB, _ := g.Vertexes.Get(vertexNameB)
	vertexB.Set(vertexNameA, weight)
}

func (g *Graph) Calculate(from []byte) (weight int32, path []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("holy cow: \n" + string(debug.Stack()))
		}
	}()

	if _, ok := g.Vertexes.Get(from); !ok {
		return 0, nil, fmt.Errorf("%s does not exist in Graph", string(from))
	}

	queue := make(chan []byte, 1)
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
			visitedNeighbor.Path = make([]byte, len(currentVertex.Path))
			copy(visitedNeighbor.Path, currentVertex.Path)
			visitedNeighbor.Path = append(visitedNeighbor.Path, neighborVertexName...)
			g.Heap.Add(visitedNeighbor)
		}

		newVertexSource := g.Heap.GetRoot()
		if newVertexSource == nil {
			return currentVertex.Weight, currentVertex.Path, nil
		} else {
			queue <- newVertexSource.Name
		}
	}

	return 0, nil, nil
}

func main() {
	fmt.Println("run tests with the following command: go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench .")
	fmt.Println("go tool pprof -http=:8084 mem.prof")
	fmt.Println("go tool pprof -http=:8084 cpu.prof")
}
