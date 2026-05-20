# Continuous Deployment Engine (cd-engine)

A custom-built, lightweight Continuous Deployment orchestration engine written in Go.

Designed using **Domain-Driven Design (DDD)** and **Hexagonal Architecture**, this engine orchestrates complex deployment strategies (Blue/Green, Rolling, Canary) against target nodes via SSH. It includes fully immutable Infrastructure as Code (Terraform) and idempotent Configuration Management (Ansible) for provisioning target environments.

## Architecture & Core Concepts

- **Hexagonal Architecture:** Strict separation between the Core Domain (Strategy Engine, FSM), the Application Layer (Orchestrator), and the Adapters (SQLite, SSH, Health Checker).
- **Saga Pattern (Rollbacks):** The orchestrator utilizes a LIFO (Last-In-First-Out) stack for compensating transactions. If a deployment fails health checks, steps are automatically rolled back in reverse order.
- **Infrastructure as Code:** Immutable Terraform modules targeting AWS and GCP Always Free tiers, utilizing partial remote backend configurations.
- **Configuration Management:** Idempotent Ansible roles configuring raw Ubuntu VMs into hardened target nodes with Docker and Caddy (Zero-Downtime Reloads).

## Quick Start

This project is entirely driven by the standard `make` utility to ensure a smooth developer experience.

Ensure you have Go 1.22+ and `sqlite3` installed on your machine.

1. **Clone the repository:**

   ```bash
   git clone git@github.com:B-Nockk/devops-cd-engine.git
   cd devops-cd-engine
   ```

2. **View available commands:**

   ```bash
   make help
   ```

3. **Run a full local simulation:**

```bash
make clean
make seed       # Builds the binary, sets up SQLite tables, and injects test data
make run-deploy # Triggers the CLI with the test release ID

```

## Usage & Operations

To keep this repository clean and maintainable, all build, test, and infrastructure commands are abstracted behind the `Makefile`.

If you prefer to run the Go, Terraform, or Ansible commands manually, or need to see exactly what parameters are being passed to the infrastructure tools, run:

```bash
make manual-commands
```

This will print a complete reference guide to your terminal detailing the exact `go build`, `terraform init -backend-config=...`, and `ansible-playbook` commands used by this platform.

## CI/CD & Code Quality

This repository utilizes GitHub Actions (`.github/workflows/ci.yml`) for automated quality gates on every push and pull request:

- **Testing:** Standard Go race condition testing.
- **Linting:** Enforced Go idioms via `golangci-lint` (`make lint`).
- **Security:** Static application security testing (SAST) via `gosec` (`make audit`).
