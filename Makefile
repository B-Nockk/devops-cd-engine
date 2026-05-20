.DEFAULT_GOAL := help
BINARY=./bin/cd-engine

## ----------------------------------------------------------------------
## cd-engine Makefile
## ----------------------------------------------------------------------

help: ## Show this help message.
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compiles the Go binary into ./bin/
	go build -o $(BINARY) ./cmd/cd-engine

seed: build ## Initializes the SQLite DB and injects test data
	$(BINARY) setup
	sqlite3 data.db < ./scripts/seed.sql

run-deploy: build ## Runs a test deployment against the seeded data
	$(BINARY) deploy rel-12345

lint: ## Runs golangci-lint (Requires local installation)
	golangci-lint run ./...

audit: ## Runs gosec security scanner (Requires local installation)
	gosec -no-fail -fmt=golint ./...

test-coverage: ## Runs unit tests and outputs HTML coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Removes built binaries and local databases
	rm -rf ./bin data.db coverage.out

manual-commands: ## Displays the raw underlying commands for manual execution
	@echo "=========================================================="
	@echo "                MANUAL COMMAND REFERENCE                  "
	@echo "=========================================================="
	@echo ""
	@echo "[Go Engine]"
	@echo "  Build:           go build -o ./bin/cd-engine ./cmd/cd-engine"
	@echo "  Setup DB:        ./bin/cd-engine setup"
	@echo "  Deploy:          ./bin/cd-engine deploy <release_id>"
	@echo ""
	@echo "[Terraform (GCP)]"
	@echo "  Init:            cd infra/terraform/gcp && terraform init -backend-config=backend.gcp.example.hcl"
	@echo "  Plan:            cd infra/terraform/gcp && terraform plan"
	@echo "  Apply:           cd infra/terraform/gcp && terraform apply -auto-approve"
	@echo ""
	@echo "[Terraform (AWS)]"
	@echo "  Init:            cd infra/terraform/aws && terraform init -backend-config=backend.aws.example.hcl"
	@echo "  Plan:            cd infra/terraform/aws && terraform plan"
	@echo "  Apply:           cd infra/terraform/aws && terraform apply -auto-approve"
	@echo ""
	@echo "[Ansible]"
	@echo "  Run Playbook:    cd infra/ansible && ansible-playbook -i inventory.sample.ini playbook.yml"
	@echo ""
	@echo "=========================================================="
