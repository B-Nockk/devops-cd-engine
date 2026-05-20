BINARY=./bin/cd-engine

build:
	go build -o $(BINARY) ./cmd/cd-engine

seed: build
	# Run the Go setup command to create tables
	$(BINARY) setup
	# Inject the test data
	sqlite3 data.db < ./scripts/seed.sql

run-deploy: build
	$(BINARY) deploy rel-12345

test:
	go test ./...

clean:
	rm -rf ./bin data.db

lint:
	# Requires golangci-lint installed locally
	golangci-lint run ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.outlint:
	# Requires golangci-lint installed locally
	golangci-lint run ./...

audit:
	# Requires gosec installed locally
	gosec -no-fail -fmt=golint ./...

# test-coverage:
# 	go test -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out
