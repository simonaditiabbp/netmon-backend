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

	// Repository, Usecase, and Handler initialization
	repo := repository.NewDeviceRepository(database)
	usecase := usecase.NewDeviceUsecase(repo)
	handler := delivery.NewDeviceHandler(usecase)

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

	r.GET("/devices", handler.GetAllDevices)
	r.GET("/sse", handler.SSE)
	r.GET("/live", handler.GetAllLiveDevices) // without sse
	r.POST("/devices", handler.InsertDevice)
	r.PUT("/devices/:id", handler.UpdateDevice)
	r.GET("/devices/:id", handler.GetDeviceByID)
	r.DELETE("/devices/:id", handler.DeleteDevice)

	// Periodically check device statuses
	go func() {
		log.Println("Starting device status update...")
		for {
			handler.Usecase.UpdateDeviceStatus()
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
