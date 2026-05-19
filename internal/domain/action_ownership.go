// internal/domain/action_ownership.go
package domain

type InitiatorType string

const (
	InitiatorHuman  InitiatorType = "human"
	InitiatorSystem InitiatorType = "system"
)

type Initiator struct {
	Type   InitiatorType
	UserID string // could be human user ID or system ID
	Name   string // optional display name
}
