package domain

import (
	"errors"
	"fmt"
)

type HealthCheckType string

const (
	HealthCheckHTTP    HealthCheckType = "http"
	HealthCheckTCP     HealthCheckType = "tcp"
	HealthCheckCommand HealthCheckType = "command"
)

type HealthCheckConfig struct {
	Type               HealthCheckType
	Target             string
	IntervalSeconds    int
	TimeoutSeconds     int
	HealthyThreshold   int
	UnhealthyThreshold int
}

func NewHealthCheckConfig(t HealthCheckType, target string, interval, timeout, healthy, unhealthy int) (HealthCheckConfig, error) {
	if target == "" {
		return HealthCheckConfig{}, errors.New("health check requires target")
	}
	if t != HealthCheckHTTP && t != HealthCheckTCP && t != HealthCheckCommand {
		return HealthCheckConfig{}, fmt.Errorf("invalid health check type: %s", t)
	}
	return HealthCheckConfig{
		Type:               t,
		Target:             target,
		IntervalSeconds:    interval,
		TimeoutSeconds:     timeout,
		HealthyThreshold:   healthy,
		UnhealthyThreshold: unhealthy,
	}, nil
}
