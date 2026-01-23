package service

type Gender int32

const (
	GenderOther Gender = iota
	GenderMale
	GenderFemale
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	default:
		return "other"
	}
}

func GenderFromString(s string) Gender {
	switch s {
	case "male":
		return GenderMale
	case "female":
		return GenderFemale
	default:
		return GenderOther
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (g *Gender) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"male"`:
		*g = GenderMale
	case `"female"`:
		*g = GenderFemale
	default:
		*g = GenderOther
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (g Gender) MarshalJSON() ([]byte, error) {
	return []byte(`"` + g.String() + `"`), nil
}
