# go-rabbimq-mongodb-api

**To start out containers using docker-compose**
```
docker-compose up --build
```
1. Starting read messages:
```
cd receive
go run receive.go
```
2. if  you need to send a message:
```
cd send
 go run sendler.go
```
_After that, the sent message will be recorded in MongoDB_

3. To start the API Server:
```
cd api
go run api.go
```

_To get data through the API, you need to make a post request by passing the username and password (if everything is correct, you will receive a JWT).
After that, you can make a GET request with a JWT in the header.
In response, you will receive data from mongodb_