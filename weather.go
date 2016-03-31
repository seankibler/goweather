package main

import(
  "fmt"
  "flag"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "log"
  "os"
)

const apiURL = "http://api.openweathermap.org/data/2.5/weather?zip=%s,%s&APPID=%s"

var weatherAPIToken = flag.String("token", os.Getenv("OPENWEATHERMAP_APPID"), "openweathermap APPID")
var zipCode = flag.String("zip", "", "zip code for weather geography")
var countryCode = flag.String("country", "us", "country code (see http://openweathermap.org/API for valid possibilities)")
var testJson = flag.String("json", "", "read from json file instead of API")

func GetWeather() (fcst string) {
  var body []byte
  var err error

  if *zipCode == "" {
    log.Fatal("need a zip/postcal code")
  }

  if *testJson != "" {
    body, err = ioutil.ReadFile(*testJson)

    if err != nil {
      log.Fatalf("io error: %s\n", err.Error())
    }
  } else {
    if *weatherAPIToken == "" {
      log.Fatal("need an API token to use openweathermap API")
    }

    url := fmt.Sprintf(apiURL, *zipCode, *countryCode, weatherAPIToken)
    resp, err := http.Get(url)
    if err != nil {
      log.Fatalf("received http error %s", err)
    }
    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
  }

  type WeatherData struct {
    Id int `json:"id"`
    Code int32 `json:"cod"`
    Message string `json:"message"`
    Weather []struct {
      Id uint32 `json:"id"`
      Summary string `json:"main"`
      Description string `json:"description"`
    }
  }

  weatherData := WeatherData{}

  if err := json.Unmarshal(body, &weatherData); err != nil {
    log.Fatalf("json error: %s", err.Error())
  } else {
    if weatherData.Code != 200 {
      fcst = weatherData.Message
    } else {
      fcst = weatherData.Weather[0].Description
    }
  }

  return
}

// Says hello and prints today's weather forecast
func main() {
  flag.Parse()
  fmt.Printf("%s\n", GetWeather())
}
