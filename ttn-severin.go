package main

import (
    "github.com/TheThingsNetwork/go-utils/log"
    "github.com/TheThingsNetwork/go-utils/log/apex"
    "github.com/TheThingsNetwork/ttn/core/types"
    "github.com/TheThingsNetwork/ttn/mqtt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
    "fmt"
    "reflect"
    "os"
)

var configuration Configuration

func main() {
  getConfig(&configuration)

  ctx := apex.Stdout().WithField("Application", configuration.APP_ID)
  log.Set(ctx)
  //connect to TTN's app with MQTT
  client := mqtt.NewClient(ctx, "ttnctl", configuration.APP_ID, configuration.APP_KEY, configuration.TTN_URL)
  if err := client.Connect(); err != nil {
    ctx.WithError(err).Fatal("Could not connect")
  }
  //subscribe to and store upink messages in a channel
  msg := make(chan types.UplinkMessage)
  token := client.SubscribeDeviceUplink(configuration.APP_ID, configuration.DEV_ID, func(client mqtt.Client, appID string, devID string, req types.UplinkMessage) {
    msg <- req
  })
  token.Wait();
  if err := token.Error(); err != nil {
    ctx.WithError(err).Fatal("Could not subscribe")
  }
  //analyse and send payload to opensensors.io
  for i := 0 ; configuration.LOOP_NB == -1 || i < configuration.LOOP_NB; i++ {
  	payload := <-msg
  	browseAndForward(payload.PayloadFields, configuration.ROOT_TOPIC)
	}
	client.Disconnect()
}

func getConfig(config *Configuration) {
  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)
  err := decoder.Decode(&configuration)
  if err != nil {
    fmt.Println("error:", err)
  }
}

//browseAndForward browse the payload to find data and return values with their path
func browseAndForward(fields map[string]interface{}, path string) {
	for key, value := range fields {
		if reflect.TypeOf(value).String() == "map[string]interface {}" {
			browseAndForward(value.(map[string]interface{}), path + "/" + key)
		} else {
			fmt.Println("Key:", key, "Value:", value, "Path:", path)
      dataBytes, _ := json.Marshal(value)
      dataString := string(dataBytes[:])
      publish(path + "/" + key, dataString)
		}
	}
}

//publish establish connexon with opensensors to the right topic
func publish(path string, data string) {
	login := "?client-id=" + configuration.OS_ID + "&password=" + configuration.OS_PWD
	url := configuration.OS_ENDPOINT + path + login

  mapData := map[string]interface{}{"data": data}
  fmt.Println(mapData)
  jsonStr, _ := json.Marshal(mapData)

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "api-key " + configuration.API_KEY)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
}
