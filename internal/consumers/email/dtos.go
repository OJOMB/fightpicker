package email

type EventUserCreation struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
}
