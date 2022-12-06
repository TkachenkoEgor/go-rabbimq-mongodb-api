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
go run api.go jwt.go
```

_To get data through the API, you need to make a post request by passing the username and password (if everything is correct, you will receive a JWT).
After that, you can make a GET request with a JWT in the header.
In response, you will receive data from mongodb_


4. To send a HTTP request with a jwt in the header using the GO code
```
cd api-request
go run request.go
```
# IMPORTANT 

it is necessary to correctly create the request body (url 2)
specify at least one date and the required collection.
otherwise, you will see error (Invalid Namespace)
or an empty structure