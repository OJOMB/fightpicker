package id

import "github.com/gofrs/uuid/v5"

type UUIDV7Generator struct{}

func NewUUIDV7Generator() *UUIDV7Generator {
	return &UUIDV7Generator{}
}

func (g *UUIDV7Generator) Generate() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
