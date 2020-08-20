package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if r.Method == "HEAD" {
		fmt.Println("/device/info HEAD Request recieved!")
	} else {
		fmt.Println("/device/info GET Request recieved!")
	}
}

func returnAllDeviceInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All DeviceInfos Endpoint")
	json.NewEncoder(w).Encode(DeviceInfos)
	if r.Method == "HEAD" {
		fmt.Println("/geo HEAD Request recieved!")
	} else {
		fmt.Println("/geo GET Request recieved!")
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!\n\n")
	fmt.Fprint(w, "To see all the records of type /geo go to:\t\thttp://localhost:10000/geos\nTo see all the records of type /device/info go to:\thttp://localhost:10000/device/infos\n\n")
	fmt.Fprint(w, "To see a specific /geo record, go to:\t\thttp://localhost:10000/geo/{device_id}\nTo see a specific /device/info record, go to:\thttp://localhost:10000/device/info/{device_id}")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/geos", returnAllGeos)
	myRouter.HandleFunc("/device/infos", returnAllDeviceInfos)
	myRouter.HandleFunc("/geo", returnAllGeos).Methods("HEAD")
	myRouter.HandleFunc("/device/info", returnAllDeviceInfos).Methods("HEAD")
	myRouter.HandleFunc("/geo", createNewGeo).Methods("POST")
	myRouter.HandleFunc("/device/info", createNewDeviceInfo).Methods("POST")
	myRouter.HandleFunc("/geo/{device_id}", deleteGeo).Methods("DELETE")
	myRouter.HandleFunc("/device/info/{device_id}", deleteDeviceInfo).Methods("DELETE")
	myRouter.HandleFunc("/geo/{device_id}", updateGeo).Methods("PATCH")
	myRouter.HandleFunc("/device/info/{device_id}", updateDeviceInfo).Methods("PATCH")
	myRouter.HandleFunc("/geo/{device_id}", updateGeo).Methods("PUT")
	myRouter.HandleFunc("/device/info/{device_id}", updateDeviceInfo).Methods("PUT")
	myRouter.HandleFunc("/geo/{device_id}", returnSingleGeo)
	myRouter.HandleFunc("/device/info/{device_id}", returnSingleDeviceInfo)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnSingleGeo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["device_id"]
	fmt.Fprintf(w, "Key: "+key)
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
	for _, device := range DeviceInfos {
		if device.DeviceID == key {
			json.NewEncoder(w).Encode(device)
		}
	}
}

func createNewGeo(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request and return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var geo Geo
	json.Unmarshal(reqBody, &geo)

	//Print to console
	fmt.Printf("/geo POST Request: %+v\n", geo)

	//Append to list of Geos and print to local host
	Geos = append(Geos, geo)
	json.NewEncoder(w).Encode(geo)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func createNewDeviceInfo(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request and return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var deviceInfo DeviceInfo
	json.Unmarshal(reqBody, &deviceInfo)

	//Print to console
	fmt.Printf("/device/info POST Request: %+v\n", deviceInfo)

	//Append to list of DeviceInfos and print to local host
	DeviceInfos = append(DeviceInfos, deviceInfo)
	json.NewEncoder(w).Encode(deviceInfo)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func deleteGeo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// extract the `id` of the geo we want to delete
	id := vars["device_id"]

	// loop through all our geos, if id matches, print to console and delete
	for index, geo := range Geos {
		if geo.DeviceID == id {
			//Print to console
			fmt.Printf("/geo DELETE Request: %+v\n", geo)
			Geos = append(Geos[:index], Geos[index+1:]...)
		}
	}
}

func deleteDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// extract the `id` of the device info we want to delete
	id := vars["device_id"]

	// loop through all our device infos, if id matches, print to console and delete
	for index, deviceInfo := range DeviceInfos {
		if deviceInfo.DeviceID == id {
			//Print to console
			fmt.Printf("/device/info DELETE Request: %+v\n", deviceInfo)
			DeviceInfos = append(DeviceInfos[:index], DeviceInfos[index+1:]...)
		}
	}
}

func updateGeo(w http.ResponseWriter, r *http.Request) {
	geoID := mux.Vars(r)["device_id"]
	var updatedGeo Geo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Only enter data with the latitude, longitude and ip_address in order to update")
	}
	json.Unmarshal(reqBody, &updatedGeo)
	//Print to console
	fmt.Print("/geo PATCH Request for DeviceID " + geoID + ": ")
	fmt.Printf("%+v\n", updatedGeo)

	for i, singleGeo := range Geos {
		if singleGeo.DeviceID == geoID {
			singleGeo.Latitude = updatedGeo.Latitude
			singleGeo.Longitude = updatedGeo.Longitude
			singleGeo.IPAddress = updatedGeo.IPAddress
			Geos = append(Geos[:i], singleGeo)
			json.NewEncoder(w).Encode(singleGeo)
		}
	}
}

func updateDeviceInfo(w http.ResponseWriter, r *http.Request) {
	deviceInfoID := mux.Vars(r)["device_id"]
	var updatedDeviceInfo DeviceInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Only enter data with the user_agent, ip_address, battery_level, and screen_orientation in order to update")
	}
	json.Unmarshal(reqBody, &updatedDeviceInfo)
	//Print to console
	fmt.Print("/device/info PATCH Request for DeviceID " + deviceInfoID + ": ")
	fmt.Printf("%+v\n", updatedDeviceInfo)

	for i, singleDeviceInfo := range DeviceInfos {
		if singleDeviceInfo.DeviceID == deviceInfoID {
			singleDeviceInfo.UserAgent = updatedDeviceInfo.UserAgent
			singleDeviceInfo.IPAddress = updatedDeviceInfo.IPAddress
			singleDeviceInfo.BatteryLevel = updatedDeviceInfo.BatteryLevel
			singleDeviceInfo.ScreenOrientation = updatedDeviceInfo.ScreenOrientation
			DeviceInfos = append(DeviceInfos[:i], singleDeviceInfo)
			json.NewEncoder(w).Encode(singleDeviceInfo)
		}
	}
}

func main() {
	fmt.Println("Kayla's Test Rest API Running")
	fmt.Println("Go to http://localhost:10000/ to see homepage")
	DeviceInfos = []DeviceInfo{
		DeviceInfo{DeviceID: "JENNA", UserAgent: "Kayla", IPAddress: "129.232.23.121", BatteryLevel: "87%", ScreenOrientation: "vertical"},
		DeviceInfo{DeviceID: "KAYLA", UserAgent: "Jenna", IPAddress: "120.112.19.333", BatteryLevel: "17%", ScreenOrientation: "horizontal"},
	}
	Geos = []Geo{
		Geo{DeviceID: "1234", Latitude: "48.121", Longitude: "127.12", IPAddress: "129.232.23.121"},
	}
	handleRequests()
}
