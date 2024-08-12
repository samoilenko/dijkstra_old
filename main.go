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

type Graph struct {
	Vertexes map[string]map[string]int32
	Visited  map[string]*VisitedVertex
	Heap     *HeapMin
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

func (g *Graph) Calculate(from string) (weight int32, path string, err error) {
	if _, ok := g.Vertexes[from]; !ok {
		return 0, "", fmt.Errorf("%s does not exist in Graph", from)
	}

	queue := make(chan string, 1)
	defer close(queue)
	queue <- from

	for currentVertexName := range queue {
		if _, ok := g.Visited[currentVertexName]; !ok {
			g.Visited[currentVertexName] = &VisitedVertex{Path: currentVertexName, Name: currentVertexName}
		}

		for neighborVertexName, neighborWeight := range g.Vertexes[currentVertexName] {
			destinationWeight := g.Visited[currentVertexName].Weight + neighborWeight

			// if vertex has been visited and new weight is bigger than current weight then go to the next neighbor
			if _, ok := g.Visited[neighborVertexName]; ok && destinationWeight >= g.Visited[neighborVertexName].Weight {
				continue
			}

			if _, ok := g.Visited[neighborVertexName]; !ok {
				g.Visited[neighborVertexName] = &VisitedVertex{Name: neighborVertexName}
			} else {
				g.Heap.Delete(g.Visited[neighborVertexName].Index)
			}

			g.Visited[neighborVertexName].Weight = destinationWeight
			g.Visited[neighborVertexName].Path = g.Visited[currentVertexName].Path + neighborVertexName
			g.Heap.Add(g.Visited[neighborVertexName])
		}

		newVertexSource := g.Heap.GetRoot()
		if newVertexSource == nil {
			// close(queue)
			return g.Visited[currentVertexName].Weight, g.Visited[currentVertexName].Path, nil
		} else {
			queue <- newVertexSource.Name
		}
	}

	return 0, "", nil
}

func main() {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	fmt.Println("run tests with the following command: go test -race -v main.go main_test.go heapMin.go")
}
