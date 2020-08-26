# *Test API Project in GoLang*

A GoLang RESTful API Application that accepts records/requests of type /geo and /device/info and prints them out to the console. It accepts GET, POST, PATCH, PUT, DELETE, HEAD, and OPTIONS requests. It then marshalls the requests into a GoStruct (geo object, and device info object) and sends them to MongoDB for storage. *Skills implemented include: GoLang, HTTP requests, MongoDB* 


## How to Run

```
$ go run main.go 
```

## Important
Add and env variable ```GOOGLE_APPLICATION_CREDENTIALS="/Users/kforemski/go/src/git.dev.kochava.com/KaylaAPI/key-file.json"```
with the path to your gcloud key-file.json
