// Package graphlib provides a simple library for working with directed, weighted graphs.
package graphlib

import (
    "container/heap"
    "fmt"
    "math"
    "os"
)

// Node represents a node in the graph.
type Node struct {
    Key int
}


// Edge represents an edge in the graph.
type Edge struct {
    Destination *Node
    Weight float64
}

// Graph represents the graph with nodes and edges.
type Graph struct {
    nodes map[int]*Node
    edges map[int]map[int]*Edge
}

// NewGraph creates a new instance of a Graph.
func NewGraph() *Graph {
    return &Graph{
        nodes: make(map[int]*Node),
        edges: make(map[int]map[int]*Edge),
    }
}

// AddNode adds a node with the given key to the graph.
// If the node already exists, it does nothing.
func (g *Graph) AddNode(key int) {
    if _, exists := g.nodes[key]; !exists {
        g.nodes[key] = &Node{Key: key}
    }
}

// RemoveNode removes a node and its associated edges from the graph.
// If the node does not exist, it does nothing.
func (g *Graph) RemoveNode(key int) {
    if _, exists := g.nodes[key]; !exists {
        return
    }
    delete(g.edges, key)
    for srcKey := range g.edges {
        delete(g.edges[srcKey], key)
    }
    delete(g.nodes, key)
}

// AddEdge adds an edge from a source to a destination node with the specified weight.
// If the source or destination nodes do not exist, it does nothing.
func (g *Graph) AddEdge(sourceKey, destinationKey int, weight float64) {
    source, sourceExists := g.nodes[sourceKey]
    destination, destinationExists := g.nodes[destinationKey]

    if !sourceExists || !destinationExists {
        return
    }

    if g.edges[sourceKey] == nil {
        g.edges[sourceKey] = make(map[int]*Edge)
    }
    g.edges[sourceKey][destinationKey] = &Edge{Destination: destination, Weight: weight}
}

// RemoveEdge removes an edge from the source to the destination node.
// If the source or destination nodes or the edge do not exist, it does nothing.
func (g *Graph) RemoveEdge(sourceKey, destinationKey int) {
    if _, exists := g.nodes[sourceKey]; !exists {
        return
    }
    if _, exists := g.nodes[destinationKey]; !exists {
        return
    }
    delete(g.edges[sourceKey], destinationKey)
}

// ExportToDot exports the graph to a file in DOT format.
// The filename specifies the output file.
// Returns an error if the file cannot be created or written to.
func (g *Graph) ExportToDot(filename string) error {
    // ... DOT format export implementation ...
    // Include the code for exporting to DOT format here.
}

// DijkstraData stores data for Dijkstra's algorithm.
type DijkstraData struct {
    Distance    map[int]float64
    Predecessor map[int]int
}

// FindShortestPath finds the shortest path between two nodes using Dijkstra's algorithm.
// Returns the path as a slice of node keys, the total distance, and a boolean indicating success.
// If no path exists, the returned slice is nil and the distance is zero.
func (g *Graph) FindShortestPath(sourceKey, destinationKey int) ([]int, float64, bool) {
    if _, exists := g.Nodes[sourceKey]; !exists {
        return nil, 0, false  // Source node doesn't exist
    }

    if _, exists := g.Nodes[destinationKey]; !exists {
        return nil, 0, false  // Destination node doesn't exist
    }

    dijkstraData := DijkstraData{
        Distance: make(map[int]float64),
        Predecessor: make(map[int]int),
    }
    for key := range g.Nodes {
        dijkstraData.Distance[key] = math.Inf(1)
    }
    dijkstraData.Distance[sourceKey] = 0

    pq := make(PriorityQueue, 0)
    heap.Init(&pq)
    heap.Push(&pq, &Item{Value: sourceKey, Priority: 0})

    for pq.Len() > 0 {
        current := heap.Pop(&pq).(*Item)
        currentNodeKey := current.Value

        if currentNodeKey == destinationKey {
            break
        }

        for _, edge := range g.Edges[currentNodeKey] {
            alt := dijkstraData.Distance[currentNodeKey] + edge.Weight
            if alt < dijkstraData.Distance[edge.Destination.Key] {
                dijkstraData.Distance[edge.Destination.Key] = alt
                dijkstraData.Predecessor[edge.Destination.Key] = currentNodeKey
                heap.Push(&pq, &Item{Value: edge.Destination.Key, Priority: alt})
            }
        }
    }

    path, totalDistance := reconstructPath(dijkstraData, sourceKey, destinationKey)
    return path, totalDistance, true
}

// reconstructPath reconstructs the shortest path from Dijkstra's data.
func reconstructPath(data DijkstraData, sourceKey, destinationKey int) ([]int, float64) {
    path := []int{}
    totalDistance, exists := data.Distance[destinationKey]
    if !exists || math.IsInf(totalDistance, 1) {
        return nil, 0  // Path doesn't exist
    }

    for at := destinationKey; at != sourceKey; at = data.Predecessor[at] {
        path = append([]int{at}, path...)
    }
    path = append([]int{sourceKey}, path...)

    return path, totalDistance
}

type Item struct {
    Value int    // The node key
    Priority float64    // The priority of the item in the queue
    Index int    // The index of the item in the heap
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].Index = i
    pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.Index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil  // avoid memory leak
    item.Index = -1 // for safety
    *pq = old[0 : n-1]
    return item
}

// ExportToDot exports the graph to a DOT format file.
func (g *Graph) ExportToDot(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString("digraph G {\n")
    if err != nil {
        return err
    }

    // Write nodes
    for _, node := range g.Nodes {
        _, err := file.WriteString(fmt.Sprintf("    %d;\n", node.Key))
        if err != nil {
            return err
        }
    }

    // Write edges
    for sourceKey, edgesMap := range g.Edges {
        for destinationKey, edge := range edgesMap {
            _, err := file.WriteString(fmt.Sprintf("    %d -> %d [label=\"%.2f\"];\n", sourceKey, destinationKey, edge.Weight))
            if err != nil {
                return err
            }
        }
    }

    _, err = file.WriteString("}\n")
    return err
}
