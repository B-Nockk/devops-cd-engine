// internal/ports/outBound/deployment_store.go
package ports

import (
	domain "cd-engine/internal/domain"
	"context"
)

type DeploymentStore interface {
	Create(
		ctx context.Context,
		deployment domain.Deployment,
	) error

	UpdateStatus(
		ctx context.Context,
		id domain.ID,
		status domain.DeploymentStatus,
	) error
}
