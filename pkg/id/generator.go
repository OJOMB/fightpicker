package id

import "github.com/gofrs/uuid"

type Generator interface {
	Generate() uuid.UUID
}
