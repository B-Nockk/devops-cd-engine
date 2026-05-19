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
