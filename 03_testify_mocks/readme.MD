# Mockery: Automated Mock Generation for Go

Mockery is a tool that **auto-generates mocks for Go interfaces**.

It’s built on top of `testify/mock` and saves you from writing boilerplate by hand, keeping your mocks in sync when interfaces change.

---

## 1. Why `mock` is Needed?

When testing, you often don’t want to:

- Hit a real database
- Make actual API calls
- Depend on external services

Instead, you create a **mock object** that pretends to be the dependency. This way, you can:

- Define what the mock should return
- Verify if certain methods were called (and how many times, with what arguments)

⚠️ If you rename methods or add new ones, your hand-written mocks won’t compile until you update them. Mockery keeps them in sync automatically.

---

## 2. The Core of `testify/mock`

The key type is `mock.Mock`.

You **embed it** into your mock struct, which represents the dependency you’re faking.

---

## 3. Creating a Mock Manually

Say you have this interface:

```go
type UserRepository interface {
    GetUser(id int) (string, error)
}
```

With `testify/mock`, you’d write your own mock like this:

```go
import "github.com/stretchr/testify/mock"

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetUser(id int) (string, error) {
    args := m.Called(id) // track that it was called with `id`
    return args.String(0), args.Error(1)
}
```

---

## 4. Using the Mock in a Test

```go
func TestServiceGetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)

    // Tell the mock: "If GetUser(42) is called, return 'Alice', nil"
    mockRepo.On("GetUser", 42).Return("Alice", nil)

    // Inject mock into your service
    service := UserService{Repo: mockRepo}

    user, err := service.GetUser(42)

    assert.NoError(t, err)
    assert.Equal(t, "Alice", user)

    // Verify that expectations were met
    mockRepo.AssertExpectations(t)
}
```

---

## 5. Key Methods in `mock`

- **`On(method string, args...)`** → define expected calls
- **`Return(values...)`** → define return values for those calls
- **`Called(args...)`** → used inside your mock to track a call
- **`AssertExpectations(t)`** → verify all expectations were met
- **`AssertCalled(t, method, args...)`** → verify a method was called
- **`AssertNotCalled(t, method, args...)`** → verify a method was *not* called
- **`Once()`, `Twice()`, `Times(n)`** → control how many times a method should be called
- **`Run(func(args mock.Arguments))`** → add custom behavior (e.g., inspect or modify inputs)

### Example with `Run`

```go
mockRepo.On("GetUser", mock.Anything).
    Return("Bob", nil).
    Run(func(args mock.Arguments) {
        id := args.Int(0)
        fmt.Println("GetUser called with ID:", id)
    })
```

You can do all this manually. It works, but as interfaces grow, writing and updating mocks becomes repetitive and error-prone. That’s where **Mockery** comes in.

---

## 6. Why Mockery?

Mockery **automatically generates these mocks for you**, saving time and reducing errors.

---

## Installation

```bash
go install github.com/vektra/mockery/v2@latest
```

Verify:

```bash
mockery --version
```

---

## Generating Mocks

### Basic Usage

```bash
mockery --name=UserRepository
```

➡️ Produces `mocks/UserRepository.go` with a mock struct.

### Specifying Source

```bash
mockery --name=UserRepository --srcpkg=./path/to/package
```

Generate all mocks in a package:

```bash
mockery --all --output=./mocks
```

---

## Standard vs Expecter Mocks

Mockery supports two styles:

### 1. Standard (default)

```go
mockRepo := &mocks.UserRepository{}
mockRepo.On("GetUser", 42).Return("Alice", nil)
```

Advantages: Familiar, smaller, classic `testify/mock` API.

---

### 2. With Expecters (`with-expecter`)

```bash
mockery --name=UserRepository --with-expecter
```

```go
mockRepo := &mocks.UserRepository{}
mockRepo.EXPECT().GetUser(42).Return("Alice", nil)
```

Advantages: **Type-safe**, more readable, IDE-friendly, supports chaining.

---

### Side-by-Side

**Standard**

```go
mockRepo.On("GetUser", 42).Return("Alice", nil).Once()
```

**With Expecter**

```go
mockRepo.EXPECT().GetUser(42).Return("Alice", nil).Times(1)
```

💡 The **Expecter style** is preferred because it gives compile-time safety if method names or signatures change.

