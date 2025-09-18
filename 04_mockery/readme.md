# Mockery: Automated Mock Generation for Go

Mockery is a powerful Go tool built on top of `testify/mock` that automatically generates mock implementations for your interfaces. It eliminates the need to write boilerplate mock code by hand, ensuring your mocks stay up-to-date with interface changes.

## Installation

Install Mockery using Go's module system:

```bash
go install github.com/vektra/mockery/v2@latest
```

This installs the latest version of Mockery globally. You can verify the installation by running:

```bash
mockery --version
```

## Generating Mocks

### Basic Usage

To generate a mock for an interface, run Mockery with the `--name` flag:

```bash
mockery --name=UserRepository
```

This generates a file like `mocks/UserRepository.go` containing a mock struct that implements the `UserRepository` interface.

### Specifying Source Files

If Mockery can't find the interface automatically, specify the source file:

```bash
mockery --name=UserRepository --srcpkg=./path/to/package
```

Or generate mocks for all interfaces in a package:

```bash
mockery --all --output=./mocks
```

## Mockery Options and Benefits

Mockery offers various options to customize mock generation. The most important distinction is between standard mocks and mocks with expecters.

### Standard Mocks (Without `--with-expecter`)

This is the default mode. Generated mocks use the classic `testify/mock` API:

```go
mockRepo := &mocks.UserRepository{}
mockRepo.On("GetUser", 42).Return("Alice", nil)
```

#### Benefits:
- **Simplicity**: Familiar API if you're already using `testify/mock`
- **Flexibility**: Direct access to all `testify/mock` features
- **Smaller Code**: Less generated code compared to expecter mocks

### Mocks with Expecters (`--with-expecter`)

Use the `--with-expecter` flag to generate mocks with a fluent expecter API:

```bash
mockery --name=UserRepository --with-expecter
```

This generates additional code allowing a more readable, type-safe API:

```go
mockRepo := &mocks.UserRepository{}
mockRepo.EXPECT().GetUser(42).Return("Alice", nil)
```

#### Benefits:
- **Type Safety**: Compile-time checking of method names and parameters
- **Readability**: Fluent API that's easier to read and write
- **IDE Support**: Better autocomplete and refactoring support
- **Chaining**: Allows chaining expectations in a single statement

#### Comparison:

**Standard:**
```go
mockRepo.On("GetUser", 42).Return("Alice", nil).Once()
mockRepo.On("SaveUser", mock.AnythingOfType("*User")).Return(nil)
```

**With Expecter:**
```go
mockRepo.EXPECT().GetUser(42).Return("Alice", nil).Times(1)
mockRepo.EXPECT().SaveUser(mock.AnythingOfType("*User")).Return(nil)
```

### Other Useful Options

- `--output=./mocks`: Specify output directory
- `--outpkg=mocks`: Set the package name for generated mocks
- `--filename=user_repo_mock.go`: Custom filename for the generated mock
- `--structname=MockUserRepo`: Custom struct name
- `--case=snake`: Use snake_case for method names (default is camelCase)
- `--disable-version-string`: Remove version strings from generated code

Example with multiple options:

```bash
mockery \
  --name=UserRepository \
  --with-expecter \
  --output=./mocks \
  --outpkg=mocks \
  --filename=user_repo_mock.go \
  --structname=MockUserRepo
```

## Using .mockery.yaml Configuration File

Mockery supports configuration via a `.mockery.yaml` file in your project root. This allows you to set default options and avoid repeating flags.

### Basic Configuration

Create a `.mockery.yaml` file:

```yaml
with-expecter: true
output: mocks
outpkg: mocks
filename: "{{.InterfaceName}}_mock.go"
structname: "Mock{{.InterfaceName}}"
case: camel
disable-version-string: true
```

### Advanced Configuration

You can configure different settings for different packages or interfaces:

```yaml
packages:
  github.com/your/project/repository:
    config:
      with-expecter: true
      output: mocks
    interfaces:
      UserRepository:
        config:
          structname: MockUserRepo
      Database:
        config:
          with-expecter: false
```

### Configuration Options

- `with-expecter`: Generate expecter methods (default: false)
- `output`: Output directory (default: "mocks")
- `outpkg`: Package name for mocks (default: "mocks")
- `filename`: Template for filename (supports `{{.InterfaceName}}`)
- `structname`: Template for struct name (supports `{{.InterfaceName}}`)
- `case`: Case style for methods ("camel" or "snake")
- `disable-version-string`: Remove version comments (default: false)
- `mock-build-tags`: Build tags for generated mocks

### Using Configuration

With a `.mockery.yaml` file, you can run Mockery with minimal flags:

```bash
mockery --name=UserRepository
```

Mockery will use the configuration from `.mockery.yaml` automatically.

## Examples

### Complete Workflow

1. Define your interface:

```go
type UserRepository interface {
    GetUser(id int) (User, error)
    SaveUser(user *User) error
    DeleteUser(id int) error
}
```

2. Generate mock with expecters:

```bash
mockery --name=UserRepository --with-expecter --output=./mocks
```

3. Use in tests:

```go
func TestUserService_GetUser(t *testing.T) {
    mockRepo := &mocks.UserRepository{}
    service := NewUserService(mockRepo)

    expectedUser := User{ID: 42, Name: "Alice"}
    mockRepo.EXPECT().GetUser(42).Return(expectedUser, nil)

    user, err := service.GetUser(42)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

## Benefits Summary

- **Productivity**: Eliminates manual mock writing
- **Maintenance**: Automatically updates mocks when interfaces change
- **Consistency**: Standardized mock implementations
- **Flexibility**: Multiple generation options for different needs
- **Integration**: Works seamlessly with `testify/mock` ecosystem

For more advanced usage and options, refer to the [official Mockery documentation](https://github.com/vektra/mockery).
