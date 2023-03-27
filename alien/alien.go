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

// goroutine for the corresponding alien to move through the city graph
func (alien *Alien) AlienServiceWorker(worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList, alienlist *[]*Alien, wg *sync.WaitGroup, mutex *sync.RWMutex) {
	defer wg.Done()
	state := Running // Begin in the running state.
	for {
		select {
		case state = <-alien.Ws:
			switch state {
			case Stopped:
				alien.mu.Lock()
				close(alien.Ws)
				alien.Destroyed = true
				alien.mu.Unlock()
				return
			}
		default:
			runtime.Gosched()
			if state == Paused {
				break
			}

			for {
				// we pause the movement of other aliens before making any move
				SetState(*alienlist, 1, alien.Index)
				rand.Seed(time.Now().UnixNano())
				landingcity := ""

				// The city gets allocated randomly to the alien just landed on world.
				// If the alien's total moves > 0, then we fetch the list of connected cities to the current city and allocate city randomly from those connected cities
				if alien.Totalmoves == 0 && alien.Currentcity == "" {
					mutex.Lock()
					totalcities := len(*citylist)
					if totalcities == 0 {
						SetState(*alienlist, 2, alien.Index)
						mutex.Unlock()
						return
					}
					landingindex := rand.Intn(totalcities)
					landingcity = (*citylist)[landingindex]
					mutex.Unlock()
				} else {
					citiesToMove := make([]string, 0)
					for _, city := range (*worldmap)[alien.Currentcity] {
						if city != "" && city != alien.Currentcity {
							citiesToMove = append(citiesToMove, city)
						}
					}
					var landingindex int
					if len(citiesToMove) == 0 {
						SetState(*alienlist, 2, alien.Index)
						return
					} else {
						landingindex = rand.Intn(len(citiesToMove))
					}
					landingcity = citiesToMove[landingindex]
				}

				if len((*simulation)[landingcity]) < 1 {
					if len((*simulation)[alien.Currentcity]) == 1 && (*simulation)[alien.Currentcity][0] == alien.Index {
						(*simulation)[alien.Currentcity] = make([]int, 0)
					}
					if alien.Totalmoves > 10000 {
						// if more than 10000 moves are completed, stop the goroutine and destroy the alien
						alien.Destroyed = true
						for i, a := range *alienlist {
							if a.Index == alien.Index {
								*alienlist = append((*alienlist)[:i], (*alienlist)[i+1:]...)
								break
							}
						}
						close(alien.Ws)
						SetState(*alienlist, 2, alien.Index)
						return

					}
					(*simulation)[landingcity] = append((*simulation)[landingcity], alien.Index)
					alien.Totalmoves += 1
					alien.Currentcity = landingcity
				} else {
					// when two aliens land on same city, both of them are destroyed along with the city
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
					alien.DestroyCity(landingcity, worldmap, simulation, citylist, mutex)
					close(alien.Ws)
					alien.Destroyed = true
					SetState(*alienlist, 2, alien.Index)
					return
				}
				SetState(*alienlist, 2, alien.Index)
			}
		}
	}

}

func (alien *Alien) DestroyCity(cityName string, worldmap *citymap.WorldMap, simulation *citymap.SimulationTrack, citylist *citymap.CityList, mutex *sync.RWMutex) {
	// when the city gets destroyed, remove the entry of city from the graph and all the links conneted to it
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
	mutex.RLock()
	for i, v := range *citylist {
		if v == cityName {
			*citylist = append((*citylist)[:i], (*citylist)[i+1:]...)
			break
		}
	}
	mutex.RUnlock()
	fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", cityName, existingAlien, alien.Index)
}

// fan out process which signals other alien goroutines to stop/pause/resume
func SetState(aliens []*Alien, state int, sentBy int) {
	defer func() {
		// recover from panic if occured due to closed channel
		_ = recover()
	}()
	for _, a := range aliens {
		if a.Index == sentBy {
			continue
		}
		if !a.Destroyed {
			select {
			case a.Ws <- state:
			case <-time.After(time.Millisecond):
			}
		}
	}
}
