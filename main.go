package main

import (
	"fmt"
	"sync"

	"github.com/Prajjawalk/Alien-invade/alien"
	citymap "github.com/Prajjawalk/Alien-invade/map"
	"github.com/Prajjawalk/Alien-invade/utils"
)

func main() {
	cities, err := utils.ReadFile("./cityDetails.txt")
	if err != nil {
		fmt.Println(err)
	}
	wg := new(sync.WaitGroup)
	citygraph, simulation, citylist := citymap.CreateCityGraph(cities)
	alienList := make([](*alien.Alien), 0)
	var mutex = &sync.RWMutex{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		if len(citygraph) == 0 {
			break
		}
		alien := &alien.Alien{
			Index:       i,
			Totalmoves:  0,
			Currentcity: "",
			Ws:          make(chan int),
			Destroyed:   false,
		}
		alienList = append(alienList, alien)
		go alien.AlienServiceWorker(&citygraph, &simulation, &citylist, &alienList, wg, mutex)
	}
	alien.SetState(alienList, 2)
	wg.Wait()

	s := ""
	for city, links := range citygraph {
		s += fmt.Sprintf("%v ", city)
		for i, c := range links {
			if c == "" {
				continue
			}
			switch i {
			case int(citymap.North):
				s += fmt.Sprintf("north=%v ", c)
			case int(citymap.South):
				s += fmt.Sprintf("south=%v ", c)
			case int(citymap.East):
				s += fmt.Sprintf("east=%v ", c)
			case int(citymap.West):
				s += fmt.Sprintf("west=%v ", c)
			}
		}
		s += "\n"
	}
	fmt.Print("output: \n" + s)
}
