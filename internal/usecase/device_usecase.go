package usecase

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	// "encoding/json"

	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/repository"
)

type DeviceUsecase struct {
	Repo             *repository.DeviceRepository
	TypeMapRepo      *repository.DeviceTypeMapRepository
	TypeRepo         *repository.DeviceTypeRepository
	BroadcastChannel chan gin.H
}

func NewDeviceUsecase(repo *repository.DeviceRepository, typeMapRepo *repository.DeviceTypeMapRepository, typeRepo *repository.DeviceTypeRepository) *DeviceUsecase {
	return &DeviceUsecase{
		Repo:             repo,
		TypeMapRepo:      typeMapRepo,
		TypeRepo:         typeRepo,
		BroadcastChannel: make(chan gin.H, 30), // Initialize buffered channel
	}
}

func (u *DeviceUsecase) GetAllDevices() ([]domain.Device, error) {
	return u.Repo.GetAllDevices()
}

func (u *DeviceUsecase) GetAllDevicesWithTypes() ([]domain.Device, error) {
	devices, err := u.Repo.GetAllDevices()
	if err != nil {
		return nil, err
	}
	for i := range devices {
		types, _ := u.TypeMapRepo.GetDeviceTypes(devices[i].ID)
		devices[i].Types = types
	}
	return devices, nil
}

func (u *DeviceUsecase) InsertDevice(device *domain.Device) error {
	return u.Repo.InsertDevice(device)
}

func (u *DeviceUsecase) InsertDeviceWithTypes(device *domain.Device, typeIDs []uint) error {
	if err := u.Repo.InsertDevice(device); err != nil {
		return err
	}
	return u.TypeMapRepo.AddDeviceTypes(device.ID, typeIDs)
}

func (u *DeviceUsecase) UpdateDevice(device *domain.Device) error {
	return u.Repo.UpdateDevice(device)
}

func (u *DeviceUsecase) UpdateDeviceWithTypes(device *domain.Device, typeIDs []uint) error {
	if err := u.Repo.UpdateDevice(device); err != nil {
		return err
	}
	return u.TypeMapRepo.UpdateDeviceTypes(device.ID, typeIDs)
}

func (u *DeviceUsecase) UpdateDeviceStatus() {
	devices, err := u.Repo.GetAllDevices()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	for _, device := range devices {
		// fmt.Println("Checking device:", device.Name)
		// Check device status using HTTP or ping
		isOnline := false

		if strings.HasPrefix(device.IP, "http") {
			resp, err := http.Get(device.IP)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				isOnline = true
			} else {
				log.Printf("HTTP request failed for URL %s: %v", device.IP, err)
			}
		} else {
			cmd := exec.Command("ping", "-n", "1", device.IP)
			if err := cmd.Run(); err == nil {
				isOnline = true
			} else {
				log.Printf("Ping failed for IP %s: %v", device.IP, err)
			}
		}

		// fmt.Println("Device:", device.Name, "isOnline:", isOnline)

		// Update status and broadcast changes
		oldStatus := device.Status
		device.Status = map[bool]string{true: "online", false: "offline"}[isOnline]
		// fmt.Println("Device:", device.Name, "oldStatus:", oldStatus, "newStatus:", device.Status)
		if oldStatus != device.Status {
			log := domain.Log{
				DeviceID:  device.ID,
				OldStatus: oldStatus,
				NewStatus: device.Status,
				Logtime:   time.Now(),
			}
			// fmt.Println("Logging status change: ", log)
			u.Repo.CreateLog(&log)
		}
		device.LastOnline = time.Now()
		// Only update status and lastonline, not other fields
		err := u.Repo.DB.Model(&domain.Device{}).Where("id = ?", device.ID).Updates(map[string]interface{}{
			"status":     device.Status,
			"lastonline": time.Now(),
		}).Error
		if err != nil {
			log.Printf("Error updating device %s: %v", device.Name, err)
		}
		// Trigger broadcast to clients after updating status
		// u.BroadcastDevices()
	}
}

// func (u *DeviceUsecase) BroadcastDevices() {
// 	fmt.Println("Broadcasting device statuses to clients...")
// 	devices, err := u.Repo.GetAllDevices()
// 	// fmt.Println("Fetched devices for broadcast:", devices)
// 	if err != nil {
// 		log.Printf("Error fetching devices for broadcast: %v", err)
// 		return
// 	}

// 	// Calculate statistics
// 	total := len(devices)
// 	online := 0
// 	for _, device := range devices {
// 		if device.Status == "online" {
// 			online++
// 		}
// 	}
// 	offline := total - online

// 	// Broadcast to all clients
// 	u.BroadcastChannel <- gin.H{
// 		"total":   total,
// 		"online":  online,
// 		"offline": offline,
// 		"devices": devices,
// 	}
// 	fmt.Println("Broadcast completed.")
// }

