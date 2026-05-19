package domain

import (
	"errors"
	"fmt"
)

type Strategy string

const (
	StrategyBlueGreen Strategy = "blue_green"
	StrategyRolling   Strategy = "rolling"
	StrategyCanary    Strategy = "canary"
)

type Environment struct {
	ID             ID
	TenantID       string
	Name           string
	Strategy       Strategy
	HealthCheck    HealthCheckConfig
	RollbackPolicy RollbackPolicy
	MetaData       MetaData
}

func NewEnvironment(
	tenantID string,
	name string,
	st Strategy,
	hc HealthCheckConfig,
	rp RollbackPolicy,
	note string,
) (*Environment, error) {
	if tenantID == "" || name == "" {
		return nil, errors.New("environment requires tenantID & name")
	}

	if st != StrategyBlueGreen && st != StrategyRolling && st != StrategyCanary {
		return nil, fmt.Errorf("invalid strategy: %s", st)
	}
	return &Environment{
		ID:             NewID(), // SNID/KSUID
		TenantID:       tenantID,
		Name:           name,
		Strategy:       st,
		HealthCheck:    hc,
		RollbackPolicy: rp,
		MetaData:       NewMetaData(note),
	}, nil
}
