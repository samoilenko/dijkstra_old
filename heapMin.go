package main

type HeapMin struct {
	tree []*VisitedVertex
}

func (h *HeapMin) getParentIndex(nodeIndex int) int {
	return int(nodeIndex / 2)
}

func (h *HeapMin) Delete(nodeIndex int) {
	h.swap(nodeIndex, len(h.tree)-1)
	h.tree = h.tree[:len(h.tree)-1]
	h.Heapify(nodeIndex)
}

func (h *HeapMin) Add(vertex *VisitedVertex) {
	h.tree = append(h.tree, vertex)
	childIndex := len(h.tree) - 1
	h.tree[len(h.tree)-1].Index = childIndex
	if len(h.tree) == 1 {
		return
	}

	for {
		parentIndex := h.getParentIndex(childIndex)
		if parentIndex == childIndex {
			return
		}

		if h.tree[parentIndex].GetValue() <= h.tree[childIndex].GetValue() {
			return
		}
		h.swap(parentIndex, childIndex)
		childIndex = parentIndex
	}
}

func (h *HeapMin) swap(i, j int) {
	if i == j {
		return
	}

	h.tree[i], h.tree[j] = h.tree[j], h.tree[i]
	h.tree[j].SetIndex(j)
	h.tree[i].SetIndex(i)
}

func (h *HeapMin) Heapify(nodeIndex int) {
	heapSize := len(h.tree)
	// queue := make(chan int, 1)
	// queue <- nodeIndex
	// for parentIndex := range queue {
	parentIndex := nodeIndex
	for {
		leftChildIndex := 2*parentIndex + 1
		rightChildIndex := 2*parentIndex + 2

		// swap with the smallest child
		smallestIndex := parentIndex
		if heapSize > rightChildIndex && h.tree[smallestIndex].GetValue() > h.tree[rightChildIndex].GetValue() { // swap with right child
			smallestIndex = rightChildIndex
		}
		if heapSize > leftChildIndex && h.tree[smallestIndex].GetValue() > h.tree[leftChildIndex].GetValue() {
			smallestIndex = leftChildIndex
		}

		if smallestIndex != parentIndex {
			h.swap(parentIndex, smallestIndex)
			// queue <- smallestIndex
			parentIndex = smallestIndex
		} else {
			break
			// close(queue)
		}
	}
}

func (h *HeapMin) GetRoot() *VisitedVertex {
	if len(h.tree) == 0 {
		return nil
	}

	root := h.tree[0]
	h.Delete(0)

	return root
}
