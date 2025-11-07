package usecase

import (
	"fmt"
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
	BroadcastChannel chan gin.H
}

func NewDeviceUsecase(repo *repository.DeviceRepository) *DeviceUsecase {
	return &DeviceUsecase{
		Repo:             repo,
		BroadcastChannel: make(chan gin.H, 30), // Initialize buffered channel
	}
}

func (u *DeviceUsecase) GetAllDevices() ([]domain.Device, error) {
	return u.Repo.GetAllDevices()
}

func (u *DeviceUsecase) CreateDevice(device *domain.Device) error {
	return u.Repo.CreateDevice(device)
}

func (u *DeviceUsecase) UpdateDeviceStatus() {
	log.Println("Starting UpdateDeviceStatus...")
	defer log.Println("Finished UpdateDeviceStatus.")

	devices, err := u.Repo.GetAllDevices()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	log.Println("Checking device statuses...")

	for _, device := range devices {
		fmt.Println("Checking device:", device.Name)
		// Check device status using HTTP or ping
		isOnline := false

		if strings.HasPrefix(device.IP, "http") {
			log.Printf("Sending HTTP request to URL: %s", device.IP)
			resp, err := http.Get(device.IP)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				isOnline = true
			} else {
				log.Printf("HTTP request failed for URL %s: %v", device.IP, err)
			}
		} else {
			log.Printf("Pinging IP: %s", device.IP)
			cmd := exec.Command("ping", "-n", "1", device.IP)
			if err := cmd.Run(); err == nil {
				isOnline = true
			} else {
				log.Printf("Ping failed for IP %s: %v", device.IP, err)
			}
		}

		fmt.Println("Device:", device.Name, "isOnline:", isOnline)

		// Update status and broadcast changes
		oldStatus := device.Status
		device.Status = map[bool]string{true: "online", false: "offline"}[isOnline]
		fmt.Println("Device:", device.Name, "oldStatus:", oldStatus, "newStatus:", device.Status)
		if oldStatus != device.Status {
			log := domain.Log{
				DeviceID:  device.ID,
				OldStatus: oldStatus,
				NewStatus: device.Status,
				Timestamp: time.Now(),
			}
			fmt.Println("Logging status change: ", log)
			u.Repo.CreateLog(&log)
		}
		device.LastOnline = time.Now()
		err := u.Repo.UpdateDevice(&device)
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
