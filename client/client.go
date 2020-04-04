package main

import (
	"context"
	"flag"
	"fmt"

	as "github.com/llschuster/qasmoke/appservices"
	"google.golang.org/grpc"
)

func main() {
	srvAddress := flag.String("address", "localhost:5000", "address for server format host:port")
	flag.Parse()

	con, err := grpc.Dial(*srvAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Error while running rpc client ", err)
		return
	}
	defer con.Close()
	client := as.NewQasmoskeClient(con)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var response *as.Response
	response, err = client.GetMessage(ctx, &as.Request{Message: "you gay"})
	if err != nil {
		fmt.Println("error while calling endpoint ", err)
	}
	fmt.Println(response.Message)
}
