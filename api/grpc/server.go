package grpc

import (
	"context"
	"log"
	"strconv"

	pb "gitlab.skypicker.com/platform/security/iam/api/grpc/v1"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

type userDataService interface {
	GetUser(string) (okta.User, error)
	AddPermissions(*okta.User, string) error
}

type Server struct {
	userService userDataService
}

func CreateServer(userServiceClient userDataService) (*Server, error) {
	return &Server{userService: userServiceClient}, nil
}

func (s *Server) User(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	log.Println("Hi")
	user, _ := s.userService.GetUser(in.Email)
	employeNumber, _ := strconv.ParseInt(user.EmployeeNumber, 10, 64)

	return &pb.UserResponse{
		EmployeeNumber: employeNumber,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Position:       user.Position,
		Department:     user.Department,
		Location:       user.Location,
		Manager:        user.Manager,
		TeamMembership: user.TeamMembership,
	}, nil
}
