package strategy

import (
	"cd-engine/internal/domain"
	"errors"
)

// CanaryPlanner implements Planner for canary deployments.
type CanaryPlanner struct{}

func (p *CanaryPlanner) Plan(release *domain.Release, env *domain.Environment) ([]Step, error) {
	if env.Strategy != domain.StrategyCanary {
		return nil, errors.New("environment is not canary strategy")
	}

	steps := []Step{}
	canaryTarget := "canary-slot" // placeholder

	// 1. Pull Artifact to canary
	pull := Step{
		Name:       "Pull Artifact (Canary)",
		ActionType: ActionPullArtifact,
		Target:     canaryTarget,
	}
	pull.RollbackStep = &Step{
		Name:       "Remove Artifact (Canary)",
		ActionType: ActionStopContainer,
		Target:     canaryTarget,
	}
	steps = append(steps, pull)

	// 2. Start Container on canary
	start := Step{
		Name:       "Start Container (Canary)",
		ActionType: ActionStartContainer,
		Target:     canaryTarget,
	}
	start.RollbackStep = &Step{
		Name:       "Stop Container (Canary)",
		ActionType: ActionStopContainer,
		Target:     canaryTarget,
	}
	steps = append(steps, start)

	// 3. Canary Health Check / Wait
	wait := Step{
		Name:       "Canary Validation",
		ActionType: ActionCheckHealth, // could also define ActionWait
		Target:     canaryTarget,
	}
	wait.RollbackStep = nil
	steps = append(steps, wait)

	// 4. Rollout to main targets (simplified example)
	mainTargets := []string{"main1", "main2"}
	for _, tgt := range mainTargets {
		pullMain := Step{
			Name:       "Pull Artifact (Main)",
			ActionType: ActionPullArtifact,
			Target:     tgt,
		}
		pullMain.RollbackStep = &Step{
			Name:       "Remove Artifact (Main)",
			ActionType: ActionStopContainer,
			Target:     tgt,
		}
		steps = append(steps, pullMain)

		startMain := Step{
			Name:       "Start Container (Main)",
			ActionType: ActionStartContainer,
			Target:     tgt,
		}
		startMain.RollbackStep = &Step{
			Name:       "Stop Container (Main)",
			ActionType: ActionStopContainer,
			Target:     tgt,
		}
		steps = append(steps, startMain)

		checkMain := Step{
			Name:       "Check Health (Main)",
			ActionType: ActionCheckHealth,
			Target:     tgt,
		}
		checkMain.RollbackStep = nil
		steps = append(steps, checkMain)
	}

	return steps, nil
}
