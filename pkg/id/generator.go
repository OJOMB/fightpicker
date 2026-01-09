package id

import "github.com/gofrs/uuid/v5"

type Generator interface {
	Generate() uuid.UUID
}
