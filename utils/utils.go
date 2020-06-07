package utils

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

func IsGnuTar() bool {
  out, _ := exec.Command("tar", "--version").CombinedOutput()
  return strings.Contains(string(out), "GNU")
}

func ContainFile(dir string) (bool, error) {
  var contain bool

  if _, err := os.Stat(dir); err != nil {
    return contain, err
  }

  filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if info.Mode().IsRegular() {
      contain = true
    }
    return nil
  })

  return contain, nil
}
