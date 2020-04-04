package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	as "github.com/llschuster/qasmoke/appservices"
	"google.golang.org/grpc"
)

type servicesObject struct {
	Message *as.Response
	Request *as.Request
}

// the parameter before function_name indicates that the function will be part of the type servicesObject
func (s *servicesObject) GetMessage(context context.Context, request *as.Request) (*as.Response, error) {
	if request != nil {
		fmt.Println("executing request ", request.Message)
		return &as.Response{Message: "no, you gay"}, nil
	}
	return &as.Response{Message: "u"}, nil
}

func main() {
	port := flag.Int("port", 5000, "server port")
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panic(err)
	}
	log.Println(fmt.Sprintf("grpc server running on port %v", *port))

	grpcServer := grpc.NewServer()
	as.RegisterQasmoskeServer(grpcServer, &servicesObject{})
	grpcServer.Serve(listen)

}
