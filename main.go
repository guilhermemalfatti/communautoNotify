package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/guilhermemalfatti/communautowatcher"
	"github.com/umahmood/haversine"
	"golang.org/x/sync/errgroup"
)

func main() {
	group, groupCtx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		communautowatcher.StartWatcher(groupCtx, communautowatcher.WatcherOptions{
			Interval:              time.Minute * 5,
			Watcher:               &Watcher{},
			IsEnableFetchStations: false,
			IsEnableFetchFlexCars: true,
		})
		return fmt.Errorf("Watcher was interrupted.")
	})

	if err := group.Wait(); err != nil {
		fmt.Printf("failed to read from communauto API: %s", err)
	}
}

type Watcher struct{}

func (w *Watcher) GetQueries() []communautowatcher.CarQuery {
	// Could fetch your "queries" from a database or a config file
	startDate, _ := time.Parse("2006-01-02T15:04", "2022-10-11T11:00")
	endDate, _ := time.Parse("2006-01-02T15:04", "2022-10-11T11:30")

	return []communautowatcher.CarQuery{
		{
			StartDate:     startDate,
			EndDate:       endDate,
			FromLatitude:  "45.5393407",
			FromLongitude: "-73.6307189",
			CityID:        string(communautowatcher.Montreal),
		},
	}
}

func (w *Watcher) GetFlexCarQuery() communautowatcher.CarQuery {

	return communautowatcher.CarQuery{
		BranchID:   "1",
		LanguageID: "1",
		CityID:     string(communautowatcher.Montreal),
	}
}

func (w *Watcher) OnCarAvailable(query communautowatcher.CarQuery, cars []communautowatcher.Car) {
	// todo
}

func (w *Watcher) OnFlexCarAvailable(cars []communautowatcher.Car) {
	mtlHomeCoord := haversine.Coord{Lat: 45.540615, Lon: -73.636537}
	fmt.Printf("From Lat: %f and Long: %f find bellow the closest cars\n", mtlHomeCoord.Lat, mtlHomeCoord.Lon)

	filteredcars := []communautowatcher.Car{}
	for _, car := range cars {
		carCoord := haversine.Coord{Lat: car.Latitude, Lon: car.Longitude}

		_, km := haversine.Distance(mtlHomeCoord, carCoord)

		// filter cars less then 1 km of distance
		if km < 1.0 {
			car.Distance = km
			filteredcars = append(filteredcars, car)
		}
	}

	// order by distance
	sort.SliceStable(filteredcars, func(i, j int) bool {
		return filteredcars[i].Distance < filteredcars[j].Distance
	})

	for _, car := range filteredcars {
		fmt.Printf("car No: %d Plate: %s Distance: %f  lat: %f long: %f\n", car.CarNo, car.CarPlate, car.Distance, car.Latitude, car.Longitude)
	}
}
