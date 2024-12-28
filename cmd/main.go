package main

import (
	"log"
	"user-service/pkg/router"
)

func main() {

	r := router.SetUpRouter()

	if err := r.Run(":8030"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
