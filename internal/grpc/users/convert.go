package users

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	typespb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/common/types"
	userspb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
)

func toProtoUser(u usersservice.User) *userspb.User {
	return &userspb.User{
		Id: &typespb.UUID7{
			Value: u.ID.String(),
		},
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
