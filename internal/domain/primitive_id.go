// internal/domain/primitive_id.go
package domain

import (
	"errors"

	"github.com/segmentio/ksuid"
)

// ID is a strongly typed identifier based on KSUID.
type ID string

// NewID generates a new KSUID-based ID.
func NewID() ID {
	return ID(ksuid.New().String())
}

// IDFromString recreates an ID from a string, with basic validation.
func IDFromString(s string) (ID, error) {
	if s == "" {
		return "", errors.New("ID cannot be empty")
	}
	// Validate format using ksuid.Parse
	if _, err := ksuid.Parse(s); err != nil {
		return "", errors.New("invalid ID format")
	}
	return ID(s), nil
}

// IsEmpty checks if the ID is unset.
func (id ID) IsEmpty() bool {
	return string(id) == ""
}

// MarshalJSON ensures IDs serialize as strings.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(id) + `"`), nil
}

// UnmarshalJSON parses IDs from JSON strings.
func (id *ID) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	parsed, err := IDFromString(str)
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}
