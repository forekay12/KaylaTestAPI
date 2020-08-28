package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Geo struct {
	DeviceID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Latitude  string             `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude string             `json:"longitude,omitempty" bson:"longitude,omitempty"`
	IPAddress string             `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
}

type DeviceInfo struct {
	DeviceID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserAgent         string             `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	IPAddress         string             `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	BatteryLevel      string             `json:"battery_level,omitempty" bson:"battery_level,omitempty"`
	ScreenOrientation string             `json:"screen_orientation,omitempty" bson:"screen_orientation,omitempty"`
}

var client *mongo.Client

// explicit reads credentials from the specified path.
func explicit(jsonPath, projectID string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}

//publish pushes messages into the specified pubsub topic
func publish(b []byte) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "cloud-test-287516")
	if err != nil {
		fmt.Println(err)
	}
	topic := client.Topic("kayla")
	defer topic.Stop()
	var results []*pubsub.PublishResult
	r := topic.Publish(ctx, &pubsub.Message{Data: b})
	results = append(results, r)
	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
		fmt.Printf("And request body: %s\n", b)
	}
}

//returnAllGeos is the GET endpoint for /geos
func returnAllGeos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var geos []Geo
	collection := client.Database("kaylatestapi").Collection("geo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var geo Geo
		cursor.Decode(&geo)
		geos = append(geos, geo)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(geos)
	fmt.Println("Endpoint Hit: All Geos Endpoint")
	fmt.Println("/geos "+r.Method+" request recieved: ", geos)
}

//returnAllDeviceInfos is the GET endpoint for /device/infos
func returnAllDeviceInfos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var deviceInfos []DeviceInfo
	collection := client.Database("kaylatestapi").Collection("deviceinfo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var deviceInfo DeviceInfo
		cursor.Decode(&deviceInfo)
		deviceInfos = append(deviceInfos, deviceInfo)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(deviceInfos)
	fmt.Println("Endpoint Hit: All DeviceInfos Endpoint")
	fmt.Println("/device/infos "+r.Method+" request recieved: ", deviceInfos)
}

//homePage is the endpoint for /
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!\n\n")
	fmt.Fprint(w, "To see all the records of type /geo go to:\t\thttp://localhost:11000/geos\nTo see all the records of type /device/info go to:\thttp://localhost:11000/device/infos\n\n")
	fmt.Fprint(w, "To see a specific /geo record, go to:\t\thttp://localhost:11000/geo/{device_id}\nTo see a specific /device/info record, go to:\thttp://localhost:11000/device/info/{device_id}")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {
	explicit("/Users/kforemski/go/src/git.dev.kochava.com/KaylaAPI/key-file.json", "cloud-test-287516")
	fmt.Println("Starting Kayla's Test Rest API...")
	fmt.Println("Go to http://localhost:11000/ to see homepage")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	handleRequests()
}

//handleRequests handles all types of requests and directs them to appropriate endpoints
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/geos", returnAllGeos)
	myRouter.HandleFunc("/device/infos", returnAllDeviceInfos)
	myRouter.HandleFunc("/geo", returnAllGeos).Methods("HEAD")
	myRouter.HandleFunc("/device/info", returnAllDeviceInfos).Methods("HEAD")
	myRouter.HandleFunc("/geo", createNewGeo).Methods("POST")
	myRouter.HandleFunc("/device/info", createNewDeviceInfo).Methods("POST")
	myRouter.HandleFunc("/geo/{id}", deleteGeo).Methods("DELETE")
	myRouter.HandleFunc("/device/info/{id}", deleteDeviceInfo).Methods("DELETE")
	myRouter.HandleFunc("/geo/{id}", updateGeo).Methods("PATCH")
	myRouter.HandleFunc("/device/info/{id}", updateDeviceInfo).Methods("PATCH")
	myRouter.HandleFunc("/geo/{id}", updateGeo).Methods("PUT")
	myRouter.HandleFunc("/device/info/{id}", updateDeviceInfo).Methods("PUT")
	myRouter.HandleFunc("/geo/{id}", returnSingleGeo)
	myRouter.HandleFunc("/device/info/{id}", returnSingleDeviceInfo)
	log.Fatal(http.ListenAndServe(":11000", myRouter))
}

//returnSingleGeo is the GET endpoint for /geo/{id}
func returnSingleGeo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	key, _ := primitive.ObjectIDFromHex(vars["id"])
	var geo Geo
	collection := client.Database("kaylatestapi").Collection("geo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Geo{DeviceID: key}).Decode(&geo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(geo)
	//Print to console
	fmt.Printf("/geo GET Request: %+v\n", geo)
}

