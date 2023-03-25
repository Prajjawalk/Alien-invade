package alien

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	citymap "github.com/Prajjawalk/Alien-invade/map"
)

type Alien struct {
	Index       int
	Totalmoves  int
	Currentcity string
}

func (alien *Alien) AlienServiceWorker(worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		totalcities := len(*citylist)
		rand.Seed(time.Now().UnixNano())
		if alien.Totalmoves == 0 && alien.Currentcity == "" {
			landingindex := rand.Intn(totalcities)
			landingcity := (*citylist)[landingindex]
			if len((*simulation)[landingcity]) < 2 {
				(*simulation)[landingcity] = append((*simulation)[landingcity], alien.Index)
				alien.Totalmoves += 1
				alien.Currentcity = landingcity
				fmt.Printf("Alien: %v, current city is %v\n", alien.Index, alien.Currentcity)
			}
		} else {
			citiesToMove := make([]string, 0)
			for _, city := range (*worldmap)[alien.Currentcity] {
				if city != "" {
					citiesToMove = append(citiesToMove, city)
				}
			}
			fmt.Printf("Alien: %v, cities to move %v, current city: %v\n", alien.Index, citiesToMove, alien.Currentcity)
			var landingindex int
			if len(citiesToMove) == 0 {
				break
			} else if len(citiesToMove) == 1 {
				landingindex = 0
			} else {
				landingindex = rand.Intn(len(citiesToMove))
			}
			landingcity := citiesToMove[landingindex]
			if len((*simulation)[landingcity]) < 1 {
				(*simulation)[landingcity] = append((*simulation)[landingcity], alien.Index)
				alien.Totalmoves += 1
				alien.Currentcity = landingcity
				fmt.Printf("Alien: %v, current city is %v\n", alien.Index, alien.Currentcity)
			} else {
				alien.DestroyCity(landingcity, worldmap, simulation, citylist)
				break
			}
		}
	}
}

func (alien *Alien) DestroyCity(cityName string, worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList) {
	existingAlien := (*simulation)[cityName][0]
	connectedcities := (*worldmap)[cityName]
	for idx, city := range connectedcities {
		if city != "" {
			switch idx {
			case int(citymap.North):
				(*worldmap)[city][citymap.South] = ""
			case int(citymap.South):
				(*worldmap)[city][citymap.North] = ""
			case int(citymap.East):
				(*worldmap)[city][citymap.West] = ""
			case int(citymap.West):
				(*worldmap)[city][citymap.East] = ""

			}
		}
	}
	delete(*worldmap, cityName)
	delete(*simulation, cityName)
	for i, v := range *citylist {
		if v == cityName {
			*citylist = append((*citylist)[:i], (*citylist)[i+1:]...)
			break
		}
	}
	fmt.Printf("%v has been destroyed by alien %v and alien %v!", cityName, existingAlien, alien.Index)
}
