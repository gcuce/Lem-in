package main

// dfs function finds all paths between start and end nodes using depth-first search algorithm
func dfs(capacity map[string]map[string]int, tunnels map[string][]string, start, end string, visited map[string]bool, path []string, paths *[][]string) {
	// Mark the current node as visited
	visited[start] = true

	// Update the path
	path = append(path, start)

	// If the end node is reached, copy and append the found path to paths slice
	if start == end {
		newPath := make([]string, len(path))
		copy(newPath, path)
		*paths = append(*paths, newPath)
	} else {
		// If the end node is not reached yet, visit the neighboring nodes
		for _, neighbor := range tunnels[start] {
			if !visited[neighbor] && capacity[start][neighbor] > 0 {
				dfs(capacity, tunnels, neighbor, end, visited, path, paths)
			}
		}
	}

	// Unmark the current node as visited
	delete(visited, start)

	// Update the path: backtrack
	path = path[:len(path)-1]
}
