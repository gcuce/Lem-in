package main

import (
	"fmt"
	"math"
	"strings"
)

// Path struct represents a path taken by an ant
type Path struct {
	Steps []string
}

// copyGraph creates a deep copy of the capacity graph
func copyGraph(original map[string]map[string]int) map[string]map[string]int {
	copy := make(map[string]map[string]int)
	for u, neighbors := range original {
		copy[u] = make(map[string]int)
		for v, capacity := range neighbors {
			copy[u][v] = capacity
		}
	}
	return copy
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FindAllPaths finds all paths from start to end using a modified Edmonds-Karp approach
func (g *Graph) FindAllPaths(start, end string) []Path {
	var paths []Path
	residualGraph := copyGraph(g.Capacity)
	for {
		path := bfs(residualGraph, g.Tunnels, start, end)
		if path == nil {
			break
		}
		paths = append(paths, Path{Steps: path})
		flow := math.MaxInt32
		for i := 0; i < len(path)-1; i++ {
			u := path[i]
			v := path[i+1]
			flow = min(flow, residualGraph[u][v])
		}
		for i := 0; i < len(path)-1; i++ {
			u := path[i]
			v := path[i+1]
			residualGraph[u][v] -= flow
			residualGraph[v][u] += flow
		}
	}
	return paths
}

// FindNonOverlappingPaths finds non-overlapping paths, preferring shorter ones when initial steps overlap
func (g *Graph) FindNonOverlappingPaths() []Path {
	allPaths := g.FindAllPaths(g.Start, g.End)
	roomUsed := make(map[string]bool)
	var selectedPaths []Path

	for _, path := range allPaths {
		conflict := false
		for _, room := range path.Steps {
			if room == g.Start || room == g.End {
				continue
			}
			if roomUsed[room] {
				conflict = true
				break
			}
		}
		if !conflict {
			selectedPaths = append(selectedPaths, path)
			for _, room := range path.Steps {
				if room == g.Start || room == g.End {
					continue
				}
				roomUsed[room] = true
			}
		}
	}

	// Resolve conflicts in initial steps by preferring shorter paths
	for i := 0; i < len(selectedPaths)-1; i++ {
		for j := i + 1; j < len(selectedPaths); j++ {
			if selectedPaths[i].Steps[1] == selectedPaths[j].Steps[1] {
				if len(selectedPaths[i].Steps) > len(selectedPaths[j].Steps) {
					selectedPaths[i], selectedPaths[j] = selectedPaths[j], selectedPaths[i]
				}
			}
		}
	}

	return selectedPaths
}

// printPathLevels prints the levels of paths taken by the ants
func printPathLevels(paths []Path, antCount int) {
	maxPathLength := 0
	for _, path := range paths {
		if len(path.Steps) > maxPathLength {
			maxPathLength = len(path.Steps)
		}
	}

	antPositions := make([]string, antCount)
	for step := 0; step < maxPathLength+antCount-1; step++ {
		movement := []string{}
		for ant := 0; ant < antCount; ant++ {
			pathIdx := ant % len(paths)
			if step-ant >= 0 && step-ant < len(paths[pathIdx].Steps)-1 { // Skip the start room
				antPositions[ant] = paths[pathIdx].Steps[step-ant+1]
				movement = append(movement, fmt.Sprintf("L%d-%s", ant+1, antPositions[ant]))
			}
		}
		if len(movement) > 0 {
			fmt.Println(strings.Join(movement, " "))
		}
	}
}
