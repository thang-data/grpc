package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"gitlab.com/dem1/dem1/pb"
	"gitlab.com/dem1/dem1/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")

	flag.Parse()

	log.Printf("Starting server on port %d", port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore,imageStore, ratingStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Print("cannot connect to server: ", err)
	}

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Print("cannot connect to server: ", err)
	}
}
