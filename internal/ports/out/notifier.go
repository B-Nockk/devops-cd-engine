// internal/ports/ports.go
package ports

import (
	domain "cd-engine/internal/domain"
	"context"
)

type Notifier interface {
	Notify(
		ctx context.Context,
		tenantID domain.ID,
		message string,
	) error
}
