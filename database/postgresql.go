package database

import (
  "backuper/logger"
  "errors"
  "os/exec"
  "path/filepath"
)

type PostgreSQL struct {
  Name string   `yaml:"name"`
  Env  []string `yaml:"env"`
  Args []string `yaml:"args"`
  Gzip bool     `yaml:"gzip"`
}

func (db PostgreSQL) Validate() error {
  if db.Name != "" && len(db.Args) != 0 {
    return nil
  }
  return errors.New("PostgreSQL required `name` and `args`")
}

func (db PostgreSQL) Perform(dir string) error {
  if err := db.Validate(); err != nil {
    return err
  }

  file := "PostgreSQL_" + db.Name + ".sql"
  path := filepath.Join(dir, file)

  dump := exec.Command("pg_dump", "--file="+path)
  dump.Args = append(dump.Args, db.Args...)
  dump.Env = append(dump.Env, db.Env...)
  logger.Info("RUN %v", dump)
  if out, err := dump.CombinedOutput(); err != nil {
    return errors.New(string(out))
  }

  if db.Gzip {
    gzip := exec.Command("gzip", path)
    logger.Info("RUN %v", gzip)
    if out, err := gzip.CombinedOutput(); err != nil {
      return errors.New(string(out))
    }
  }

  return nil
}
