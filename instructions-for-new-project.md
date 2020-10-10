# Overview
This file contains instructions for a sample project with login functionality. 
It contains different layers as described below. 

It uses these dependencies
* github.com/devingen/api-core
* go.mongodb.org/mongo-driver
* github.com/gorilla/mux
* github.com/aws/aws-lambda-go

**Service**

Services are responsible for data storage and retrieval. The service we have
* finds user with email
* creates a new user

**Controller**

Controllers are responsible for handling incoming requests and returning responses. 
They contain all the logic. The controller we have uses the service to find and create users. It
* registers a new user
  * find user with email (service)
    * returns an error if user with the email exists
  * create a new user (service)
  * returns the details of the new user

### Basic files

Create a new `go.mod` file with content

```
module github.com/devingen/PROJECT_NAME

go 1.12

require (
    github.com/devingen/api-core v0.0.8
    go.mongodb.org/mongo-driver v1.3.2
)

```

### Models

This API will contain a User model. So create a folder `model` and put this `user.go` file in it.

```go
package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {  
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FirstName string             `json:"firstName"`
    LastName  string             `json:"lastName"`
    Email     string             `json:"email"`
}
```

### Services

Create a folder `service` in root for our service related files. 

#### Definition 

Create `service/definition.go` with the content below.

```go
package service

type DevingenService interface {
    CreateUser(base, firstName, lastName, email string) (*model.User, error)
    FindUserUserWithEmail(base, email string) (*model.User, error)
}
```

#### Implementation

Create `service/database-service` folder and these files in it.


```go
// database-service.go

package database_service

// DatabaseService implements DataService interface with database connection
type DatabaseService struct {
    Database *database.Database
}
```

```go
// find-user-with-email.go

package database_service

func (service DatabaseService) FindUserUserWithEmail(base, email string) (*model.User, error) {
    fmt.Println("service.FindUserUserWithEmail", base, email)
    return nil, nil
}
```

```go
// create-user.go

package database_service

func (service DatabaseService) CreateUser(base, firstName, lastName, email string) (*model.User, error) {
    fmt.Println("service.CreateUser", base, firstName, lastName, email)
    return &model.User{
        ID:        primitive.NewObjectID(),
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
    }, nil
}
```

#### Init helper

Create `service/generator.go` that will be used to create this service.

```go
package service

// NewDatabaseService generates new DatabaseService
func NewDatabaseService(database *database.Database) *database_service.DatabaseService {
    return &database_service.DatabaseService{
        Database: database,
    }
}

```

### Data Transfer Objects

We'll need register related models for request and response. 
Create a folder `dto` and a `register.go` file with content

```go
package dto

type RegisterWithEmailRequest struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

type RegisterWithEmailResponse struct {
    UserID string `json:"userId"`
    JWT    string `json:"jwt"`
}
```

### Controllers

First, create folder `controller` in root. 

#### Definition

Create `controller/definition.go` file with content below. This file defines the input/output of our functions.

```
package controller

type IServiceController interface {
	RegisterWithEmail(base string, request *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error)
}
```

#### Implementation

Create `controller/service-controller` folder and these files in it.

```go
// service-controller.go

package service_controller

// ServiceController implements IServiceController interface by using DevingenService
type ServiceController struct {
    Service service.IServiceController
}
```

```go
// register.go

package service_controller

func (controller ServiceController) RegisterWithEmail(base string, request *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error) {
    userWithSameEmail, err := controller.Service.FindUserUserWithEmail(base, email)
    if err != nil {
        return nil, err
    }
    
    if userWithSameEmail != nil {
        return nil, coremodel.NewStatusError(http.StatusConflict)
    }
    
    user, err := controller.Service.CreateUser(
        base,
        request.FirstName,
        request.LastName,
        email,
    )
    if err != nil {
        return nil, err
    }

}
```

#### Init helper

Create `controller/generator.go` that will be used to create this controller.

```go
package controller

// NewServiceController generates new ServiceController
func NewServiceController(service service.IDevingenService) *service_controller.ServiceController {
    return &service_controller.ServiceController{
        Service: service,
    }
}
```

### Web Server

Web server is the server that'll run on port `8080` and listen for the requests.
Create `server` folder that'll contain the web server related files.
 
#### Request Handler

Create `server/handler` folder and these files in it.

```go
// definition.go

package handler

type ServerHandler struct {
    Controller controller.IServiceController
    Router     *mux.Router
}
```

```go
// register.go

package handler

func (handler ServerHandler) register(w http.ResponseWriter, r *http.Request) {
    pathVariables := mux.Vars(r)
    base := pathVariables["base"]
    
    var body dto.RegisterWithEmailRequest
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    result, err := handler.Controller.RegisterWithEmail(base, &body)
    response, err := util.BuildResponse(http.StatusCreated, result, err)
    server.ReturnResponse(w, response, err)
}
```

