package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/usecase"
)

type DeviceHandler struct {
	Usecase *usecase.DeviceUsecase
}

func NewDeviceHandler(usecase *usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{Usecase: usecase}
}

func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	devices, err := h.Usecase.GetAllDevicesWithTypes()
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

func (h *DeviceHandler) GetAllLiveDevices(c *gin.Context) {
	// Query devices from the database
	devices, err := h.Usecase.GetAllDevicesWithTypes()
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

	data := gin.H{
		"total":   total,
		"online":  online,
		"offline": offline,
		"devices": devices,
	}

	c.JSON(http.StatusOK, data)
}

func (h *DeviceHandler) InsertDevice(c *gin.Context) {
	var device domain.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.Usecase.InsertDeviceWithTypes(&device, device.TypeIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Device inserted successfully"})
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var device domain.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if idUint, err := strconv.ParseUint(id, 10, 32); err == nil {
		device.ID = uint(idUint)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	if err := h.Usecase.UpdateDeviceWithTypes(&device, device.TypeIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully"})
}

func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.Usecase.GetDeviceByIDWithTypes(uint(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	if err := h.Usecase.DeleteDevice(uint(idUint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func (h *DeviceHandler) GetDevicesByType(c *gin.Context) {
	typeIDStr := c.Query("type_id")
	if typeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type_id query param is required"})
		return
	}
	typeID, err := strconv.ParseUint(typeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type_id"})
		return
	}
	devices, err := h.Usecase.GetDevicesByType(uint(typeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) GetDevicesByTypeMulti(c *gin.Context) {
	typeIDsStr := c.QueryArray("type_ids")
	if len(typeIDsStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type_ids query param is required"})
		return
	}
	var typeIDs []uint
	for _, s := range typeIDsStr {
		id, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type_id: " + s})
			return
		}
		typeIDs = append(typeIDs, uint(id))
	}
	devices, err := h.Usecase.GetDevicesByTypeMulti(typeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) GetAllDevicesWithTypesAndLocation(c *gin.Context) {
	devices, err := h.Usecase.GetAllDevicesWithTypesAndLocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) GetDevicesByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location_id"})
		return
	}
	devices, err := h.Usecase.GetDevicesByLocation(uint(locationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}
