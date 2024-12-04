.PHONY: dev prod test cover

dev:
	docker-compose -f docker-compose.dev.yml up --build

prod:
	docker-compose up --build

# Run tests with coverage
test:
	go test ./... -coverprofile=coverage.out

# Generate and open the coverage report
cover: test
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html # For macOS; use `xdg-open` for Linux or `start` for Windows