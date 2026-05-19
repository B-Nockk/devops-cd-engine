package strategy

import (
	"cd-engine/internal/domain"
	"errors"
)

// Planner defines strategy planning interface.
type Planner interface {
	Plan(release *domain.Release, env *domain.Environment) ([]Step, error)
}

// BlueGreenPlanner implements Planner for blue/green strategy.
type BlueGreenPlanner struct{}

func (p *BlueGreenPlanner) Plan(release *domain.Release, env *domain.Environment) ([]Step, error) {
	if env.Strategy != domain.StrategyBlueGreen {
		return nil, errors.New("environment is not blue/green strategy")
	}

	steps := []Step{}

	inactiveSlot := "green" // placeholder: logic to detect active slot
	activeSlot := "blue"

	// 1. Pull Artifact
	pull := Step{
		Name:       "Pull Artifact",
		ActionType: ActionPullArtifact,
		Target:     inactiveSlot,
	}
	pull.RollbackStep = &Step{
		Name:       "Remove Artifact",
		ActionType: ActionStopContainer, // simplified inverse
		Target:     inactiveSlot,
	}
	steps = append(steps, pull)

	// 2. Start Container
	start := Step{
		Name:       "Start Container",
		ActionType: ActionStartContainer,
		Target:     inactiveSlot,
	}
	start.RollbackStep = &Step{
		Name:       "Stop Container",
		ActionType: ActionStopContainer,
		Target:     inactiveSlot,
	}
	steps = append(steps, start)

	// 3. Check Health
	check := Step{
		Name:       "Check Health",
		ActionType: ActionCheckHealth,
		Target:     inactiveSlot,
	}
	check.RollbackStep = nil // read-only
	steps = append(steps, check)

	// 4. Switch Traffic
	switchTraffic := Step{
		Name:       "Switch Traffic",
		ActionType: ActionSwitchTraffic,
		Target:     inactiveSlot,
	}
	switchTraffic.RollbackStep = &Step{
		Name:       "Switch Traffic Back",
		ActionType: ActionSwitchTraffic,
		Target:     activeSlot,
	}
	steps = append(steps, switchTraffic)

	// 5. Stop Old Container
	stopOld := Step{
		Name:       "Stop Old Container",
		ActionType: ActionStopContainer,
		Target:     activeSlot,
	}
	stopOld.RollbackStep = &Step{
		Name:       "Restart Old Container",
		ActionType: ActionStartContainer,
		Target:     activeSlot,
	}
	steps = append(steps, stopOld)

	return steps, nil
}
