// internal/ports/out/release_store.go
package ports

import (
	domain "cd-engine/internal/domain"
	"context"
)

type ReleaseStore interface {
	Create(
		ctx context.Context,
		release domain.Release,
	) error

	UpdateStatus(
		ctx context.Context,
		id domain.ID,
		status domain.ReleaseStatus,
	) error

	GetRelease(
		ctx context.Context,
		id domain.ID,
	) (domain.Release, error)
}
