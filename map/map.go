package citymap

import (
	"strings"
)

type Direction int

type WorldMap map[string][]string
type SimulationTrack map[string][]int
type CityList []string

const (
	North Direction = iota
	South
	East
	West
)

func CreateCityGraph(cities []string) (WorldMap, SimulationTrack, CityList) {
	worldmap := make(map[string][]string)
	simulationTrack := make(map[string][]int)
	citylist := make([]string, 0)
	for _, city := range cities {
		chunks := strings.Split(city, " ")
		cityName := chunks[0]
		links := make([]string, 4)
		for _, roads := range chunks[1:] {
			dir := strings.Split(roads, "=")
			switch dir[0] {
			case "north":
				links[North] = dir[1]
			case "south":
				links[South] = dir[1]
			case "east":
				links[East] = dir[1]
			case "west":
				links[West] = dir[1]
			}
		}
		citylist = append(citylist, cityName)
		worldmap[cityName] = links
		simulationTrack[cityName] = make([]int, 0)

	}
	return worldmap, simulationTrack, citylist
}
