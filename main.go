package main;

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "flag"
  "io/ioutil"
  "log"
  "net/http"
  "strconv"
  "time"

  "github.com/golang/protobuf/proto"
  "github.com/julienschmidt/httprouter"
  "github.com/nats-io/go-nats"
)

var listen string
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
  // Parse flags
  flag.StringVar(&listen, "listen", ":8080", "host and port to listen on")
  flag.StringVar(&apiKey, "api-key", "AIzaSyDxSXDefzw9gXCQaVzOCYlRn_vcC9Da9Q0", "Google Cloud API key")
  flag.StringVar(&natsHost, "nats", "nats://localhost:4222", "host and port of NATS")
  flag.Parse()

  //NATS
  nc, err := nats.Connect(natsHost)
  if err != nil {
    log.Fatal(err)
    return
  }

  nc.Subscribe("new_bite", NewBite)

  // Routes
	router := httprouter.New()

  router.GET("/transcription/:key/scan", ScanTranscription) // Scanning
  router.GET("/transcription/:key/start/:start", GetTranscription)

  // Start server
  log.Printf("starting server on %s", listen)
  log.Fatal(http.ListenAndServe(listen, router))
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
    errRes := Response {
      Code: 500,
      Message: []byte(http.StatusText(http.StatusInternalServerError)),
      Client: bite.Client,
    }
    errResBytes, err := proto.Marshal(&errRes)
    if err == nil {
      nc.Publish("res", errResBytes)
    }
    return
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println(err)
    errRes := Response {
      Code: 500,
      Message: []byte(http.StatusText(http.StatusInternalServerError)),
      Client: bite.Client,
    }
    errResBytes, err := proto.Marshal(&errRes)
    if err == nil {
      nc.Publish("res", errResBytes)
    }
    return
  }

  results := AudioResults{}
  err = json.Unmarshal(body, &results)
  if err != nil {
    log.Println(err)
    errRes := Response {
      Code: 500,
      Message: []byte(http.StatusText(http.StatusInternalServerError)),
      Client: bite.Client,
    }
    errResBytes, err := proto.Marshal(&errRes)
    if err == nil {
      nc.Publish("res", errResBytes)
    }
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
    errRes := Response {
      Code: 404,
      Message: []byte(http.StatusText(http.StatusNotFound)),
      Client: bite.Client,
    }
    errResBytes, err := proto.Marshal(&errRes)
    if err == nil {
      nc.Publish("res", errResBytes)
    }
  }
}

// Route handlers
func ParseStartString(start string) (uint64, error) {
	return strconv.ParseUint(start, 10, 64)
}

func ScanTranscription(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  from, err := ParseStartString(r.FormValue("from"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	to, err := ParseStartString(r.FormValue("to"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

  scanRequest := ScanRequest {
    Key: p.ByName("key"),
    From: from,
    To: to,
    Type: "transcription",
  }

  drBytes, err := proto.Marshal(&scanRequest);
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  msg, err := nc.Request("scan_store", drBytes, 10 * time.Second) // 10s timeout
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  res := Response {}
  if err := proto.Unmarshal(msg.Data, &res); err != nil {
    log.Println(err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  if res.Code == 200 {
    w.Header().Set("Content-Type", "application/json")
  	w.Write(res.Message)
  } else if len(res.Message) == 0 {
    http.Error(w, http.StatusText(int(res.Code)), int(res.Code))
  } else {
    http.Error(w, string(res.Message), int(res.Code))
  }
}

func GetTranscription(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  start, err := ParseStartString(p.ByName("start"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

  dataRequest := DataRequest {
    Key: p.ByName("key"),
    Start: start,
    Type: "transcription",
  }

  drBytes, err := proto.Marshal(&dataRequest);
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  msg, err := nc.Request("request_store", drBytes, 10 * time.Second) // 10s timeout
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  res := Response {}
  if err := proto.Unmarshal(msg.Data, &res); err != nil {
    log.Println(err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  if res.Code == 200 {
    w.Header().Add("Content-Type", "text/plain")
  	w.Write(res.Message)
  } else if len(res.Message) == 0 {
    http.Error(w, http.StatusText(int(res.Code)), int(res.Code))
  } else {
    http.Error(w, string(res.Message), int(res.Code))
  }
}
