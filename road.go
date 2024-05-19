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

func printPathLevels(paths []Path, antCount int) {
	if len(paths) == 0 {
		fmt.Println("Başlangıç ve bitiş noktası arasında yol bulunamadı.")
		return
	}

	// Karıncaların takip edeceği yolları antPaths dizisine atar
	antPaths := make([][]string, antCount)
	for i := 0; i < antCount; i++ {
		antPaths[i] = paths[i%len(paths)].Steps
	}

	// En uzun yolun uzunluğunu hesapla
	maxPathLength := 0
	for _, path := range paths {
		if len(path.Steps) > maxPathLength {
			maxPathLength = len(path.Steps)
		}
	}

	// Karıncaların pozisyonlarını, adım sayılarını ve düğüm işgal durumlarını tutan diziler oluştur
	antPositions := make([]int, antCount)
	nodeOccupied := make(map[string]bool)
	antSteps := make([]int, antCount) // Karıncaların attığı adım sayısını takip eder

	// Başlangıçta tüm karıncaların pozisyonlarını ve adım sayılarını başlat
	for i := 0; i < antCount; i++ {
		antPositions[i] = 1
		antSteps[i] = 1
	}

	round := 1 // Mevcut turu takip eder

	// Başlangıç düğümünden çıkan bağlantı sayısını belirler
	startNodeConnections := len(paths)

	// Tüm karıncalar bitiş düğümüne ulaşana kadar simülasyonu döngüde tut
	for {
		allAntsFinished := true // Tüm karıncaların bitip bitmediğini kontrol eder
		roundOutput := []string{}

		// Her turda başlangıç düğümünden hareket eden karıncaların sayısını sınırla
		antsMovingFromStart := 0

		// Her karınca için hareketi kontrol eder
		for i := 0; i < antCount; i++ {
			// Eğer karınca bitiş düğümüne ulaştıysa devam et
			if antPositions[i] >= len(antPaths[i]) {
				continue // Eğer bu karınca bitiş düğümüne ulaşmışsa atla
			}

			// Karınca adımlarını en uzun yola göre kontrol et
			if antSteps[i] < maxPathLength {
				nextNode := antPaths[i][antPositions[i]] // Karıncanın gideceği bir sonraki düğüm

				// Bir önceki düğümün artık boş olduğunu belirt
				if antPositions[i] > 1 && antPositions[i]-1 < len(antPaths[i]) {
					nodeOccupied[antPaths[i][antPositions[i]-1]] = false
				}

				// Başlangıç düğümünden çıkan karıncaların sayısını sınırla
				if antPositions[i] == 1 {
					if antsMovingFromStart >= startNodeConnections {
						continue // Eğer sınır aşıldıysa bu karıncayı atla
					}
					antsMovingFromStart++
				}

				// Eğer düğüm işgal edilmemişse veya bitiş düğümüyse karıncayı hareket ettir
				if !nodeOccupied[nextNode] || nextNode == antPaths[i][len(antPaths[i])-1] {
					roundOutput = append(roundOutput, fmt.Sprintf("L%d-%s", i+1, nextNode))
					nodeOccupied[nextNode] = true
					antPositions[i]++ // Karıncayı bir sonraki düğüme taşı
					antSteps[i]++     // Karıncanın adım sayısını arttır
				}

				// Eğer karınca henüz yolunu tamamlamadıysa
				if antPositions[i] < len(antPaths[i]) {
					allAntsFinished = false // En az bir karınca hala hareket ediyor
				}
			} else {
				allAntsFinished = false // Bu karınca diğerlerini bekliyor
			}
		}

		// Her turun çıktısını yazdır
		if len(roundOutput) > 0 {
			fmt.Println(strings.Join(roundOutput, " "))
		}
		round++

		// Eğer tüm karıncalar bitiş düğümüne ulaştıysa döngüden çık
		if allAntsFinished {
			break // Tüm karıncalar yollarını tamamladı
		}
	}
}
