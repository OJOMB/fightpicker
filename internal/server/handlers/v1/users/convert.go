package users

import (
	"github.com/mikhalytch/eggs/deref"
	"github.com/oapi-codegen/runtime/types"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	service "github.com/OJOMB/fightpicker/internal/service/users"
)

func userCreateRequestDTOtoIDO(ucr dtos.UserCreateReq) service.User {
	return service.User{
		Email:        string(ucr.Email),
		FirstName:    ucr.FirstName,
		LastName:     ucr.LastName,
		Username:     ucr.Username,
		DOB:          ucr.Dob.Time,
		Location:     deref.OrDefault(ucr.Location),
		Bio:          deref.OrDefault(ucr.Bio),
		PasswordHash: ucr.Password,
	}
}

func userUpdateDTOToIDO(ucr dtos.UserUpdateReq) service.UserUpdate {
	var email *string
	if ucr.Email != nil {
		e := string(*ucr.Email)
		email = &e
	}

	return service.UserUpdate{
		Email:     email,
		FirstName: ucr.FirstName,
		LastName:  ucr.LastName,
		Username:  ucr.Username,
		Location:  ucr.Location,
		Bio:       ucr.Bio,
		Password:  ucr.Password,
	}
}

// userIDOToDTO converts a service.User to a dtos.UserResponse.
// It omits sensitive fields such as Email, LastName, Password, and DOB
// if they are empty or zeroed.
func userIDOToDTO(u service.User) dtos.UserResponse {
	var email *types.Email
	if u.Email != "" {
		e := types.Email(u.Email)
		email = &e
	}

	var dob *types.Date
	if !u.DOB.IsZero() {
		dob = &types.Date{Time: u.DOB}
	}

	var first_name *string
	if u.FirstName != "" {
		first_name = &u.FirstName
	}

	var last_name *string
	if u.LastName != "" {
		last_name = &u.LastName
	}

	var location *string
	if u.Location != "" {
		location = &u.Location
	}

	var profilePicture *string
	if u.ProfilePicture != "" {
		profilePicture = &u.ProfilePicture
	}

	return dtos.UserResponse{
		Id:             u.ID,
		Email:          email,
		FirstName:      first_name,
		LastName:       last_name,
		Username:       u.Username,
		Dob:            dob,
		Location:       location,
		Bio:            u.Bio,
		ProfilePicture: profilePicture,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
