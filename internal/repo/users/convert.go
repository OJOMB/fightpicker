package users

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/internal/service"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

func UserIDOtoCreateUserParamsDBO(ido usersservice.User) postgres.CreateUserParams {
	return postgres.CreateUserParams{
		ID:        ido.ID,
		Email:     ido.Email,
		FirstName: ido.FirstName,
		LastName:  ido.LastName,
		Username:  ido.Username,
		Dob: pgtype.Date{
			Time:  ido.DOB.UTC(),
			Valid: true,
		},
		Gender:         postgres.Gender(ido.Gender.String()),
		Location:       pgtype.Text{String: ido.Location, Valid: true},
		Bio:            pgtype.Text{String: ido.Bio, Valid: true},
		ProfilePicture: pgtype.Text{String: ido.ProfilePicture, Valid: true},
		PasswordHash:   ido.PasswordHash,
		CreatedAt:      pgtype.Timestamptz{Time: ido.CreatedAt.UTC(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: ido.UpdatedAt.UTC(), Valid: true},
	}
}

func userByEmailDBOToUserIDO(dbo postgres.GetUserByEmailRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		PasswordHash:   dbo.PasswordHash,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}

func userByUsernameDBOToUserIDO(dbo postgres.GetUserByUsernameRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		PasswordHash:   dbo.PasswordHash,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}

func userByIDDBOToUserIDO(dbo postgres.GetUserByIDRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		PasswordHash:   dbo.PasswordHash,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}

func UserIDOtoUpdateUserParamsDBO(dbo usersservice.User) postgres.UpdateUserByIDParams {
	return postgres.UpdateUserByIDParams{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		Dob:            pgtype.Date{Time: dbo.DOB.UTC(), Valid: true},
		Gender:         postgres.Gender(dbo.Gender.String()),
		Location:       pgtype.Text{String: dbo.Location, Valid: true},
		Bio:            pgtype.Text{String: dbo.Bio, Valid: true},
		ProfilePicture: pgtype.Text{String: dbo.ProfilePicture, Valid: true},
		PasswordHash:   dbo.PasswordHash,
		UpdatedAt:      pgtype.Timestamptz{Time: dbo.UpdatedAt.UTC(), Valid: true},
	}
}

func listUsersRowDBOToUserIDO(dbo postgres.ListUsersRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}

func listFollowersRowDBOToUserIDO(dbo postgres.ListFollowersRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}

func listFolloweesRowDBOToUserIDO(dbo postgres.ListFolloweesRow) usersservice.User {
	return usersservice.User{
		ID:             dbo.ID,
		Email:          dbo.Email,
		FirstName:      dbo.FirstName,
		LastName:       dbo.LastName,
		Username:       dbo.Username,
		DOB:            dbo.Dob.Time.UTC(),
		Gender:         service.GenderFromString(string(dbo.Gender)),
		Location:       dbo.Location.String,
		Bio:            dbo.Bio.String,
		ProfilePicture: dbo.ProfilePicture.String,
		CreatedAt:      dbo.CreatedAt.Time.UTC(),
		UpdatedAt:      dbo.UpdatedAt.Time.UTC(),
		UpdatedBy:      dbo.UpdatedBy.Bytes,
	}
}
