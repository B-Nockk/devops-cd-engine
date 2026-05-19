// internal/adapters/healthcheck/checker.go
package healthcheck

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	out "cd-engine/internal/ports/out"
)

type checker struct {
	httpClient *http.Client
}

func NewChecker() out.HealthChecker {
	return &checker{
		httpClient: &http.Client{
			// Prevent infinite redirects
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				return nil
			},
		},
	}
}

// Check inspects the target and performs either HTTP or TCP health check
func (c *checker) Check(ctx context.Context, target string) (bool, string) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		// HTTP health check
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return false, fmt.Sprintf("failed to create request: %v", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return false, fmt.Sprintf("http request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return true, fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		return false, fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// TCP health check
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		return false, fmt.Sprintf("tcp connection failed: %v", err)
	}
	defer conn.Close()

	return true, "tcp connection successful"
}
