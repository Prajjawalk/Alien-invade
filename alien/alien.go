package alien

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	citymap "github.com/Prajjawalk/Alien-invade/map"
)

type Alien struct {
	Index       int
	Totalmoves  int
	Currentcity string
	Ws          chan int
	Destroyed   bool
	mu          sync.Mutex
}

// Possible worker states.
const (
	Stopped = iota
	Paused
	Running
)

func (alien *Alien) AlienServiceWorker(worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList, alienlist *[]*Alien, wg *sync.WaitGroup, mutex *sync.RWMutex) {
	defer wg.Done()
	state := Paused // Begin in the paused state.
	for {
		select {
		case state = <-alien.Ws:
			switch state {
			case Stopped:
				alien.mu.Lock()
				defer alien.mu.Unlock()
				close(alien.Ws)
				alien.Destroyed = true
				return
			}
		default:
			runtime.Gosched()
			if state == Paused {
				break
			}

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
						if city != "" && city != alien.Currentcity {
							citiesToMove = append(citiesToMove, city)
						}
					}
					fmt.Printf("Alien: %v, cities to move %v, current city: %v\n", alien.Index, citiesToMove, alien.Currentcity)
					var landingindex int
					if len(citiesToMove) == 0 {
						return
					} else {
						landingindex = rand.Intn(len(citiesToMove))
					}
					landingcity := citiesToMove[landingindex]
					if len((*simulation)[landingcity]) < 1 {
						alien.mu.Lock()
						SetState(*alienlist, 1)
						if len((*simulation)[alien.Currentcity]) == 1 && (*simulation)[alien.Currentcity][0] == alien.Index {
							mutex.Lock()
							(*simulation)[alien.Currentcity] = make([]int, 0)
							mutex.Unlock()
						}
						if alien.Totalmoves > 10000 {
							alien.Destroyed = true
							for i, a := range *alienlist {
								if a.Index == alien.Index {
									*alienlist = append((*alienlist)[:i], (*alienlist)[i+1:]...)
									break
								}
							}
							close(alien.Ws)
							return

						}
						(*simulation)[landingcity] = append((*simulation)[landingcity], alien.Index)
						alien.Totalmoves += 1
						alien.Currentcity = landingcity
						fmt.Printf("Alien: %v, current city is %v\n", alien.Index, alien.Currentcity)
						SetState(*alienlist, 2)
						alien.mu.Unlock()
					} else {
						for i, a := range *alienlist {
							if a.Index == (*simulation)[landingcity][0] {
								*alienlist = append((*alienlist)[:i], (*alienlist)[i+1:]...)
								break
							}
						}
						for i, a := range *alienlist {
							if a.Index == alien.Index {
								*alienlist = append((*alienlist)[:i], (*alienlist)[i+1:]...)
								break
							}
						}
						alien.mu.Lock()
						SetState(*alienlist, 1)
						alien.DestroyCity(landingcity, worldmap, simulation, citylist)
						SetState(*alienlist, 2)
						close(alien.Ws)
						alien.Destroyed = true
						alien.mu.Unlock()
						return
					}
				}
			}
		}
	}

}

func (alien *Alien) DestroyCity(cityName string, worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList) {
	if len((*simulation)[cityName]) == 0 {
		return
	}
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
	fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", cityName, existingAlien, alien.Index)
}

func SetState(aliens []*Alien, state int) {
	defer func() {
		// recover from panic if occured due to closed channel
		_ = recover()
	}()
	for _, a := range aliens {
		if !a.Destroyed {
			select {
			case a.Ws <- state:
			case <-time.After(time.Millisecond):
			}
		}
	}
}
