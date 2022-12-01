# go-rabbimq-mongodb-api

**To start out containers using docker-compose**

> docker-compose up --build

1. Starting read messages:
> cd receive 
> go run receive.go

2. if  you need to send a message:
> cd send
> go run sendler.go

_After that, the sent message will be recorded in MongoDB_