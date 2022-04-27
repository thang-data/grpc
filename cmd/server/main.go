package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"gitlab.com/dem1/dem1/pb"
	"gitlab.com/dem1/dem1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println(" --> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println(" --> stram interceptor: ", info.FullMethod)
	return handler(srv, stream)
}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")

	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)

	if err != nil {
		return err
	}
	return userStore.Save(user)
}
func main() {
	port := flag.Int("port", 0, "the server port")

	flag.Parse()

	log.Printf("Starting server on port %d", port)

	userStores := service.NewInMemoryUserStore()
	err := seedUsers(userStores)
	if err != nil {
		log.Fatal("cannot seed users")
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStores, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

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
