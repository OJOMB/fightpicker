package auth

type AuthTool struct {
	*BCryptPasswordTool
	*TokenGenerator
}

func NewAuthTool(bcryptCost int) *AuthTool {
	return &AuthTool{
		BCryptPasswordTool: NewBCryptPasswordTool(bcryptCost),
		TokenGenerator:     &TokenGenerator{},
	}
}
