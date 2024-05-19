package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . file.txt")
		os.Exit(1)
	}

	filename := os.Args[1]
	antCount, rooms, tunnels, startRoom, endRoom, err := ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Check if start and end rooms are defined
	if (startRoom == Room{}) || (endRoom == Room{}) {
		fmt.Println("Error: Start or end room not defined.")
		os.Exit(1)
	}

	// Create the graph
	graph := NewGraph(antCount, rooms, tunnels, startRoom, endRoom)

	// Print the graph
	fmt.Println("Graph:", graph)

	// Find all paths
	allPaths := graph.FindAllPaths(startRoom.Name, endRoom.Name)
	fmt.Println("All Paths:")
	for _, path := range allPaths {
		fmt.Println(path)
	}

	// Find non-overlapping paths
	paths := graph.FindNonOverlappingPaths()
	fmt.Println("Non-overlapping paths:")
	for _, path := range paths {
		fmt.Println(path)
	}

	// Print the path levels
	printPathLevels(paths, antCount)
}
