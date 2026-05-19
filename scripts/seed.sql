-- scripts/seed.sql
INSERT OR IGNORE INTO tenants (id, name, slug, ssh_host, ssh_user, ssh_key_ref, notify_channel, metadata)
VALUES ('tenant-1', 'Acme Corp', 'acme', '127.0.0.1', 'deploy', 'acme-key', 'slack', '{}');

-- INSERT OR IGNORE INTO environments (id, tenant_id, name, strategy, healthcheck, rollback_policy, metadata)
-- VALUES (
--     'env-1', 'tenant-1', 'production', 'blue_green',
--     '{"Target": "http://localhost:8080/health", "HealthyThreshold": 1, "UnhealthyThreshold": 1}',
--     '{"Auto": true}',
--     '{}'
-- );

INSERT OR IGNORE INTO environments (id, tenant_id, name, strategy, healthcheck, rollback_policy, metadata)
VALUES (
    'env-1', 'tenant-1', 'production', 'blue_green',
    '{"Target": "https://example.com", "HealthyThreshold": 1, "UnhealthyThreshold": 1}',
    '{"Auto": true}',
    '{}'
);

INSERT OR IGNORE INTO releases (id, environment_id, artifact, git_tag, initiated_by, status, strategy_used, started_at, completed_at, release_notes, metadata)
VALUES (
    'rel-12345', 'env-1', 'acme/app:v2', '',
    '{"Type": "human", "UserID": "admin", "Name": "Admin"}',
    'pending', 'blue_green', '2026-05-19T18:00:00Z', NULL, '', '{}'
);
