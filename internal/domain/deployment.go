package domain

import (
	"errors"
	"time"
)

type DeploymentStatus string

const (
	DeploymentStarting  DeploymentStatus = "starting"
	DeploymentRunning   DeploymentStatus = "running"
	DeploymentHealthy   DeploymentStatus = "healthy"
	DeploymentUnhealthy DeploymentStatus = "unhealthy"
	DeploymentStopped   DeploymentStatus = "stopped"
)

type Deployment struct {
	ID          ID
	ReleaseID   ID
	Slot        string
	ServerHost  string
	Status      DeploymentStatus
	StartedAt   time.Time
	InitiatedBy Initiator
	MetaData    MetaData
}

func NewDeployment(
	releaseID ID,
	slot string,
	serverHost string,
	initiator Initiator,
	note string,
) (*Deployment, error) {
	if releaseID.IsEmpty() || slot == "" {
		return nil, errors.New("deployment requires releaseID and slot")
	}

	return &Deployment{
		ID:          NewID(),
		ReleaseID:   releaseID,
		Slot:        slot,
		ServerHost:  serverHost,
		Status:      DeploymentStarting,
		StartedAt:   time.Now(),
		InitiatedBy: initiator,
		MetaData:    NewMetaData(note),
	}, nil
}
