package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guilhermemalfatti/communautowatcher"
	"golang.org/x/sync/errgroup"
)

func main() {
	group, groupCtx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		err := communautowatcher.StartWatcher(groupCtx, communautowatcher.WatcherOptions{
			Interval:        time.Minute * 5,
			Watcher:         &Watcher{},
			IsFetchStations: false,
			IsFetchFlexCars: true,
		})
		return err
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

func (w *Watcher) OnFlexCarAvailable(cars []communautowatcher.FlexCarAvailabilityResp) {
	// todo
}
