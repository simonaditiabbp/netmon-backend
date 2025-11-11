package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simonaditiabbp/netmon-backend/internal/db"
	"github.com/simonaditiabbp/netmon-backend/internal/delivery"
	"github.com/simonaditiabbp/netmon-backend/internal/repository"
	"github.com/simonaditiabbp/netmon-backend/internal/usecase"
)

func main() {
	// Initialize database connection
	database := db.InitDB()

	deviceRepo := repository.NewDeviceRepository(database)
	deviceTypeRepo := repository.NewDeviceTypeRepository(database)
	deviceTypeMapRepo := repository.NewDeviceTypeMapRepository(database)

	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo, deviceTypeMapRepo, deviceTypeRepo)
	deviceTypeUsecase := usecase.NewDeviceTypeUsecase(deviceTypeRepo)

	deviceHandler := delivery.NewDeviceHandler(deviceUsecase)
	deviceTypeHandler := delivery.NewDeviceTypeHandler(deviceTypeUsecase)

	// Gin router setup
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.GET("/devices", deviceHandler.GetAllDevices)
	r.GET("/sse", deviceHandler.SSE)
	r.GET("/live", deviceHandler.GetAllLiveDevices) // without sse
	r.POST("/devices", deviceHandler.InsertDevice)
	r.PUT("/devices/:id", deviceHandler.UpdateDevice)
	r.GET("/devices/:id", deviceHandler.GetDeviceByID)
	r.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	r.POST("/devices_types", deviceTypeHandler.CreateDeviceType)
	r.PUT("/devices_types/:id", deviceTypeHandler.UpdateDeviceType)
	r.GET("/devices_types", deviceTypeHandler.GetAllDeviceTypes)
	r.GET("/devices_types/:id", deviceTypeHandler.GetDeviceTypeByID)
	r.DELETE("/devices_types/:id", deviceTypeHandler.DeleteDeviceType)

	r.GET("/devices/by-type", deviceHandler.GetDevicesByType)
	r.GET("/devices/by-types", deviceHandler.GetDevicesByTypeMulti)

	// Periodically check device statuses
	go func() {
		log.Println("Starting device status update...")
		for {
			deviceHandler.Usecase.UpdateDeviceStatus()
			log.Println("Device status update completed.")
			time.Sleep(5 * time.Second)
		}
	}()

	// Start a goroutine to process broadcast messages
	// go func() {
	// 	for message := range usecase.BroadcastChannel {
	// 		log.Printf("Broadcast message: %+v", message)
	// 	}
	// }()

	// Start server
	r.Run(":8082")
}
