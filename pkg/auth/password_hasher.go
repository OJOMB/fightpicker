package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasherVerifier interface {
	HashPassword(password string) (string, error)
	Verify(hashedPassword, password string) (bool, error)
}

type BCryptPasswordTool struct {
	cost int
}

func NewBCryptPasswordTool(cost int) *BCryptPasswordTool {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BCryptPasswordTool{
		cost: cost,
	}
}

func (ph *BCryptPasswordTool) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (ph *BCryptPasswordTool) Verify(hashedPassword, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
