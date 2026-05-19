// internal/ports/outBound/health_checker.go
package ports

import "context"

type HealthChecker interface {
	Check(
		ctx context.Context,
		target string,
	) (bool, string)
}
