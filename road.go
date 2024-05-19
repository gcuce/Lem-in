package main

import (
	"fmt"
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

// FindAllPaths finds all paths between start and end nodes
func (g *Graph) FindAllPaths(start, end string) [][]string {
	var paths [][]string
	residualGraph := copyGraph(g.Capacity)
	visited := make(map[string]bool)
	dfs(residualGraph, g.Tunnels, start, end, visited, []string{}, &paths)
	return paths
}

// FindNonOverlappingPaths finds non-overlapping paths, preferring shorter ones when initial steps overlap
func (g *Graph) FindNonOverlappingPaths() []Path {
	allPaths := g.FindAllPaths(g.Start, g.End)
	roomUsed := make(map[string]bool)
	var selectedPaths []Path

	for _, path := range allPaths {
		conflict := false
		for _, room := range path {
			if room == g.Start || room == g.End {
				continue
			}
			if roomUsed[room] {
				conflict = true
				break
			}
		}
		if !conflict {
			selectedPaths = append(selectedPaths, Path{Steps: path})
			for _, room := range path {
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

func printPathLevels(paths []Path, antCount int) {
	if len(paths) == 0 {
		fmt.Println("Başlangıç ve bitiş noktası arasında yol bulunamadı.")
		return
	}

	// Assign paths for ants to follow to antPaths slice
	antPaths := make([][]string, antCount)
	for i := 0; i < antCount; i++ {
		antPaths[i] = paths[i%len(paths)].Steps
	}

	// Calculate the length of the longest path
	maxPathLength := 0
	for _, path := range paths {
		if len(path.Steps) > maxPathLength {
			maxPathLength = len(path.Steps)
		}
	}

	// Arrays to keep track of ant positions, step counts, and node occupancy
	antPositions := make([]int, antCount)
	nodeOccupied := make(map[string]bool)
	antSteps := make([]int, antCount)

	// Initialize ant positions and step counts at the beginning
	for i := 0; i < antCount; i++ {
		antPositions[i] = 1
		antSteps[i] = 1
	}

	round := 1 // Track the current round

	// Determine the number of connections leaving the start node
	startNodeConnections := len(paths)

	// Keep the simulation loop until all ants reach the end node
	for {
		allAntsFinished := true // Check if all ants have finished
		roundOutput := []string{}

		antsMovingFromStart := 0

		// Check the movement for each ant
		for i := 0; i < antCount; i++ {
			// Skip if ant reached the end node
			if antPositions[i] >= len(antPaths[i]) {
				continue
			}

			// Check ant movements according to the longest path
			if antSteps[i] < maxPathLength {
				nextNode := antPaths[i][antPositions[i]] // Next node for the ant

				// Mark the previous node as unoccupied
				if antPositions[i] > 1 && antPositions[i]-1 < len(antPaths[i]) {
					nodeOccupied[antPaths[i][antPositions[i]-1]] = false
				}

				// Limit ants moving from the start node in each round
				if antPositions[i] == 1 {
					if antsMovingFromStart >= startNodeConnections {
						continue // Skip this ant if the limit is reached
					}
					antsMovingFromStart++
				}

				// Move the ant if the node is unoccupied or it's the end node
				if !nodeOccupied[nextNode] || nextNode == antPaths[i][len(antPaths[i])-1] {
					roundOutput = append(roundOutput, fmt.Sprintf("L%d-%s", i+1, nextNode))
					nodeOccupied[nextNode] = true
					antPositions[i]++ // Move the ant to the next node
					antSteps[i]++     // Increment ant's step count
				}

				// If the ant hasn't completed its path yet
				if antPositions[i] < len(antPaths[i]) {
					allAntsFinished = false // At least one ant is still moving
				}
			} else {
				allAntsFinished = false // This ant is waiting for others
			}
		}

		// Print the output for each round
		if len(roundOutput) > 0 {
			fmt.Println(strings.Join(roundOutput, " "))
		}
		round++

		// Exit loop if all ants reached the end node
		if allAntsFinished {
			break // All ants have completed their paths
		}
	}
}
