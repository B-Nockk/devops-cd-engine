package strategy

import (
	"cd-engine/internal/domain"
	"errors"
)

// RollingPlanner implements Planner for rolling deployments.
type RollingPlanner struct{}

func (p *RollingPlanner) Plan(release *domain.Release, env *domain.Environment) ([]Step, error) {
	if env.Strategy != domain.StrategyRolling {
		return nil, errors.New("environment is not rolling strategy")
	}

	steps := []Step{}
	// Example: assume env has a list of server hosts
	servers := []string{"server1", "server2", "server3"} // placeholder

	for _, srv := range servers {
		// Pull Artifact
		pull := Step{
			Name:       "Pull Artifact",
			ActionType: ActionPullArtifact,
			Target:     srv,
		}
		pull.RollbackStep = &Step{
			Name:       "Remove Artifact",
			ActionType: ActionStopContainer, // simplified inverse
			Target:     srv,
		}
		steps = append(steps, pull)

		// Start Container
		start := Step{
			Name:       "Start Container",
			ActionType: ActionStartContainer,
			Target:     srv,
		}
		start.RollbackStep = &Step{
			Name:       "Stop Container",
			ActionType: ActionStopContainer,
			Target:     srv,
		}
		steps = append(steps, start)

		// Check Health
		check := Step{
			Name:       "Check Health",
			ActionType: ActionCheckHealth,
			Target:     srv,
		}
		check.RollbackStep = nil
		steps = append(steps, check)
	}

	return steps, nil
}
