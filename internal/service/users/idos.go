package users

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/internal/service"
)

// User represents a user in the system.
type User struct {
	ID             uuid.UUID
	Email          string
	FirstName      string
	LastName       string
	Username       string
	DOB            time.Time // format: DD-MM-YYYY
	Gender         service.Gender
	Location       string
	Bio            string
	ProfilePicture string
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UpdatedBy      uuid.UUID
}

// removePI removes personally identifiable information from the User struct.
func (u *User) removePI() {
	u.LastName = ""
	u.Email = ""
	u.PasswordHash = ""
	u.DOB = time.Time{}
}

// IsZero checks if the User struct is its zero value.
func (u User) IsZero() bool {
	return u.ID == uuid.Nil
}

// Equals checks if two User structs are equal.
func (u User) Equals(other User) bool {
	return u.ID == other.ID &&
		u.Email == other.Email &&
		u.FirstName == other.FirstName &&
		u.LastName == other.LastName &&
		u.Username == other.Username &&
		u.DOB.Equal(other.DOB) &&
		u.Gender == other.Gender &&
		u.Location == other.Location &&
		u.Bio == other.Bio &&
		u.ProfilePicture == other.ProfilePicture &&
		u.PasswordHash == other.PasswordHash
}

// Update applies the given UserUpdate to the User and returns the updated User.
func (u User) Update(updates UserUpdate) User {
	if updates.Email != nil {
		u.Email = *updates.Email
	}

	if updates.FirstName != nil {
		u.FirstName = *updates.FirstName
	}

	if updates.LastName != nil {
		u.LastName = *updates.LastName
	}

	if updates.Username != nil {
		u.Username = *updates.Username
	}

	if updates.Location != nil {
		u.Location = *updates.Location
	}

	if updates.Bio != nil {
		u.Bio = *updates.Bio
	}

	if updates.Password != nil {
		u.PasswordHash = *updates.Password
	}

	u.UpdatedAt = updates.UpdatedAt
	u.UpdatedBy = updates.UpdatedBy

	return u
}

// UserUpdate represents the fields that can be updated for a user.
type UserUpdate struct {
	Email     *string
	FirstName *string
	LastName  *string
	Username  *string
	Location  *string
	Bio       *string
	Gender    *service.Gender
	Password  *string
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
}

// IsZero checks if the UserUpdate struct has no effective fields set.
func (uu UserUpdate) IsZero() bool {
	return uu.Email == nil &&
		uu.FirstName == nil &&
		uu.LastName == nil &&
		uu.Username == nil &&
		uu.Location == nil &&
		uu.Bio == nil &&
		uu.Gender == nil &&
		uu.Password == nil
}
