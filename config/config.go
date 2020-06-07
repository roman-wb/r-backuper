package config

import (
  "backuper/archive"
  "backuper/database"
  "backuper/logger"
  "backuper/notifier"
  "backuper/storage"
  "flag"
  "io/ioutil"
  "log"

  yaml "gopkg.in/yaml.v3"
)

type Database struct {
  PostgreSQL []database.PostgreSQL `yaml:"postgresql"`
  MySQL      []database.MySQL      `yaml:"mysql"`
}

type Storage struct {
  YandexDisk []storage.YandexDisk `yaml:"yandex_disk"`
}

type Notifier struct {
  Telegram []notifier.Telegram `yaml:"telegram"`
}

type Config struct {
  Name      string            `yaml:"name"`
  Archives  []archive.Archive `yaml:"archives"`
  Databases Database          `yaml:"databases"`
  Storages  Storage           `yaml:"storages"`
  Notifies  Notifier          `yaml:"notifiers"`
}

func (c Config) IsArchive() bool {
  return len(c.Archives) != 0
}

func (c Config) IsDatabase() bool {
  return len(c.Databases.PostgreSQL) != 0 || len(c.Databases.MySQL) != 0
}

func (c Config) IsStorage() bool {
  return len(c.Storages.YandexDisk) != 0
}

func Setup() Config {
  var path string
  flag.StringVar(&path, "config", "config.yml", "Configuration file")
  flag.Parse()

  logger.Info("Load %v", path)
  config := Load(path)

  return config
}

func Load(path string) Config {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal(err)
  }

  config := Config{}
  if err := yaml.Unmarshal(data, &config); err != nil {
    log.Fatal(err)
  }

  return config
}
