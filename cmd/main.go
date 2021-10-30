package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"log"
	"os"
	"russt.io/varanus"
)

func main() {
	funcframework.RegisterEventFunction("/", varanus.CheckSiteAvailability)

	port := "8080"
	envPort := os.Getenv("FUNCTIONS_PORT")
	if envPort != "" {
		port = envPort
	}
	err := funcframework.Start(port)
	if err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
