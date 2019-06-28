package main;

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "os"

  . "transcription/backend-protobuf/go"

  "github.com/joho/godotenv"
  "github.com/golang/protobuf/proto"
  "github.com/nats-io/go-nats"
)

var apiKey string
var natsHost string

var nc *nats.Conn

type AudioOptions struct {
  Content []byte `json:"content"`
}

type AudioConfigOptions struct {
  LanguageCode string `json:"languageCode"`
}

type AudioConfig struct {
  Audio AudioOptions `json:"audio"`
  Config AudioConfigOptions `json:"config"`
}

type AudioTranscript struct {
  Transcript string `json:"transcript"`
  Confidence float32 `json:"confidence"`
}

type AudioResult struct {
  Alternatives []AudioTranscript `json:"alternatives"`
}

type AudioResults struct {
  Results []AudioResult `json:"results"`
}

func main() {
  // Load .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  natsHost = os.Getenv("NATS")
  apiKey = os.Getenv("API_KEY")

  //NATS
  nc, err := nats.Connect(natsHost)
  if err != nil {
    log.Fatal(err)
    return
  }

  nc.Subscribe("bite", NewBite)

  log.Printf("listening on nats")
  select { }
}

func NewBite(m *nats.Msg) {
  bite := Bite{}
  if err := proto.Unmarshal(m.Data, &bite); err != nil {
    log.Println(err)
    return
  }

  //TODO: Check cache (store) for existing transcription

  // Base64 encode audio bytes
  audioEncoded := base64.StdEncoding.EncodeToString(bite.Data)

  config := AudioConfig {
    Audio: AudioOptions {
      Content: []byte(audioEncoded),
    },
    Config: AudioConfigOptions {
      LanguageCode: "en-US",
    },
  }

  configJson, err := json.Marshal(config)
  if err != nil {
    log.Println(err)
    return
  }

  url := "https://speech.googleapis.com/v1/speech:recognize?key=" + apiKey
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(configJson))
  req.Header.Set("Content-Type", "application/json")

  client := http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Println(err)
    return
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println(err)
    return
  }

  results := AudioResults{}
  err = json.Unmarshal(body, &results)
  if err != nil {
    log.Println(err)
    return
  }

  if len(results.Results) > 0 && len(results.Results[0].Alternatives) > 0 {
    bite.Data = []byte(results.Results[0].Alternatives[0].Transcript)

    storeRequest := Store {
      Type: "transcription",
      Bite: &bite,
    }
    storeRequestBytes, err := proto.Marshal(&storeRequest)
    if err != nil {
      return
    }

    nc.Publish("new_store", storeRequestBytes)
  } else {
    // 404
    log.Println("google api could not be reached")
  }
}

