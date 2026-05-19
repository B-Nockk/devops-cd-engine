// internal/app/orchestrator.go
package app

import (
	"context"
	"errors"
	"time"

	domain "cd-engine/internal/domain"
	strategy "cd-engine/internal/domain/strategy"
	out "cd-engine/internal/ports/out"
)

type DeploymentOrchestrator struct {
	tenantStore   out.TenantStore
	envStore      out.EnvironmentStore
	releaseStore  out.ReleaseStore
	planner       strategy.Planner
	executor      out.Executor
	healthChecker out.HealthChecker
}

func NewDeploymentOrchestrator(
	tenantStore out.TenantStore,
	envStore out.EnvironmentStore,
	releaseStore out.ReleaseStore,
	planner strategy.Planner,
	executor out.Executor,
	healthChecker out.HealthChecker,
) *DeploymentOrchestrator {
	return &DeploymentOrchestrator{
		tenantStore:   tenantStore,
		envStore:      envStore,
		releaseStore:  releaseStore,
		planner:       planner,
		executor:      executor,
		healthChecker: healthChecker,
	}
}

func (o *DeploymentOrchestrator) RunDeploy(ctx context.Context, releaseID domain.ID) error {
	// 1. Fetch State
	release, err := o.releaseStore.GetRelease(ctx, releaseID)
	if err != nil {
		return err
	}
	env, err := o.envStore.GetEnvironment(ctx, domain.ID(release.EnvironmentID))
	if err != nil {
		return err
	}
	// tenant, err := o.tenantStore.Get(ctx, domain.ID(env.TenantID))
	// if err != nil {
	// 	return err
	// }

	// 2. Update Status -> in_progress
	release.Status = domain.ReleaseStatusInProgress
	if err := o.releaseStore.UpdateStatus(ctx, releaseID, release.Status); err != nil {
		return err
	}

	// 3. Plan
	steps, err := o.planner.Plan(&release, &env)
	if err != nil {
		return err
	}

	// 4. Rollback Stack
	var rollbackStack []strategy.Step

	// 5. Execution Loop
	for _, step := range steps {
		if step.ActionType == strategy.ActionCheckHealth {
			healthyCount := 0
			unhealthyCount := 0
			backoff := time.Second

			for {
				// Use the HealthCheck target defined in the Environment
				healthy, _ := o.healthChecker.Check(ctx, env.HealthCheck.Target)

				if healthy {
					healthyCount++
					unhealthyCount = 0
				} else {
					unhealthyCount++
					healthyCount = 0
				}

				if healthyCount >= env.HealthCheck.HealthyThreshold {
					break
				}
				if unhealthyCount >= env.HealthCheck.UnhealthyThreshold {
					return o.rollback(ctx, &release, rollbackStack, errors.New("health check failed"))
				}

				// The correct Go way to sleep while respecting context cancellation
				select {
				case <-ctx.Done():
					return o.rollback(ctx, &release, rollbackStack, ctx.Err())
				case <-time.After(backoff):
				}

				backoff *= 2
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
			}
		} else {
			// Pass the target (e.g., IP address) and the command name to the executor
			if _, err := o.executor.Execute(ctx, step.Target, string(step.ActionType)); err != nil {
				return o.rollback(ctx, &release, rollbackStack, err)
			}
		}

		if step.RollbackStep != nil {
			rollbackStack = append(rollbackStack, *step.RollbackStep)
		}
	}

	// 7. Success
	release.Status = domain.ReleaseStatusHealthy
	return o.releaseStore.UpdateStatus(ctx, releaseID, release.Status)
}

func (o *DeploymentOrchestrator) rollback(ctx context.Context, release *domain.Release, stack []strategy.Step, originalErr error) error {
	release.Status = domain.ReleaseStatusFailed
	_ = o.releaseStore.UpdateStatus(ctx, release.ID, release.Status)

	// LIFO Rollback
	for i := len(stack) - 1; i >= 0; i-- {
		step := stack[i]
		// Best-effort rollback execution
		_, _ = o.executor.Execute(context.Background(), step.Target, string(step.ActionType))
	}
	return originalErr
}
