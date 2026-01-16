package users

import (
	"github.com/mikhalytch/eggs/deref"
	"google.golang.org/protobuf/types/known/timestamppb"

	userspb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	"github.com/OJOMB/fightpicker/internal/service"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
)

func userIDOtoDTO(u usersservice.User) *userspb.User {
	return &userspb.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		Username:  u.Username,
		Dob:       timestamppb.New(u.DOB),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Location:  u.Location,
		Bio:       u.Bio,
		Gender:    userspb.Gender(u.Gender),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func createUserRequestDTOtoIDO(req *userspb.CreateUserRequest) usersservice.User {
	return usersservice.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: req.Password,
		DOB:          req.Dob.AsTime(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Location:     deref.OrDefault(req.Location),
		Bio:          deref.OrDefault(req.Bio),
		Gender:       service.Gender(req.Gender),
	}
}

func updateUserRequestDTOtoIDO(req *userspb.UpdateUserRequest) usersservice.UserUpdate {
	return usersservice.UserUpdate{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Location:  req.Location,
		Bio:       req.Bio,
		Gender:    (*service.Gender)(req.Gender),
	}
}
