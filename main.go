package main

import (
	"log"
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

	r.GET("/devices", handler.GetAllDevices)
	r.GET("/sse", handler.SSE)
	r.GET("/events", handler.SSE)

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
