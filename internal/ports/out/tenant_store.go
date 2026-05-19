// internal/ports/outBound/tenant_store.go
package ports

import (
	domain "cd-engine/internal/domain"
	"context"
)

type TenantStore interface {
	Get(
		ctx context.Context,
		id domain.ID,
	) (domain.Tenant, error)

	ListTenants(ctx context.Context) ([]domain.Tenant, error)
}
