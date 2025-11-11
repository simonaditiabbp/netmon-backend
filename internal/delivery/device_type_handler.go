package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/usecase"
)

type DeviceTypeHandler struct {
	Usecase *usecase.DeviceTypeUsecase
}

func NewDeviceTypeHandler(usecase *usecase.DeviceTypeUsecase) *DeviceTypeHandler {
	return &DeviceTypeHandler{Usecase: usecase}
}

func (h *DeviceTypeHandler) CreateDeviceType(c *gin.Context) {
	var dt domain.DeviceType
	if err := c.ShouldBindJSON(&dt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.Usecase.CreateDeviceType(&dt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Device type created successfully"})
}

func (h *DeviceTypeHandler) UpdateDeviceType(c *gin.Context) {
	id := c.Param("id")
	var dt domain.DeviceType
	if err := c.ShouldBindJSON(&dt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device type ID"})
		return
	}
	dt.ID = uint(idUint)
	if err := h.Usecase.UpdateDeviceType(&dt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device type updated successfully"})
}

func (h *DeviceTypeHandler) GetAllDeviceTypes(c *gin.Context) {
	types, err := h.Usecase.GetAllDeviceTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *DeviceTypeHandler) GetDeviceTypeByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device type ID"})
		return
	}
	dt, err := h.Usecase.GetDeviceTypeByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dt)
}

func (h *DeviceTypeHandler) DeleteDeviceType(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device type ID"})
		return
	}

	// Cek apakah masih dipakai di device_type_maps
	count, err := h.Usecase.CountDeviceTypeUsage(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device type masih digunakan oleh device, tidak bisa dihapus."})
		return
	}

	if err := h.Usecase.DeleteDeviceType(uint(idUint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device type deleted successfully"})
}