//returnSingleGeo is the GET endpoint for /geo/{id}
func returnSingleDeviceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	key, _ := primitive.ObjectIDFromHex(vars["id"])
	var deviceInfo DeviceInfo
	collection := client.Database("kaylatestapi").Collection("deviceinfo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, DeviceInfo{DeviceID: key}).Decode(&deviceInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(deviceInfo)
	//Print to console
	fmt.Printf("/device/info GET Request: %+v\n", deviceInfo)
}

//createNewGeo is a POST endpoint for /geo
//It also publishes the request to pubsub
func createNewGeo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var geo Geo
	_ = json.NewDecoder(r.Body).Decode(&geo)
	if geo.Latitude == "" || geo.Longitude == "" {
		fmt.Println("/geo POST Request FAILED")
		fmt.Println("latitude and longitude are required values!")
		return
	}
	//Print to console
	fmt.Printf("/geo POST Request: %+v\n", geo)
	//Publish request to gcloud pubsub
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(geo)
	publish(reqBodyBytes.Bytes())

	collection := client.Database("kaylatestapi").Collection("geo")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, geo)
	json.NewEncoder(w).Encode(result)
}

//createNewDeviceInfo is a POST endpoint for /deviceInfo
//It also publishes the request to pubsub
func createNewDeviceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var deviceInfo DeviceInfo
	_ = json.NewDecoder(r.Body).Decode(&deviceInfo)
	if deviceInfo.UserAgent == "" {
		fmt.Println("/device/info POST Request FAILED")
		fmt.Println("user agent is a required value!")
		return
	}
	//Print to console
	fmt.Printf("/device/info POST Request: %+v\n", deviceInfo)
	//Publish request to gcloud pubsub
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(deviceInfo)
	publish(reqBodyBytes.Bytes())

	collection := client.Database("kaylatestapi").Collection("deviceinfo")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, deviceInfo)
	json.NewEncoder(w).Encode(result)
}

//deleteGeo is a DELETE request for the endpoint /geo/{id}
func deleteGeo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	collection := client.Database("kaylatestapi").Collection("geo")
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
	//Print to console
	fmt.Printf("/geo DELETE Request with "+id.String()+": %v\n", result)
}

//deleteDeviceInfo is a DELETE request for the endpoing /device/info/{id}
func deleteDeviceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	collection := client.Database("kaylatestapi").Collection("deviceinfo")
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
	//Print to console
	fmt.Printf("/device/info DELETE Request with "+id.String()+": %v\n", result)
}

//updateGeo is the PATCH and PUT request endpoint for /geo/{id}
//It also publishes the request to pubsub
func updateGeo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var geo Geo
	_ = json.NewDecoder(r.Body).Decode(&geo)
	if geo.Latitude == "" || geo.Longitude == "" {
		fmt.Println("/geo POST Request FAILED")
		fmt.Println("latitude and longitude are required values!")
		return
	}
	//Print to console
	fmt.Printf("/geo "+r.Method+" Request: %+v\n", geo)
	//Publish request to gcloud pubsub
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(geo)
	publish(reqBodyBytes.Bytes())
	vars := mux.Vars(r)
	collection := client.Database("kaylatestapi").Collection("geo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	key, _ := primitive.ObjectIDFromHex(vars["id"])
	result := collection.FindOneAndReplace(ctx, bson.M{"_id": key}, geo)
	json.NewEncoder(w).Encode(result)
}

//updateDeviceInfo is the PATCH and PUT request endpoint for /device/info/{id}
//It also publishes the request to pubsub
func updateDeviceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var deviceInfo DeviceInfo
	_ = json.NewDecoder(r.Body).Decode(&deviceInfo)
	if deviceInfo.UserAgent == "" {
		fmt.Println("/device/info POST Request FAILED")
		fmt.Println("user agent is a required value!")
		return
	}
	//Print to console
	fmt.Printf("/device/info "+r.Method+" Request: %+v\n", deviceInfo)
	//Publish request to gcloud pubsub
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(deviceInfo)
	publish(reqBodyBytes.Bytes())
	vars := mux.Vars(r)
	collection := client.Database("kaylatestapi").Collection("deviceinfo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	result := collection.FindOneAndReplace(ctx, bson.M{"_id": id}, deviceInfo)
	json.NewEncoder(w).Encode(result)
}
