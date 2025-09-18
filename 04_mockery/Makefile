# Makefile for generating mocks with Mockery

.PHONY: help mock mock-clean mock-user mock-io install-mockery test

# Default target
help:
	@echo "Available targets:"
	@echo "  mock          - Generate all mocks using .mockery.yaml config"
	@echo "  mock-clean    - Remove all generated mock files"
	@echo "  install-mockery - Install mockery tool (run once)"
	@echo "  test          - Run tests to verify mocks work"
	@echo "  help          - Show this help message"

# Generate all mocks using the .mockery.yaml configuration
mock:
	@echo "Generating mocks for all interfaces..."
	mockery --config=.mockery.yaml

# Clean all generated mock files
mock-clean:
	@echo "Removing generated mock files..."
	rm -rf mocks/*.go

# Generate individual mocks
mock-user:
	@echo "Generating mock for UserRepository..."
	mockery --name=UserRepository --with-expecter --structname=MockUserRepositoryWithExpecter --output=mocks --outpkg=mocks

mock-io:
	@echo "Generating mocks for I/O interfaces..."
	mockery --name=Reader --structname=MockReader --output=mocks --outpkg=mocks
	mockery --name=Writer --structname=MockWriter --output=mocks --outpkg=mocks
	mockery --name=ReadWriter --structname=MockReadWriter --output=mocks --outpkg=mocks

# Install mockery (run this once)
install-mockery:
	@echo "Installing mockery..."
	go install github.com/vektra/mockery/v2@latest

# Run tests to verify mocks work
test:
	@echo "Running tests..."
	go test -v ./...

# Regenerate all mocks (clean + generate)
mock-regen: mock-clean mock
	@echo "Mocks regenerated successfully"