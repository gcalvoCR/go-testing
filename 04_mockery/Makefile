# Makefile for generating mocks with Mockery

.PHONY: help mock mock-clean mock-user mock-io install-mockery test

# Default target
help:
	@echo "Available targets:"
	@echo "  mock          - Generate all mocks using .mockery.yaml config"
	@echo "  mock-clean    - Remove all generated mock files"
	@echo "  mock-user     - Generate mock for UserRepository only"
	@echo "  mock-io       - Generate mocks for Reader/Writer/ReadWriter interfaces"
	@echo "  install-mockery - Install mockery tool (run once)"
	@echo "  test          - Run tests to verify mocks work"
	@echo "  help          - Show this help message"

# Generate all mocks using explicit commands
mock: mock-user mock-io
	@echo "All mocks generated successfully"

# Clean all generated mock files
mock-clean:
	@echo "Removing generated mock files..."
	rm -rf mocks/*.go

# Generate mock for UserRepository with expecters
mock-user:
	@echo "Generating mock for UserRepository..."
	mockery --name=UserRepository --with-expecter --structname=MockUserRepositoryWithExpecter

# Generate mocks for I/O interfaces without expecters
mock-io:
	@echo "Generating mocks for I/O interfaces..."
	mockery --name=Reader --structname=MockReader
	mockery --name=Writer --structname=MockWriter
	mockery --name=ReadWriter --structname=MockReadWriter

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