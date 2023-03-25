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
	wg.Add(3)
	citygraph, simulation, citylist := citymap.CreateCityGraph(cities)
	fmt.Print(citygraph, simulation, citylist)
	for i := 0; i < 3; i++ {
		alien := &alien.Alien{
			Index:       i,
			Totalmoves:  0,
			Currentcity: "",
		}
		go alien.AlienServiceWorker(&citygraph, &simulation, &citylist, wg)
	}
	wg.Wait()
}
