package grpc

import (
	"context"
	"log"

	pb "gitlab.skypicker.com/platform/security/iam/api/grpc/v1"
)

type Server struct {
}

func (s *Server) User(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	log.Println("Hi")

	return &pb.UserResponse{
		EmployeeNumber: 1234,
		FirstName:      "Test",
		LastName:       "Tester",
		Position:       "QA Tester",
		Department:     "tst",
		Location:       "Hell",
		Manager:        "Satan",
		TeamMembership: []string{"test", "hi"},
	}, nil
}
