package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Geo struct {
	DeviceID  string `json:"device_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	IPAddress string `json:"ip_address"`
}

type DeviceInfo struct {
	DeviceID          string `json:"device_id"`
	UserAgent         string `json:"user_agent"`
	IPAddress         string `json:"ip_address"`
	BatteryLevel      string `json:"battery_level"`
	ScreenOrientation string `json:"screen_orientation"`
}

var Geos []Geo
var DeviceInfos []DeviceInfo

func returnAllGeos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Geos Endpoint")
	json.NewEncoder(w).Encode(Geos)
}

func returnAllDeviceInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All DeviceInfos Endpoint")
	json.NewEncoder(w).Encode(DeviceInfos)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Fprint(w, "Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/geo", returnAllGeos)
	myRouter.HandleFunc("/deviceinfo", returnAllDeviceInfos)
	myRouter.HandleFunc("/geo/{device_id}", returnSingleGeo)
	myRouter.HandleFunc("/deviceinfo/{device_id}", returnSingleDeviceInfo)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnSingleGeo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["device_id"]

	fmt.Fprintf(w, "Key: "+key)
	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, geo := range Geos {
		if geo.DeviceID == key {
			json.NewEncoder(w).Encode(geo)
		}
	}
}

func returnSingleDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["device_id"]

	fmt.Fprintf(w, "Key: "+key)
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	DeviceInfos = []DeviceInfo{
		DeviceInfo{DeviceID: "1234-123919291-123-12312", UserAgent: "Kayla", IPAddress: "129.232.23.121", BatteryLevel: "87%", ScreenOrientation: "vertical"},
	}
	Geos = []Geo{
		Geo{DeviceID: "1234-123919291-123-12312", Latitude: "48.121", Longitude: "127.12", IPAddress: "129.232.23.121"},
	}
	handleRequests()
}
