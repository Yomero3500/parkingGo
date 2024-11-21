package domain

import (
	"fmt"
	"time"
)

type ParkingManager struct {
	lot *ParkingLot
}

func NewParkingManager(lot *ParkingLot) *ParkingManager {
	return &ParkingManager{lot: lot}
}

func (pm *ParkingManager) LocateEmptySlot() int {
	pm.lot.Mu.Lock()
	defer pm.lot.Mu.Unlock()

	for i := range pm.lot.Spaces {
		if !pm.lot.Spaces[i] {
			return i
		}
	}
	return -1
}

func (pm *ParkingManager) ReserveSlot(slotIndex int) {
	pm.lot.Mu.Lock()
	pm.lot.Spaces[slotIndex] = true
	pm.lot.Mu.Unlock()
}

func (pm *ParkingManager) FreeSlot(slotIndex int) {
	pm.lot.Mu.Lock()
	pm.lot.Spaces[slotIndex] = false
	pm.lot.Mu.Unlock()
}

func (pm *ParkingManager) RefreshStatus(message string) {
	pm.lot.UpdateChan <- message
}

func (pm *ParkingManager) HandleVehicleEntry(vehicle *Vehicle) (int, bool) {
	select {
	case <-pm.lot.SpacesAvailable:
		select {
		case pm.lot.Entrance <- true:
			if pm.lot.Direction == 0 || pm.lot.Direction == 1 {
				pm.lot.Direction = 1
				slot := pm.LocateEmptySlot()
				if slot != -1 {
					pm.ReserveSlot(slot)
					pm.RefreshStatus(fmt.Sprintf("Vehicle %d entering slot %d", vehicle.ID, slot))
					time.Sleep(time.Millisecond * 500)
					<-pm.lot.Entrance
					pm.lot.Direction = 0
					return slot, true
				}
			}
			<-pm.lot.Entrance
		default:
			pm.lot.SpacesAvailable <- true
			if pm.lot.Direction == -1 {
				pm.RefreshStatus(fmt.Sprintf("Vehicle %d waiting to enter", vehicle.ID))
				return -1, false
			}
		}
	default:
		pm.RefreshStatus(fmt.Sprintf("Vehicle %d waiting for a slot", vehicle.ID))
		return -1, false
	}
	return -1, false
}

func (pm *ParkingManager) HandleVehicleExit(vehicle *Vehicle, slotIndex int) bool {
	select {
	case pm.lot.Entrance <- true:
		if pm.lot.Direction == 0 || pm.lot.Direction == -1 {
			pm.lot.Direction = -1
			pm.FreeSlot(slotIndex)
			pm.lot.SpacesAvailable <- true

			pm.lot.Mu.Lock()
			pm.lot.VehiclesExited++
			totalExits := pm.lot.VehiclesExited
			pm.lot.Mu.Unlock()

			message := fmt.Sprintf("Vehicle %d exited slot %d (Total exits: %d)", vehicle.ID, slotIndex, totalExits)
			pm.RefreshStatus(message)

			time.Sleep(time.Millisecond * 500)
			<-pm.lot.Entrance
			pm.lot.Direction = 0
			return true
		}
		<-pm.lot.Entrance
	default:
		if pm.lot.Direction == 1 {
			pm.RefreshStatus(fmt.Sprintf("Vehicle %d waiting to exit", vehicle.ID))
			return false
		}
	}
	return false
}
