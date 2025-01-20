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


