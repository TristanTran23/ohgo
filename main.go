package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get API key from environment variables
	apiKey := os.Getenv("OHGO_APIKEY")
	if apiKey == "" {
		log.Fatalf("API key not found in environment variables")
	}

	// Define API endpoint and parameters
	baseURL := "https://api.ohgo.com/v1/cameras"
  req, err := http.NewRequest("GET", baseURL, nil)
  if err != nil {
    log.Fatal(err)
  }

}

