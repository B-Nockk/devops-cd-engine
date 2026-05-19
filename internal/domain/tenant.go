package domain

import (
	"errors"
)

type NotifyChannel string

const (
	NotifyChannelSlack NotifyChannel = "slack"
	NotifyChannelEmail NotifyChannel = "email"
)

type Tenant struct {
	ID            ID
	Name          string
	Slug          string
	SSHHost       string
	SSHUser       string
	SSHKeyRef     string
	NotifyChannel NotifyChannel
	MetaData      MetaData
}

func NewTenant(
	name,
	slug,
	sshHost,
	sshUser,
	sshKeyRef string,
	notifyChannel NotifyChannel,
	note string) (*Tenant, error) {
	if name == "" || slug == "" {
		return nil, errors.New("tenant requires id, name, and slug")
	}
	return &Tenant{
		ID:            NewID(),
		Name:          name,
		Slug:          slug,
		SSHHost:       sshHost,
		SSHUser:       sshUser,
		SSHKeyRef:     sshKeyRef,
		NotifyChannel: notifyChannel,
		MetaData:      NewMetaData(note),
	}, nil
}
