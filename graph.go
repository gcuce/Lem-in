package main

import "fmt"

// Graph struct represents the graph structure
type Graph struct {
	Nodes    map[string]*Room
	Tunnels  map[string][]string
	Capacity map[string]map[string]int
	Start    string
	End      string
}

// NewGraph initializes a new Graph
func NewGraph(antCount int, rooms []Room, tunnels []Tunnel, startRoom Room, endRoom Room) *Graph {
	graph := &Graph{
		Nodes:    make(map[string]*Room),
		Tunnels:  make(map[string][]string),
		Capacity: make(map[string]map[string]int),
		Start:    startRoom.Name,
		End:      endRoom.Name,
	}

	// Add start and end rooms
	graph.Nodes[startRoom.Name] = &startRoom
	graph.Nodes[endRoom.Name] = &endRoom

	// Add other rooms
	for i := range rooms {
		room := rooms[i]
		graph.Nodes[room.Name] = &room
	}

	for _, tunnel := range tunnels {
		if _, ok := graph.Tunnels[tunnel.From]; !ok {
			graph.Tunnels[tunnel.From] = make([]string, 0)
		}
		graph.Tunnels[tunnel.From] = append(graph.Tunnels[tunnel.From], tunnel.To)

		if _, ok := graph.Tunnels[tunnel.To]; !ok {
			graph.Tunnels[tunnel.To] = make([]string, 0)
		}
		graph.Tunnels[tunnel.To] = append(graph.Tunnels[tunnel.To], tunnel.From)

		if _, ok := graph.Capacity[tunnel.From]; !ok {
			graph.Capacity[tunnel.From] = make(map[string]int)
		}
		graph.Capacity[tunnel.From][tunnel.To] = 1

		if _, ok := graph.Capacity[tunnel.To]; !ok {
			graph.Capacity[tunnel.To] = make(map[string]int)
		}
		graph.Capacity[tunnel.To][tunnel.From] = 1
	}

	return graph
}

// String returns a string representation of the Graph
func (g *Graph) String() string {
	result := "\nNodes:\n"
	for name, room := range g.Nodes {
		result += fmt.Sprintf("%s: %+v\n", name, *room)
	}
	result += "\nTunnels:\n"
	for from, toList := range g.Tunnels {
		for _, to := range toList {
			result += fmt.Sprintf("%s -> %s\n", from, to)
		}
	}
	return result
}
