# Alien Invasion
Warning: This repository is WIP so there might be some scope of improvement and bugs may occur.

### Problem Statement

Mad aliens are about to invade the earth!

```
You are given a map containing the names of cities in the non-existent world of
X. The map is in a file, with one city per line. The city name is first,
followed by 1-4 directions (north, south, east, or west). Each one represents a
road to another city that lies in that direction.

Given `n` aliens, these aliens start out at random places on the map, and wander around randomly,
following links. Each iteration, the aliens can travel in any of the directions
leading out of a city.

When two aliens end up in the same place, they fight, and in the process kill
each other and destroy the city. When a city is destroyed, it is removed from
the map, and so are any roads that lead into or out of it.
```

### Solution Approach

The fictious world is treated as a graph where each city is connected as nodes of the graph. The aliens have their individual go routines running and we have mapping between cities and alien positioned inside it. The aliens land sequentially in the order of their index (total aliens = N) and corresponding go routine starts which allocates random cities to them. Whenever two aliens land in same city the entry of city in graph is deleted along with entry of aliens and go routines of those two aliens are stopped. We also pause the go routines of other aliens when any particular alien makes move, and resume thereafter.

### Requirements

* [Go 1.18+](https://golang.org/dl/)

### Tests

```
$ go test ./*/ 
```

### Build

```
$ go build .
```

### Usage

```
$ ./Alien-invade --input=<INPUT_FILE> --N=<NUMBER_OF_ALIENS>
```

### Assumptions
- There is at most one adjacent city connected per direction, eg: if Foo has Bar connected to North, then no other city could be directly connected to North of Foo.
- The Aliens move independent of each other.
- A valid map is provided with all the connected cities listed, eg: if Bar is listed as North of Foo, then Foo should be listed as South of Bar
