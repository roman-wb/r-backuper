package main

import (
  "backuper/config"
  "backuper/pipeline"
  "time"
)

func main() {
  config := config.Setup()

  pipeline := pipeline.Pipeline{
    Temp:   "temp",
    Start:  time.Now(),
    Config: config,
  }
  pipeline.Run()
  pipeline.Notify()
}
