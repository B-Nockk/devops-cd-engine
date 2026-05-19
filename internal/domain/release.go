package domain

import (
	"errors"
	"fmt"
	"time"
)

type ReleaseStatus string

const (
	ReleasePending    ReleaseStatus = "pending"
	ReleaseInProgress ReleaseStatus = "in_progress"
	ReleaseHealthy    ReleaseStatus = "healthy"
	ReleaseFailed     ReleaseStatus = "failed"
	ReleaseRolledBack ReleaseStatus = "rolled_back"
)

type Release struct {
	ID            ID
	EnvironmentID ID
	Artifact      string
	GitTag        string
	InitiatedBy   Initiator
	Status        ReleaseStatus
	StrategyUsed  Strategy
	StartedAt     time.Time
	CompletedAt   *time.Time
	ReleaseNotes  string
	MetaData      MetaData
}

func NewRelease(
	envID ID,
	artifact string,
	gitTag string,
	initiator Initiator,
	strategy Strategy,
	releaseNotes string,
	note string,
) (*Release, error) {
	if envID.IsEmpty() || artifact == "" {
		return nil, errors.New("release requires environmentID and artifact")
	}
	if strategy != StrategyBlueGreen && strategy != StrategyRolling && strategy != StrategyCanary {
		return nil, fmt.Errorf("invalid strategy: %s", strategy)
	}

	return &Release{
		ID:            NewID(),
		EnvironmentID: envID,
		Artifact:      artifact,
		GitTag:        gitTag,
		InitiatedBy:   initiator,
		Status:        ReleasePending,
		StrategyUsed:  strategy,
		StartedAt:     time.Now(),
		ReleaseNotes:  releaseNotes,
		MetaData:      NewMetaData(note),
	}, nil
}
