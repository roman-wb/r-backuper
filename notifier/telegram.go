package notifier

import (
  "bytes"
  "errors"
  "io/ioutil"
  "mime/multipart"
  "net/http"
  "net/url"
  "time"
)

type Telegram struct {
  Token  string `yaml:"token"`
  ChatId string `yaml:"chat_id"`
}

func (t Telegram) Perform(name string, start time.Time, isError bool, buffer bytes.Buffer) error {
  if isError {
    if err := t.SendMessage("[Failure] " + name); err != nil {
      return err
    }
    name := name + "_" + start.Format("2006.01.02.15.04.05") + ".log"
    if err := t.SendLogs(name, buffer); err != nil {
      return err
    }
    return nil
  } else {
    return t.SendMessage("[Success] " + name)
  }
}

func (t Telegram) SendMessage(message string) error {
  values := url.Values{}
  values.Add("chat_id", t.ChatId)
  values.Add("text", message)
  values.Add("parse_mode", "html")

  url := "https://api.telegram.org/bot" + t.Token + "/sendMessage?" + values.Encode()
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  _, err = ioutil.ReadAll(resp.Body)

  return err
}

func (t Telegram) SendLogs(filename string, buffer bytes.Buffer) error {
  values := url.Values{}
  values.Add("chat_id", t.ChatId)

  url := "https://api.telegram.org/bot" + t.Token + "/sendDocument?" + values.Encode()

  body := new(bytes.Buffer)
  bodyWriter := multipart.NewWriter(body)
  part, err := bodyWriter.CreateFormFile("document", filename)
  if err != nil {
    return err
  }

  part.Write(buffer.Bytes())

  if err = bodyWriter.Close(); err != nil {
    return err
  }

  req, err := http.NewRequest("POST", url, body)
  req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
  if err != nil {
    return err
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }

  if bodyResp, err := ioutil.ReadAll(resp.Body); err != nil {
    return errors.New(string(bodyResp))
  }

  return nil
}
