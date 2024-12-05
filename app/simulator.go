package app

import (
	"math/rand"
	"sync"
	"time"
	"github.com/Yomero3500/parkingGo/domain"
)

type Simulator struct {
	parkingService *domain.ParkingManager
}

func NewSimulator(parkingService *domain.ParkingManager	) *Simulator {
	return &Simulator{
		parkingService: parkingService,
	}
}

func (s *Simulator) SimulateVehicle(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	vehicle := &domain.Vehicle{ID: id}

	var space int
	var success bool
	for {
		space, success = s.parkingService.HandleVehicleEntry(vehicle)
		if success {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	parkingTime := 1 + rand.Intn(2)
	time.Sleep(time.Duration(parkingTime) * time.Second)

	for {
		if s.parkingService.HandleVehicleExit(vehicle, space) {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
}
