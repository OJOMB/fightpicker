package users

import (
	"github.com/pkg/errors"
)

// validateCreationReq validates the fields in a User struct for creation.
func (svc *Service) validateCreationReq(u *User) error {
	// check all required fields
	if u.Email == "" {
		return errors.Wrap(ErrMissingParameter, "email")
	}

	if u.FirstName == "" {
		return errors.Wrap(ErrMissingParameter, "first_name")
	}

	if u.LastName == "" {
		return errors.Wrap(ErrMissingParameter, "last_name")
	}

	if u.Username == "" {
		return errors.Wrap(ErrMissingParameter, "username")
	}

	if u.DOB.IsZero() {
		return errors.Wrap(ErrMissingParameter, "dob")
	}

	if u.PasswordHash == "" {
		return errors.Wrap(ErrMissingParameter, "password")
	}

	// validate email format
	if !svc.regexEmail.MatchString(u.Email) {
		return errors.Wrap(ErrInvalidParameter, "email")
	}

	return nil
}

// validateUpdateReq validates the fields in a UserUpdate struct.
func (svc *Service) validateUpdateReq(u UserUpdate) error {
	if u.IsZero() {
		return errors.Wrap(ErrMissingParameter, "user")
	}

	if u.Email != nil && !svc.regexEmail.MatchString(*u.Email) {
		return errors.Wrap(ErrInvalidParameter, "email")
	}

	if u.Password != nil && *u.Password == "" {
		// TODO: basic password validation
		return errors.Wrap(ErrInvalidParameter, "password")
	}

	if u.FirstName != nil && *u.FirstName == "" {
		return errors.Wrap(ErrInvalidParameter, "first_name")
	}

	if u.LastName != nil && *u.LastName == "" {
		return errors.Wrap(ErrInvalidParameter, "last_name")
	}

	if u.Username != nil && *u.Username == "" {
		return errors.Wrap(ErrInvalidParameter, "username")
	}

	return nil
}
