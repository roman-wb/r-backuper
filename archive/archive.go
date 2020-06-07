package archive

import (
  "backuper/logger"
  "backuper/utils"
  "errors"
  "os/exec"
  "path/filepath"
)

type Archive struct {
  Name string   `yaml:"name"`
  Args []string `yaml:"args"`
  Gzip bool     `yaml:"gzip"`
}

func (a Archive) Validate() error {
  if a.Name != "" && len(a.Args) != 0 {
    return nil
  }
  return errors.New("Archive required `name` and `args`")
}

func (a Archive) Perform(dir string) error {
  if err := a.Validate(); err != nil {
    return err
  }

  file := "Archive_" + a.Name + ".tar"
  path := filepath.Join(dir, file)

  dump := exec.Command("tar", "-c", "-f", path)
  if utils.IsGnuTar() {
    dump.Args = append(dump.Args, "--ignore-failed-read")
  }
  dump.Args = append(dump.Args, a.Args...)
  logger.Info("RUN %v", dump)
  if out, err := dump.CombinedOutput(); err != nil {
    return errors.New(string(out))
  }

  if a.Gzip {
    gzip := exec.Command("gzip", path)
    logger.Info("RUN %v", gzip)
    if out, err := gzip.CombinedOutput(); err != nil {
      return errors.New(string(out))
    }
  }

  return nil
}
