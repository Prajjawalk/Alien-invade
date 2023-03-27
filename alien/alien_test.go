package alien

import (
	"sync"
	"testing"

	citymap "github.com/Prajjawalk/Alien-invade/map"
)

func TestAlienServiceWorker(t *testing.T) {
	// check that too many service workers spawned does not lead to any panic or error
	cities := []string{"Foo north=Bar west=Baz south=Qu-ux", "Bar south=Foo west=Bee", "Baz east=Foo", "Qu-ux north=Foo", "Bee east=Bar"}
	citygraph, simulation, citylist := citymap.CreateCityGraph(cities)
	wg := new(sync.WaitGroup)
	alienList := make([](*Alien), 0)
	var mutex = &sync.RWMutex{}
	numAlien := 100
	for i := 0; i < int(numAlien); i++ {
		if len(citygraph) == 0 {
			// no more cities left for invasion
			break
		}
		wg.Add(1)
		if len(citygraph) == 0 {
			break
		}
		alien := &Alien{
			Index:       i,
			Totalmoves:  0,
			Currentcity: "",
			Ws:          make(chan int),
			Destroyed:   false,
		}
		alienList = append(alienList, alien)
		go alien.AlienServiceWorker(&citygraph, &simulation, &citylist, &alienList, wg, mutex)
	}
	wg.Wait()
}

func TestDestroyCity(t *testing.T) {
	cities := []string{"Foo north=Bar west=Baz south=Qu-ux", "Bar south=Foo west=Bee", "Baz east=Foo", "Qu-ux north=Foo", "Bee east=Bar"}
	citygraph, simulation, citylist := citymap.CreateCityGraph(cities)
	cityName := "Foo"
	var mutex = &sync.RWMutex{}
	existingAlienIndex := 5
	(simulation)[cityName] = append((simulation)[cityName], existingAlienIndex)
	alien := &Alien{
		Index:       0,
		Totalmoves:  0,
		Currentcity: "",
		Ws:          make(chan int),
		Destroyed:   false,
	}
	alien.DestroyCity(cityName, &citygraph, &simulation, &citylist, mutex)

	_, exists := citygraph[cityName]
	if exists {
		t.Errorf("unable to destroy city: city %v exists even after destroying", cityName)
	}
}
