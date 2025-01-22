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

type CameraResponse struct {
    Links              []Link    `json:"links"`
    LastUpdated       string     `json:"lastUpdated"`
    AcceptedFilters   []Filter   `json:"acceptedFilters"`
    RejectedFilters   []Filter   `json:"rejectedFilters"`
    TotalPageCount    int        `json:"totalPageCount"`
    TotalResultCount  int        `json:"totalResultCount"`
    CurrentResultCount int        `json:"currentResultCount"`
    Results           []Camera   `json:"results"`
}

type Camera struct {
    Links       []Link       `json:"links"`
    ID          string       `json:"id"`
    Latitude    float64      `json:"latitude"`
    Longitude   float64      `json:"longitude"`
    Location    string       `json:"location"`
    Description string       `json:"description"`
    CameraViews []CameraView `json:"cameraViews"`
    MainRoute   string       `json:"mainRoute"`
}

type CameraView struct {
    Direction string `json:"direction"`
    SmallUrl  string `json:"smallUrl"`
    LargeUrl  string `json:"largeUrl"`
}

type Link struct {
    Href string `json:"href"`
    Rel  string `json:"rel"`
}

type Filter struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

func NewClient() *Client {
    apiKey := os.Getenv("OHGO_APIKEY")
    if apiKey == "" {
        log.Fatal("OHGO_APIKEY environment variable not set")
    }
    
    return &Client{
        baseURL: "https://publicapi.ohgo.com/api/v1",
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) GetAllCameras() ([]Camera, error) {
    var allCameras []Camera
    currentPage := 1
    
    for {
        pageURL := fmt.Sprintf("%s/cameras?page=%d", c.baseURL, currentPage)
        req, err := http.NewRequest("GET", pageURL, nil)
        if err != nil {
            return nil, fmt.Errorf("failed to create request: %w", err)
        }
        
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Authorization", fmt.Sprintf("APIKEY %s", c.apiKey))
        
        fmt.Printf("Fetching page %d...\n", currentPage)
        resp, err := c.httpClient.Do(req)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch cameras: %w", err)
        }
        
        bodyBytes, err := io.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            return nil, fmt.Errorf("failed to read response body: %w", err)
        }
        
        var response CameraResponse
        if err := json.Unmarshal(bodyBytes, &response); err != nil {
            return nil, fmt.Errorf("failed to decode response: %w\nResponse body: %s", err, string(bodyBytes))
        }
        
        allCameras = append(allCameras, response.Results...)
        fmt.Printf("Retrieved %d cameras from page %d\n", len(response.Results), currentPage)
        
        // If we've reached the last page, break
        if currentPage >= response.TotalPageCount {
            break
        }
        
        currentPage++
    }
    
    return allCameras, nil
}

func (c *Client) GetCameraImage(camera *Camera) ([]byte, error) {
    if len(camera.CameraViews) == 0 {
        return nil, fmt.Errorf("no camera views available for camera %s", camera.ID)
    }
    
    // Use the largeUrl from the first camera view
    imageURL := camera.CameraViews[0].LargeUrl
    req, err := http.NewRequest("GET", imageURL, nil)
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
    
    fmt.Println("Attempting to fetch all cameras...")
    cameras, err := client.GetAllCameras()
    if err != nil {
        log.Fatal("Error getting cameras:", err)
    }
    
    fmt.Printf("Successfully found %d cameras\n", len(cameras))
    
    // Create an output directory if it doesn't exist
    outputDir := "camera_images"
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        log.Fatal("Error creating output directory:", err)
    }
    
    for _, camera := range cameras {
        fmt.Printf("Processing camera: %s at %s\n", camera.ID, camera.Description)
        
        imageData, err := client.GetCameraImage(&camera)
        if err != nil {
            log.Printf("Error getting image for camera %s: %v\n", camera.ID, err)
            continue
        }
        
        filename := fmt.Sprintf("%s/camera_%s.jpg", outputDir, camera.ID)
        if err := os.WriteFile(filename, imageData, 0644); err != nil {
            log.Printf("Error saving image for camera %s: %v\n", camera.ID, err)
            continue
        }
        
        fmt.Printf("Saved image for camera %s to %s\n", camera.ID, filename)
    }
    
    fmt.Println("Processing complete!")
}
