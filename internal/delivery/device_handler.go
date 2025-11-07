package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simonaditiabbp/netmon-backend/internal/usecase"
)

type DeviceHandler struct {
	Usecase *usecase.DeviceUsecase
}

func NewDeviceHandler(usecase *usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{Usecase: usecase}
}

func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	devices, err := h.Usecase.GetAllDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) SSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Query devices from the database
	devices, err := h.Usecase.GetAllDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate statistics
	total := len(devices)
	online := 0
	for _, device := range devices {
		if device.Status == "online" {
			online++
		}
	}
	offline := total - online

	// Send statistics and devices
	c.SSEvent("message", gin.H{
		"total":   total,
		"online":  online,
		"offline": offline,
		"devices": devices,
	})
}
