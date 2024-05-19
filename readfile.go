package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Room struct represents a room in the ant farm
type Room struct {
	Name string
	X    int
	Y    int
}

// Tunnel struct represents a tunnel connecting two rooms
type Tunnel struct {
	From string
	To   string
}

// ReadFile function reads the input file and processes the rooms and tunnels
func ReadFile(filename string) (int, []Room, []Tunnel, Room, Room, error) {
	var antCount int
	var rooms []Room
	var tunnels []Tunnel
	var startRoom, endRoom Room

	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, nil, Room{}, Room{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			parts := strings.Fields(scanner.Text())
			startRoom = Room{
				Name: parts[0],
				X:    convertToInt(parts[1]),
				Y:    convertToInt(parts[2]),
			}
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			parts := strings.Fields(scanner.Text())
			endRoom = Room{
				Name: parts[0],
				X:    convertToInt(parts[1]),
				Y:    convertToInt(parts[2]),
			}
		} else if antCount == 0 {
			antCount = convertToInt(line)
		} else if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			tunnels = append(tunnels, Tunnel{From: parts[0], To: parts[1]})
		} else {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				room := Room{
					Name: parts[0],
					X:    convertToInt(parts[1]),
					Y:    convertToInt(parts[2]),
				}
				rooms = append(rooms, room)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, nil, nil, Room{}, Room{}, err
	}

	return antCount, rooms, tunnels, startRoom, endRoom, nil
}

// Helper function to convert string to int
func convertToInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
