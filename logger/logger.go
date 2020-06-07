package logger

import (
  "bufio"
  "bytes"

  "github.com/kpango/glg"
)

var Buffer bytes.Buffer
var bufferWriter *bufio.Writer

func init() {
  fileWriter := glg.FileWriter("backuper.log", 0644)
  bufferWriter = bufio.NewWriter(&Buffer)

  glg.Get().SetMode(glg.BOTH).AddWriter(fileWriter).AddWriter(bufferWriter)
}

func Flush() {
  bufferWriter.Flush()
}

func Info(msg string, args ...interface{}) {
  glg.Infof(msg, args...)
}

func Error(msg string, args ...interface{}) {
  glg.Errorf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
  glg.Fatalf(msg, args...)
}
