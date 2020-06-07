package storage

import (
  "backuper/logger"
  "context"
  "net/http"
  "os"
  "path/filepath"

  yadisk "github.com/nikitaksv/yandex-disk-sdk-go"
)

type YandexDisk struct {
  AccessToken string `yaml:"access_token"`
  Keep        int    `yaml:"keep"`
}

func (y YandexDisk) Perform(model, path string) error {
  context := context.Background()
  token := &yadisk.Token{AccessToken: y.AccessToken}

  yaDisk, err := yadisk.NewYaDisk(context, http.DefaultClient, token)
  if err != nil {
    return err
  }

  folder := "app:/" + model
  yaDisk.CreateResource(folder, []string{})

  resource := folder + "/" + filepath.Base(path)
  logger.Info("Upload resource %v", resource)

  link, err := yaDisk.GetResourceUploadLink(resource, []string{}, false)
  if err != nil {
    return err
  }

  data, err := os.Open(path)
  if err != nil {
    return err
  }

  client := &http.Client{}
  req, err := http.NewRequest(http.MethodPut, link.Href, data)
  if err != nil {
    return err
  }

  _, err = client.Do(req)
  if err != nil {
    return err
  }

  return y.Cycle(yaDisk, folder)
}

func (y YandexDisk) Cycle(yaDisk yadisk.YaDisk, folder string) error {
  res, err := yaDisk.GetResource(folder, []string{}, 1000, 0, false, "", "-created")
  if err != nil {
    return err
  }

  items := res.Embedded.Items
  count := len(items)
  if count <= y.Keep {
    return nil
  }

  logger.Info("Keep %v of %v", y.Keep, count)

  for _, item := range items[y.Keep:] {
    path := folder + "/" + item.Name

    logger.Info("Remove %v", path)

    _, err := yaDisk.DeleteResource(path, []string{}, false, "", true)
    if err != nil {
      return err
    }
  }

  return nil
}
