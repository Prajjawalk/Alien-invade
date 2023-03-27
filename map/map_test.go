package citymap

import "testing"

func TestCreateCityGraph(t *testing.T) {
	cities := []string{"Foo north=Bar west=Baz south=Qu-ux", "Bar south=Foo west=Bee", "Baz east=Foo", "Qu-ux north=Foo", "Bee east=Bar"}
	citygraph, _, citylist := CreateCityGraph(cities)

	expected_citygraph_length := 5
	expected_total_cities := 5
	actual_citygraph_length := len(citygraph)
	actual_total_cities := len(citylist)

	if len(citygraph) != 5 {
		t.Errorf("error while creating city graph: length of city graph does match, expected %v, got %v", expected_citygraph_length, actual_citygraph_length)
	}

	if len(citylist) != 5 {
		t.Errorf("error while creating city graph: total number of cities does match, expected %v, got %v", expected_total_cities, actual_total_cities)
	}
}