```go
// generator.go

package handler

func NewHttpServiceHandler(controller controller.IServiceController) ServerHandler {
    handler := ServerHandler{Controller: controller}
    
    handler.Router = mux.NewRouter()
    handler.Router.HandleFunc("/{base}/register", handler.register).Methods(http.MethodPost)
    
    return handler
}
```
 
 #### Server
 
 Create `main.go` in `server` folder.

```go
package main

import (
    cm "github.com/devingen/api-core/server"
    "log"
    "net/http"
)

// Runs the server that contains all the services
func main() {
    db, err := database.NewDatabase()
    if err != nil {
        log.Fatalf("Database connection failed %s", err.Error())
    }
    
    databaseService := service.NewDatabaseService(db)
    serviceController := controller.NewServiceController(databaseService)
    
    // create a Service Handler that uses Database AtamaService
    h := handler.NewHttpServiceHandler(serviceController)
    
    http.Handle("/", &cm.CORSRouterDecorator{R: h.Router})
    err = http.ListenAndServe(":8080", &cm.CORSRouterDecorator{R: h.Router})
    if err != nil {
        log.Fatalf("Listen and serve failed %s", err.Error())
    }
}

```

#### Running server

```shell script
export MONGO_ADDRESS=localhost
go run server/main.go
```

To test;

```shell script
curl -d '{"firstName":"Emir", "lastName":"Luleci", "email":"emir@devingen.io", "password":"selam"}' -X POST http://localhost:8080/devingen-dev/register
```

### AWS Lambda

The project also includes AWS Lambda configuration to deploy individual functions to AWS Lambda.

Create `aws` folder in root and create `database.go` file in it.

```go
package aws

var db *database.Database

func GetDatabase() *database.Database {
    if db == nil {
        var err error
        db, err = database.NewDatabase()
        if err != nil {
            log.Fatalf("Database connection failed %s", err.Error())
        }
    } else if !db.IsConnected() {
        err := db.ConnectWithEnvironment()
        if err != nil {
            log.Fatalf("Database connection failed %s", err.Error())
        }
    } else {
        log.Println("Database connection exists")
    }
    return db
}
```

Create `aws/register` folder and create `main.go` file in it.

```go
package main

import (
    coreaws "github.com/devingen/api-core/aws"
)

func main() {
    db := aws.GetDatabase()
    databaseService := service.NewDatabaseService(db)
    serviceController := controller.NewServiceController(databaseService)
    
    lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    
        var body dto.RegisterWithEmailRequest
        err := json.Unmarshal([]byte(req.Body), &body)
        if err != nil {
            return events.APIGatewayProxyResponse{}, err
        }
    
        base := req.PathParameters["base"]
    
        result, err := serviceController.RegisterWithEmail(base, &body)
        response, err := util.BuildResponse(http.StatusCreated, result, err)
        return coreaws.AdaptResponse(response, err)
    })
}
```

#### Serverless configuration

Create `serverless.yml` in the root folder.

```yaml
org: devingen
app: devingen-io
service: devingen-api

frameworkVersion: '>=1.28.0 <2.0.0'

provider:
  name: aws
  runtime: go1.x
  environment:
    MONGO_ADDRESS: ${param:MONGO_ADDRESS}
    MONGO_USERNAME: ${param:MONGO_USERNAME}
    MONGO_PASSWORD: ${param:MONGO_PASSWORD}

package:
  exclude:
    - ./**
  include:
    - ./bin/**

functions:

  register:
    handler: bin/register
    events:
      - http:
          path: /{base}/register
          method: post
          cors: true
          request:
            parameters:
              paths:
                base: true

```

#### Build configuration

Create `Makefile` in root folder.

```makefile
.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/register aws/register/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-devingen: clean build
	serverless deploy --stage devingen --region eu-central-1 --verbose

teardown-devingen: clean
	serverless remove --stage devingen --region eu-central-1 --verbose

deploy-devingen-dev: clean build
	serverless deploy --stage devingen-dev --region eu-central-1 --verbose

teardown-devingen-dev: clean
	serverless remove --stage devingen-dev --region eu-central-1 --verbose

```

#### Deploying AWS Lambda functions

First, configure the [serverless.com](https://www.serverless.com/) account then run `deploy-devingen-dev` command.

This commands executes the command in `Makefile` which clears the previous builds,
generates executables and deploys the AWS Functions through Serverless Framework.

```shell script
make deploy-devingen-dev
```

To test, replace the url with the one printed in the output and execute this

```shell script
curl -d '{"firstName":"Emir", "lastName":"Luleci", "email":"emir@devingen.io", "password":"selam"}' -X POST https://lg3cpwzyej.execute-api.ca-central-1.amazonaws.com/dev/devingen-dev/register
```
