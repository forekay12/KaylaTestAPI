# *Test API Project in GoLang*

A GoLang RESTful API Application that accepts records/requests of type /geo and /device/info and prints them out to the console. It accepts GET, POST, PATCH, PUT, DELETE, HEAD, and OPTIONS requests. It then marshalls the requests into a GoStruct (geo object, and device info object) and sends them to MongoDB for storage. All POST requests are sent to gcloud pubsub, and another go application (pull.go) pulls them out. *Skills implemented include: GoLang, Gorilla Mux, MongoDB, Gcloud PubSub, Docker* 


## How to Run
First ensure start a MongoDB server and intialize the DB's
I did this through Docker:
```
$ docker pull mongo:4.0.4
$ docker run -d -p 27017-27019:27017-27019 --name mongodb mongo:4.0.4
$ docker exec -it mongodb bash
$ mongo
```

Now the mogo client is loaded, here are the commands to create your DB and add entries:
```
$ show dbs
$ use {DB_NAME}
$ db.{COLLECTION_NAME}.save({ {KEY}: "{VALUE}", {KEY}: "{VALUE}" })
$ db.{COLLECTION_NAME}.find({ {KEY}: "{VALUE}" })
```

Once the DB and collection are initialized, navigate to /KaylaAPI folder:

```
$ go run main.go 
```
Then navigate to /KaylaAPI/PULL folder:
```
$ go run pull.go 
```

## Important
Add and the env variable ```GOOGLE_APPLICATION_CREDENTIALS="/Users/kforemski/go/src/git.dev.kochava.com/KaylaAPI/key-file.json"```
with the path to your gcloud key-file.json in order for the explicit func to work
