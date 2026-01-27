package service

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

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

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeName returns a normalized string suitable for
// heuristic matching (NOT identity).
//
// Example:
//
//	"José  St.-Pierre" -> "jose st pierre"
func NormalizeName(s string) string {
	if s == "" {
		return ""
	}

	// 1. Unicode normalization (decompose accents)
	//    José → José (e + combining accent)
	t := norm.NFD.String(s)

	// 2. Strip diacritical marks
	t = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, t)

	// 3. Locale-aware lowercasing
	//    Use undetermined locale to avoid Turkish İ edge cases
	t = cases.Lower(language.Und).String(t)

	// 4. Replace punctuation with spaces
	t = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		case unicode.IsSpace(r):
			return ' '
		default:
			return ' '
		}
	}, t)

	// 5. Collapse whitespace
	t = strings.Join(strings.Fields(t), " ")

	return t
}

// NormalizeUsername converts arbitrary user input into a
// deterministic, ASCII-safe username string.
func NormalizeUsername(s string) string {
	if s == "" {
		return ""
	}

	// 1. Unicode normalize + strip accents
	t := norm.NFD.String(s)
	t = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, t)

	// 2. Lowercase (locale-independent)
	t = cases.Lower(language.Und).String(t)

	// 3. Convert to allowed charset
	var b strings.Builder
	b.Grow(len(t))

	lastWasSep := false

	for _, r := range t {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastWasSep = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastWasSep = false
		case r == '.' || r == '-' || r == '_':
			if !lastWasSep {
				b.WriteRune(r)
				lastWasSep = true
			}
		case unicode.IsSpace(r):
			if !lastWasSep {
				b.WriteRune('.')
				lastWasSep = true
			}
		default:
			// drop everything else
		}
	}

	// 4. Trim separators
	out := strings.Trim(b.String(), "._-")

	return out
}

func NormalizeString(s string) string {
	return strings.TrimSpace(s)
}
