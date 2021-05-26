package transactions

import "strings"

//go:generate stringer -type=Type
type Type int

const (
	Unknown Type = iota
	General
	Withdrawal
)

func (s Type) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Type) UnmarshalText(text []byte) error {
	*s = StatusFromText(string(text))
	return nil
}

func StatusFromText(text string) Type {
	switch strings.ToLower(text) {
	default:
		return Unknown
	case "general":
		return General
	case "withdrawal":
		return Withdrawal
	}
}
