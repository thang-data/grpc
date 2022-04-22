package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"gitlab.com/dem1/dem1/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopService struct {
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopService {
	return &LaptopService{store}
}

func (server *LaptopService) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
	) (*pb.CreateLaptopResponse, error){
		laptop := req.GetLaptop()
		log.Printf("Receive a create-laptop request with id: %s", laptop.Id)
		
		if len(laptop.Id) > 0 {
			// check if it's a valid UUID
			_, err := uuid.Parse(laptop.Id)

			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "laptop Id is not a valid UUID: %v", err)
			}
		}else {
			id, err := uuid.NewRandom()
			if err != nil{
				return nil, status.Errorf(codes.Internal, "cannot generate a new laptop Id : %v", err)
			}
			laptop.Id = id.String()
		}

	// save the laptop to  store 
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal

        if errors.Is(err, ErrAlreadyExists) {
            code = codes.AlreadyExists
        }
		return nil, status.Errorf(code, "cannot save laptop to the store : %v", err)
	}
	log.Printf("save laptop with id : %v", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}