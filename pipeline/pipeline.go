package pipeline

import (
  "backuper/config"
  "backuper/logger"
  "backuper/utils"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "time"
)

type Pipeline struct {
  Error  bool
  Temp   string
  Start  time.Time
  Config config.Config
}

func (p *Pipeline) Run() {
  defer func() {
    duration := time.Since(p.Start)
    logger.Info("Finished %v", duration)
  }()

  defer func() {
    logger.Info("Cleanup %v", p.Temp)
    os.RemoveAll(p.Temp)
  }()

  p.PrepareArchives()
  p.PrepareDatabases()

  p.Build()

  if !p.Error {
    p.Store()
  }
}

func (p Pipeline) PackName() string {
  return p.Config.Name + "_" + p.Start.Format("2006.01.02.15.04.05")
}

func (p Pipeline) PackPath() string {
  return filepath.Join(p.Temp, p.PackName())
}

func (p Pipeline) BuildName() string {
  return p.PackName() + ".tar"
}

func (p Pipeline) BuildPath() string {
  return filepath.Join(p.Temp, p.BuildName())
}

func (p *Pipeline) LogError(err string) {
  p.Error = true
  logger.Error(strings.TrimSpace(err))
}

func (p *Pipeline) PrepareArchives() {
  if !p.Config.IsArchive() {
    return
  }

  logger.Info("=== Archives")

  dir := filepath.Join(p.PackPath(), "archives")
  logger.Info("Make directory %v", dir)
  if err := os.MkdirAll(dir, os.ModePerm); err != nil {
    p.LogError(err.Error())
    return
  }

  for _, db := range p.Config.Archives {
    logger.Info("Archive %v", db.Name)
    if err := db.Perform(dir); err != nil {
      p.LogError(err.Error())
    }
  }
}

func (p *Pipeline) PrepareDatabases() {
  if !p.Config.IsDatabase() {
    return
  }

  logger.Info("=== Databases")

  dir := filepath.Join(p.PackPath(), "databases")
  logger.Info("Make directory %v", dir)
  if err := os.MkdirAll(dir, os.ModePerm); err != nil {
    p.LogError(err.Error())
    return
  }

  for _, db := range p.Config.Databases.PostgreSQL {
    logger.Info("PostgreSQL %v", db.Name)
    if err := db.Perform(dir); err != nil {
      p.LogError(err.Error())
    }
  }

  for _, db := range p.Config.Databases.MySQL {
    logger.Info("MySQL %v", db.Name)
    if err := db.Perform(dir); err != nil {
      p.LogError(err.Error())
    }
  }
}

func (p Pipeline) Build() {
  logger.Info("=== Builder")

  packPath := p.PackPath()
  if ok, _ := utils.ContainFile(packPath); !ok {
    p.LogError("Empty build " + packPath)
    return
  }

  cmd := exec.Command("tar", "-c", "-P", "-f", p.BuildPath(), "-C", p.Temp)
  if utils.IsGnuTar() {
    cmd.Args = append(cmd.Args, "--ignore-failed-read")
  }
  cmd.Args = append(cmd.Args, p.PackName())
  logger.Info("RUN %v", cmd)
  if out, err := cmd.CombinedOutput(); err != nil {
    p.LogError(string(out))
    return
  }
}

func (p Pipeline) Store() {
  if !p.Config.IsStorage() {
    return
  }

  logger.Info("=== Storages")

  for i, storage := range p.Config.Storages.YandexDisk {
    logger.Info("YandexDisk %v", i)
    err := storage.Perform(p.Config.Name, p.BuildPath())
    if err != nil {
      p.LogError(err.Error())
    }
  }
}

func (p Pipeline) Notify() {
  logger.Flush()
  for i, notifier := range p.Config.Notifies.Telegram {
    logger.Info("Notify Telegram %v", i)
    err := notifier.Perform(p.Config.Name, p.Start, p.Error, logger.Buffer)
    if err != nil {
      logger.Error(err.Error())
    }
  }
}
