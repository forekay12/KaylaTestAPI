# *Test API Project in GoLang*

A GoLang RESTful API Application that accepts records/requests of type /geo and /device/info and prints them out to the console. It accepts GET, POST, PATCH, PUT, DELETE, HEAD, and OPTIONS requests. It then marshalls the requests into a GoStruct (geo object, and device info object) and sends them to MongoDB for storage. All POST requests are sent to gcloud pubsub, and another go application pulls them out. *Skills implemented include: GoLang, HTTP requests, MongoDB, gcloud pubsub* 


## How to Run
Navigate to /KaylaAPI folder:

```
$ go run main.go 
```
Navigate to /KaylaAPI/PULL folder:
```
$ go run pull.go 
```

## Important
Add and env variable ```GOOGLE_APPLICATION_CREDENTIALS="/Users/kforemski/go/src/git.dev.kochava.com/KaylaAPI/key-file.json"```
with the path to your gcloud key-file.json
