package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/volchkovski/go-rest-api/pkg/swagger/server/restapi"
	"github.com/volchkovski/go-rest-api/pkg/swagger/server/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	server.Port = 8080

	api.PingAPIHandler = operations.PingAPIHandlerFunc(Ping)
	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)
	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func Ping(operations.PingAPIParams) middleware.Responder {
	return operations.NewPingAPIOK().WithPayload("OK")
}

func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload(fmt.Sprintf("Hello, %s!", user.User))
}

func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {
	if gopher.Name == "" {
		gopher.Name = "dr-who"
	}
	URL := fmt.Sprintf("https://github.com/scraly/gophers/raw/main/%s.png", gopher.Name)

	response, err := http.Get(URL)
	if err != nil {
		log.Println(err)
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}
