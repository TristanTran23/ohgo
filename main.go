package main

import (
  "fmt"
  "os"
  "encoding/json"
  "io"
  "log"
  "net/http"
  "time"
)

type Camera struct {
  ID string `json:"id"`
  Description string `json:"description"`
  Location string `json:"location"`
  Direction string `json:"direction"`
  Roadway string `json:"roadway"`
  Status string `json:"status"`
  LastUpdated time.Time `json:"last_updated"`
  URL string `json:"url"`
}

type Location struct {
  Latitude float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
}

type Client struct {
  baseURL string
  httpClient *http.Client
}

func NewClient() *Client {
  return &Client{
    baseURL: "https://api.ohio.gov/v1",
    httpClient: &http.Client{
      Timeout: time.Second * 30,
    },
  }
}

func (c *Client) GetCameras() ([]Camera, error) {
  resp, err := c.httpClient.Get(fmt.Sprintf("%s/cameras", c.baseURL))
  if err != nil {
    return nil, fmt.Errorf("failed to get cameras: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
  }
  
  var cameras []Camera
    if err := json.NewDecoder(resp.Body).Decode(&cameras); err != nil {
    return nil, fmt.Errorf("failed to decode cameras: %w", err)
  }

  return cameras, nil
}

func (c *Client) GetCamraImage(camera *Camera) ([]byte, error) {
  resp, err := c.httpClient.Get(camera.URL)
  if err != nil {
    return nil, fmt.Errorf("failed to get camera image: %w", err)
  }

}
