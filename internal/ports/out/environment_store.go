// internal/ports/out/environment_store.go
package ports

import (
	domain "cd-engine/internal/domain"
	"context"
)

type EnvironmentStore interface {
	GetEnvironment(
		ctx context.Context,
		id domain.ID,
	) (domain.Environment, error)

	List(
		ctx context.Context,
		tenantID domain.ID,
	) ([]domain.Environment, error)
}