---

## Useful Options when Generating Mocks

| Flag | Description |
| --- | --- |
| `--output=./mocks` | Output directory |
| `--outpkg=mocks` | Package name |
| `--filename=user_repo_mock.go` | Custom file name |
| `--structname=MockUserRepo` | Custom struct name |
| `--case=snake` | Method naming style |
| `--disable-version-string` | Remove version comments |

---

## `.mockery.yaml` Config

Avoid repeating flags by creating a `.mockery.yaml`:

```yaml
with-expecter: true
output: mocks
outpkg: mocks
filename: "{{.InterfaceName}}_mock.go"
structname: "Mock{{.InterfaceName}}"
case: camel
disable-version-string: true
```

Advanced (per package/interface):

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

Then run:

```bash
mockery --config=.mockery.yaml
```

---

## Flow to Test with Mocks

1. Define an interface:

```go
type UserRepository interface {
    GetUser(id int) (User, error)
    SaveUser(user *User) error
    DeleteUser(id int) error
}
```

1. Generate the mocks:

```bash
mockery --name=UserRepository --with-expecter --output=./mocks
```

1. Use them in a test:

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

---

# Progressive Guide to Testify Mocks

When learning `testify/mock`, it’s best to go step by step.

We’ll start with the **simplest mock** and gradually add features like flexible arguments, multiple return values, call counts, custom behaviors, and advanced matchers.

---

## 1. The Simplest Mock

```go
type Greeter interface {
    Greet(name string) string
}
```

```go
import "github.com/stretchr/testify/mock"

type MockGreeter struct {
    mock.Mock
}

func (m *MockGreeter) Greet(name string) string {
    args := m.Called(name) // record the call
    return args.String(0)  // return the first value
}
```

Test:

```go
func TestSimpleMock(t *testing.T) {
    g := new(MockGreeter)

    g.On("Greet", "Alice").Return("Hello, Alice!")

    result := g.Greet("Alice")
    assert.Equal(t, "Hello, Alice!", result)

    g.AssertExpectations(t)
}
```

---

## 2. Matching with `mock.Anything`

```go
g.On("Greet", mock.Anything).Return("Hi there!")
```

---

## 3. Multiple Returns and Errors

```go
type UserRepository interface {
    GetUser(id int) (string, error)
}

type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) GetUser(id int) (string, error) {
    args := m.Called(id)
    return args.String(0), args.Error(1)
}
```

Usage:

```go
repo := new(MockUserRepository)
repo.On("GetUser", 1).Return("Alice", nil)
repo.On("GetUser", 2).Return("", errors.New("not found"))
```

---

## 4. Controlling Call Count

```go
repo.On("GetUser", 1).Return("Alice", nil).Once()
repo.On("GetUser", 2).Return("Bob", nil).Twice()
repo.On("GetUser", 3).Return("Charlie", nil).Times(5)
```

---

## 5. Verifying Calls Explicitly

```go
repo.AssertCalled(t, "GetUser", 1)
repo.AssertNotCalled(t, "GetUser", 999)
```

---

## 6. Using `Run` for Custom Behavior

```go
repo.On("GetUser", mock.Anything).
    Return("TestUser", nil).
    Run(func(args mock.Arguments) {
        id := args.Int(0)
        fmt.Println("GetUser called with:", id)
    })
```

---

## 7. Returning Dynamic Values

```go
repo.On("GetUser", mock.Anything).
    Return(func(id int) string {
        return fmt.Sprintf("User-%d", id)
    }, nil)
```

---

## 8. Chaining Expectations

```go
repo.On("GetUser", 1).Return("Alice", nil)
repo.On("GetUser", 2).Return("Bob", nil)
repo.On("GetUser", 3).Return("", errors.New("missing"))
```

---

## 9. Advanced Argument Matchers

```go
repo.On("GetUser", mock.MatchedBy(func(id int) bool {
    return id%2 == 0 // only even IDs
})).Return("EvenUser", nil)
```

---

# ✅ Summary

- `testify/mock` lets you build mocks manually.
- `mockery` generates them automatically.
- Use **standard** or **expecter** style depending on your needs.
- The **progressive guide** helps you master `mock` step by step.

## 📚 References:

- Testify/mock docs
- [Mockery GitHub](https://github.com/vektra/mockery)