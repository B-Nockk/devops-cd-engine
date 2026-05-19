// internal/ports/outBound/executor.go
package ports

import "context"

type Executor interface {
	Execute(
		ctx context.Context,
		target string,
		command string,
	) (string, error)
}
