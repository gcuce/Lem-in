package main



// bfs fonksiyonu, genişlik öncelikli arama algoritması ile başlangıç ve bitiş düğümleri arasındaki en kısa yolu bulur
func bfs(capacity map[string]map[string]int, tunnels map[string][]string, start, end string) []string {
	visited := make(map[string]bool)
	parent := make(map[string]string)
	queue := []string{start}
	visited[start] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range tunnels[u] {
			if !visited[v] && capacity[u][v] > 0 {
				queue = append(queue, v)
				visited[v] = true
				parent[v] = u
				if v == end {
					var path []string
					for v != start {
						path = append([]string{v}, path...)
						v = parent[v]
					}
					path = append([]string{start}, path...)
					return path
				}
			}
		}
	}
	return nil
}