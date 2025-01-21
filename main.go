package main

import (
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net/http"
  "time"
  "os"
)

type Camera struct {
  ID          string    `json:"id"`
  Description string    `json:"description"`
  Location    Location  `json:"location"`
  Direction   string    `json:"direction"`
  Roadway     string    `json:"roadway"`
  Status      string    `json:"status"`
  LastUpdated time.Time `json:"lastUpdated"`
  URL         string    `json:"url"`
}

type Location struct {
  Latitude  float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
}

type Client struct {
  baseURL    string
  httpClient *http.Client
}

func NewClient() *Client {
  return &Client{
    baseURL: "https://publicapi.ohgo.com/v1",
    httpClient: &http.Client{
      Timeout: 30 * time.Second,
    },
  }
}

func (c *Client) GetCameras() ([]Camera, error) {
  resp, err := c.httpClient.Get(fmt.Sprintf("%s/cameras", c.baseURL))
  if err != nil {
    return nil, fmt.Errorf("failed to fetch cameras: %w", err)
  }
  defer resp.Body.Close()

  // Read the raw response body
  bodyBytes, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read response body: %w", err)
  }

  // Print the raw response for debugging
  fmt.Printf("Raw response: %s\n", string(bodyBytes))

  // Check content type
  contentType := resp.Header.Get("Content-Type")
  fmt.Printf("Content-Type: %s\n", contentType)

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
  }

  var cameras []Camera
  if err := json.Unmarshal(bodyBytes, &cameras); err != nil {
    return nil, fmt.Errorf("failed to decode response: %w\nResponse body: %s", err, string(bodyBytes))
  }

  return cameras, nil
}

func (c *Client) GetCameraImage(camera *Camera) ([]byte, error) {
  req, err := http.NewRequest("GET", camera.URL, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request: %w", err)
  }
  
  resp, err := c.httpClient.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to fetch camera image: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
  }

  return io.ReadAll(resp.Body)
}

func main() {
  client := NewClient()
  
  fmt.Println("Attempting to fetch cameras...")
  cameras, err := client.GetCameras()
  if err != nil {
    log.Fatal("Error getting cameras:", err)
  }

  fmt.Printf("Found %d cameras\n", len(cameras))
  for _, camera := range cameras {
    fmt.Printf("Camera: %s at %s\n", camera.ID, camera.Description)
    
    imageData, err := client.GetCameraImage(&camera)
    if err != nil {
      log.Printf("Error getting image for camera %s: %v\n", camera.ID, err)
      continue
    }
    
    filename := fmt.Sprintf("camera_%s.jpg", camera.ID)
    if err := os.WriteFile(filename, imageData, 0644); err != nil {
      log.Printf("Error saving image for camera %s: %v\n", camera.ID, err)
      continue
    }
    
    fmt.Printf("Saved image for camera %s to %s\n", camera.ID, filename)
  }
}
