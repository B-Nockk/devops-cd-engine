package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // CGO-free SQLite driver

	"cd-engine/internal/domain"
)

// Store implements TenantStore, EnvironmentStore, ReleaseStore, DeploymentStore
type Store struct {
	db *sql.DB
}

// NewStore initializes SQLite DB and runs migrations.
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	migrations := `
    CREATE TABLE IF NOT EXISTS tenants (
        id TEXT PRIMARY KEY,
        name TEXT,
        slug TEXT,
        ssh_host TEXT,
        ssh_user TEXT,
        ssh_key_ref TEXT,
        notify_channel TEXT,
        metadata TEXT
    );
    CREATE TABLE IF NOT EXISTS environments (
        id TEXT PRIMARY KEY,
        tenant_id TEXT,
        name TEXT,
        strategy TEXT,
        healthcheck TEXT,
        rollback_policy TEXT,
        metadata TEXT
    );
    CREATE TABLE IF NOT EXISTS releases (
        id TEXT PRIMARY KEY,
        environment_id TEXT,
        artifact TEXT,
        git_tag TEXT,
        initiated_by TEXT,
        status TEXT,
        strategy_used TEXT,
        started_at TEXT,
        completed_at TEXT,
        release_notes TEXT,
        metadata TEXT
    );
    CREATE TABLE IF NOT EXISTS deployments (
        id TEXT PRIMARY KEY,
        release_id TEXT,
        slot TEXT,
        server_host TEXT,
        status TEXT,
        started_at TEXT,
        initiated_by TEXT,
        metadata TEXT
    );
    `
	if _, err := db.Exec(migrations); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Get(
	ctx context.Context,
	id domain.ID,
) (domain.Tenant, error) {
	var t domain.Tenant
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug, ssh_host, ssh_user, ssh_key_ref, notify_channel, metadata FROM tenants WHERE id = ?`,
		string(id))
	err := row.Scan(&t.ID, &t.Name, &t.Slug, &t.SSHHost, &t.SSHUser, &t.SSHKeyRef, &t.NotifyChannel, &t.MetaData)
	return t, err
}

// --- EnvironmentStore ---
func (s *Store) GetEnvironment(ctx context.Context, id domain.ID) (domain.Environment, error) {
	var e domain.Environment
	// Use []byte to safely extract strings/json from SQLite
	var hcJSON, rpJSON, metaJSON []byte

	row := s.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, strategy, healthcheck, rollback_policy, metadata FROM environments WHERE id = ?`,
		string(id))

	// Scan into the temporary byte slices
	err := row.Scan(&e.ID, &e.TenantID, &e.Name, &e.Strategy, &hcJSON, &rpJSON, &metaJSON)
	if err != nil {
		return e, err
	}

	// Safely unmarshal the JSON strings into the Domain Structs
	if len(hcJSON) > 0 {
		_ = json.Unmarshal(hcJSON, &e.HealthCheck)
	}
	if len(rpJSON) > 0 {
		_ = json.Unmarshal(rpJSON, &e.RollbackPolicy)
	}

	return e, nil
}

// --- ReleaseStore ---
func (s *Store) Create(
	ctx context.Context,
	r domain.Release,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO releases (
			id, environment_id, artifact,
			git_tag, initiated_by, status,
			strategy_used, started_at, completed_at,
			release_notes, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		string(r.ID),
		string(r.EnvironmentID),
		r.Artifact,
		r.GitTag,
		r.InitiatedBy,
		r.Status,
		r.StrategyUsed,
		r.StartedAt,
		r.CompletedAt,
		r.ReleaseNotes,
		r.MetaData,
	)
	return err
}

func (s *Store) UpdateStatus(
	ctx context.Context,
	id domain.ID,
	status domain.ReleaseStatus,
) error {
	_, err := s.db.ExecContext(ctx, `UPDATE releases SET status = ? WHERE id = ?`, status, string(id))
	return err
}

func (s *Store) GetRelease(ctx context.Context, id domain.ID) (domain.Release, error) {
	var r domain.Release

	// Temporary variables for complex types and nullables
	var initJSON, metaJSON []byte
	var startedAtStr string
	var completedAtStr sql.NullString // Protects against NULL in the database

	row := s.db.QueryRowContext(
		ctx,
		`SELECT id, environment_id, artifact,
		git_tag, initiated_by, status,
		strategy_used, started_at, completed_at,
		release_notes, metadata FROM releases WHERE id = ?`,
		string(id),
	)
	err := row.Scan(
		&r.ID,
		&r.EnvironmentID,
		&r.Artifact,
		&r.GitTag,
		&initJSON, // Scan into byte slice
		&r.Status,
		&r.StrategyUsed,
		&startedAtStr,   // Scan into string
		&completedAtStr, // Scan into NullString
		&r.ReleaseNotes,
		&metaJSON, // Scan into byte slice
	)
	if err != nil {
		return r, err
	}

	// Safely unmarshal JSON into Domain Structs
	if len(initJSON) > 0 {
		_ = json.Unmarshal(initJSON, &r.InitiatedBy)
	}
	if len(metaJSON) > 0 {
		_ = json.Unmarshal(metaJSON, &r.MetaData)
	}

	// Parse timestamps
	if startedAtStr != "" {
		r.StartedAt, _ = time.Parse(time.RFC3339, startedAtStr)
	}
	if completedAtStr.Valid && completedAtStr.String != "" {
		t, _ := time.Parse(time.RFC3339, completedAtStr.String)
		r.CompletedAt = &t
	}

	return r, nil
}

// --- DeploymentStore ---
func (s *Store) CreateDeployment(
	ctx context.Context,
	d domain.Deployment,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO deployments (
			id, release_id, slot,
			server_host, status, started_at,
			initiated_by, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		string(d.ID),
		string(d.ReleaseID),
		d.Slot,
		d.ServerHost,
		d.Status,
		d.StartedAt,
		d.InitiatedBy,
		d.MetaData,
	)
	return err
}

func (s *Store) UpdateStatusDeployment(
	ctx context.Context,
	id domain.ID,
	status domain.DeploymentStatus,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE deployments SET status = ? WHERE id = ?`,
		status,
		string(id),
	)
	return err
}

// List returns all tenants (Stubbed for interface compliance)
func (s *Store) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	// TODO: Implement full SELECT query when building the HTTP API
	return []domain.Tenant{}, nil
}

// ListEnvironments returns all environments for a tenant (Stubbed for interface compliance)
func (s *Store) ListEnvironments(ctx context.Context, tenantID domain.ID) ([]domain.Environment, error) {
	// TODO: Implement full SELECT query when building the HTTP API
	return []domain.Environment{}, nil
}