// Ensure BroadcastDevices runs in a separate goroutine to avoid blocking
func (u *DeviceUsecase) BroadcastDevices() {
	go u.BroadcastDevices()
}

// func (u *DeviceUsecase) BroadcastDevices() {
// 	fmt.Println("Broadcasting device statuses to clients...")
// 	devices, err := u.Repo.GetAllDevices()
// 	// fmt.Println("Fetched devices for broadcast:", devices)
// 	if err != nil {
// 		log.Printf("Error fetching devices for broadcast: %v", err)
// 		return
// 	}

// 	// Calculate statistics
// 	total := len(devices)
// 	online := 0
// 	for _, device := range devices {
// 		if device.Status == "online" {
// 			online++
// 		}
// 	}
// 	offline := total - online

// 	// Broadcast to all clients
// 	u.BroadcastChannel <- gin.H{
// 		"total":   total,
// 		"online":  online,
// 		"offline": offline,
// 		"devices": devices,
// 	}
// 	fmt.Println("Broadcast completed.")
// }

func (u *DeviceUsecase) GetDeviceByID(id uint) (*domain.Device, error) {
	return u.Repo.GetDeviceByID(id)
}

func (u *DeviceUsecase) GetDeviceByIDWithTypes(id uint) (*domain.Device, error) {
	device, err := u.Repo.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}
	types, _ := u.TypeMapRepo.GetDeviceTypes(device.ID)
	device.Types = types
	return device, nil
}

func (u *DeviceUsecase) DeleteDevice(id uint) error {
	return u.Repo.DeleteDevice(id)
}

func (u *DeviceUsecase) GetDevicesByType(typeID uint) ([]domain.Device, error) {
	var deviceIDs []uint
	var maps []domain.DeviceTypeMap
	if err := u.TypeMapRepo.DB.Where("type_id = ?", typeID).Find(&maps).Error; err != nil {
		return nil, err
	}
	for _, m := range maps {
		deviceIDs = append(deviceIDs, m.DeviceID)
	}
	var devices []domain.Device
	if err := u.Repo.DB.Where("id IN ?", deviceIDs).Find(&devices).Error; err != nil {
		return nil, err
	}
	for i := range devices {
		types, _ := u.TypeMapRepo.GetDeviceTypes(devices[i].ID)
		devices[i].Types = types
	}
	return devices, nil
}

func (u *DeviceUsecase) GetDevicesByTypeMulti(typeIDs []uint) ([]domain.Device, error) {
	var deviceIDs []uint
	var maps []domain.DeviceTypeMap
	if err := u.TypeMapRepo.DB.Where("type_id IN ?", typeIDs).Find(&maps).Error; err != nil {
		return nil, err
	}
	idsMap := make(map[uint]struct{})
	for _, m := range maps {
		idsMap[m.DeviceID] = struct{}{}
	}
	for id := range idsMap {
		deviceIDs = append(deviceIDs, id)
	}
	var devices []domain.Device
	if err := u.Repo.DB.Where("id IN ?", deviceIDs).Find(&devices).Error; err != nil {
		return nil, err
	}
	for i := range devices {
		types, _ := u.TypeMapRepo.GetDeviceTypes(devices[i].ID)
		devices[i].Types = types
	}
	return devices, nil
}

func (u *DeviceUsecase) GetAllDevicesWithTypesAndLocation() ([]domain.Device, error) {
	devices, err := u.Repo.GetAllDevices()
	if err != nil {
		return nil, err
	}
	for i := range devices {
		types, _ := u.TypeMapRepo.GetDeviceTypes(devices[i].ID)
		devices[i].Types = types
		if devices[i].LocationID != 0 {
			loc, _ := u.Repo.DB.Model(&domain.Location{}).Where("id = ?", devices[i].LocationID).First(&domain.Location{}).Rows()
			if loc != nil {
				var location domain.Location
				devices[i].Location = &location
			}
		}
	}
	return devices, nil
}

func (u *DeviceUsecase) GetDevicesByLocation(locationID uint) ([]domain.Device, error) {
	var devices []domain.Device
	if err := u.Repo.DB.Where("location_id = ?", locationID).Find(&devices).Error; err != nil {
		return nil, err
	}
	for i := range devices {
		types, _ := u.TypeMapRepo.GetDeviceTypes(devices[i].ID)
		devices[i].Types = types
		if devices[i].LocationID != 0 {
			loc, _ := u.Repo.DB.Model(&domain.Location{}).Where("id = ?", devices[i].LocationID).First(&domain.Location{}).Rows()
			if loc != nil {
				var location domain.Location
				devices[i].Location = &location
			}
		}
	}
	return devices, nil
}
