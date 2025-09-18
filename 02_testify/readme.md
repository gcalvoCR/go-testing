# 📌 Go Testing with Testify

[Testify](https://github.com/stretchr/testify) is a Go library that makes tests **more readable and expressive** by providing powerful assertions, mocking capabilities, and test suites.

It builds on top of Go's `testing` package, so you still run tests with `go test`, but with much cleaner and more expressive syntax.

## Table of Contents

1. [What is Testify?](#-what-is-testify)
2. [Installation](#-step-1-install)
3. [Project Setup](#-project-setup)
4. [Basic Assertions](#-step-2-basic-assertions)
5. [assert vs require](#-step-3-assert-vs-require)
6. [Assertion Methods](#-step-4-useful-assertion-methods)
7. [Advanced Assertions](#-advanced-assertions)
8. [Test Suites](#-step-5-suites-optional)
9. [Mocking](#-step-6-mocks)
10. [Best Practices](#-best-practices)
11. [Example Files](#-example-files)

---

## 1. 📌 What is Testify?

[Testify](https://github.com/stretchr/testify) is a Go library that makes tests **more readable and expressive**.

It builds on top of Go's `testing` package (`*testing.T`), so you still run tests with `go test`.

It mainly gives you:

1. **Assertions** → cleaner checks (`assert.Equal`, `require.Equal`, etc.)
2. **Mocks** → for simulating dependencies
3. **Suites** → for organizing tests into groups

---

## 2. 🟢 Installation

```bash
go get github.com/stretchr/testify

```

---

## 3. 🟢 Project Setup

Before using testify, set up your Go module:

```bash
# Initialize a new module (if not already done)
go mod init your-project-name

# Install testify
go get github.com/stretchr/testify

# Your go.mod will look like:
module your-project-name

go 1.21

require (
    github.com/stretchr/testify v1.8.4
)
```

---

## 4. 🟢 Basic Assertions

Without Testify:

```go
if got != want {
    t.Errorf("got %d, want %d", got, want)
}

```

With Testify:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    assert.Equal(t, want, got) // clearer, concise
}

```

---

## 5. 🟢 `assert` vs `require`

- **`assert.*`** → logs a failure but continues the test.
- **`require.*`** → logs a failure and **stops immediately**.

```go
func TestDivide(t *testing.T) {
    assert.NotZero(t, 10, "denominator should not be zero")

    require.NotZero(t, 0, "oops, test will stop here")
    t.Log("this never runs")
}

```

👉 Use `require` when a failed check makes the rest of the test meaningless.

---

## 6. 🟢 Assertion Methods

- **Equality**
    
    ```go
    assert.Equal(t, 42, result)
    assert.NotEqual(t, 99, result)
    
    ```
    
- **Nil checks**
    
    ```go
    assert.Nil(t, err)
    assert.NotNil(t, dbConn)
    
    ```
    
- **Errors**
    
    ```go
    assert.Error(t, err)
    assert.NoError(t, err)
    
    ```
    
- **Boolean**
    
    ```go
    assert.True(t, condition)
    assert.False(t, otherCondition)
    
    ```
    
- **Collections**
    
    ```go
    assert.Contains(t, []string{"a", "b"}, "a")
    assert.Len(t, mySlice, 3)
    
    ```
    

👉 Testify has dozens more (`Greater`, `Subset`, `ElementsMatch`, etc.).

---

## 7. 🟢 Advanced Assertions

Testify provides many more powerful assertions:

- **Numeric comparisons**
    
    ```go
    assert.Greater(t, 10, 5)
    assert.GreaterOrEqual(t, 10, 10)
    assert.Less(t, 5, 10)
    assert.InDelta(t, 3.14, 3.14159, 0.01) // within delta
    ```
    
- **String operations**
    
    ```go
    assert.Contains(t, "hello world", "world")
    assert.NotContains(t, "hello", "goodbye")
    assert.HasPrefix(t, "prefix_test", "prefix_")
    assert.HasSuffix(t, "test_suffix", "_suffix")
    ```
    
- **Slice and map operations**
    
    ```go
    assert.Len(t, []int{1, 2, 3}, 3)
    assert.Empty(t, []string{})
    assert.NotEmpty(t, []int{1})
    assert.ElementsMatch(t, []int{1, 2}, []int{2, 1}) // same elements, any order
    ```
    
- **Type assertions**
    
    ```go
    assert.IsType(t, "string", result)
    assert.Implements(t, (*fmt.Stringer)(nil), myStruct)
    ```
    
- **HTTP and JSON**
    
    ```go
    assert.JSONEq(t, `{"key": "value"}`, responseBody)
    assert.HTTPStatusCode(t, http.StatusOK, resp.StatusCode)
    ```

---

## 8. 🟢 Test Suites

If you want setup/teardown logic and grouped tests, you can use **suites**.

```go
import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type MathSuite struct {
    suite.Suite
}

func (s *MathSuite) SetupTest() {
    // runs before each test
}

func (s *MathSuite) TestAdd() {
    s.Equal(4, Add(2, 2))
}

func (s *MathSuite) TestSub() {
    s.Equal(0, Sub(2, 2))
}

func TestMathSuite(t *testing.T) {
    suite.Run(t, new(MathSuite))
}

```

---

## 9. 🟢 Mocking

Testify’s **mocking framework** lets you simulate dependencies.

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

// Service interface
type Notifier interface {
    Send(msg string) error
}

// Mock implementation
type MockNotifier struct {
    mock.Mock
}

func (m *MockNotifier) Send(msg string) error {
    args := m.Called(msg)
    return args.Error(0)
}

// Test with mock
func TestNotifier(t *testing.T) {
    m := new(MockNotifier)
    m.On("Send", "hello").Return(nil)

    err := m.Send("hello")

    assert.NoError(t, err)
    m.AssertExpectations(t) // verifies all expected calls happened
}

```
---

## 10. 💡 Best Practices

### When to Use Testify
- **Use `assert`** for most cases - tests continue on failure
- **Use `require`** for setup/precondition checks that make further testing meaningless
- **Use suites** when you need setup/teardown logic across multiple tests
- **Use mocks** when testing code with external dependencies

### Test Organization
- **Import aliases** for cleaner code:
    
    ```go
    import (
        "testing"
        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
        "github.com/stretchr/testify/suite"
    )
    ```
    
- **Group related assertions** in helper functions
- **Use descriptive test names** that explain what they're testing
- **Keep test files focused** on testing one package/component

### Mocking Best Practices
- **Define interfaces** for dependencies you want to mock
- **Use `On().Return()`** to set up expected calls and return values
- **Call `AssertExpectations()`** to verify all expected calls were made
- **Use `AssertNotCalled()`** to ensure certain methods weren't called

### Error Handling
- **Test error conditions** thoroughly
- **Use `assert.Error()`** and `assert.NoError()` appropriately
- **Check error messages** when relevant
- **Use `require.NoError()`** for setup that must succeed

### Performance Considerations
- **Mock expensive operations** (database calls, HTTP requests, file I/O)
- **Use `t.Parallel()`** with testify assertions (they're thread-safe)
- **Avoid over-mocking** - sometimes integration tests are better

---

## 11. 📁 Example Files

This folder contains practical examples demonstrating testify usage:

### Core Examples
- **[`math.go`](math.go)** - Simple mathematical functions
- **[`math_test.go`](math_test.go)** - Tests using testify assertions
- **[`string_utils.go`](string_utils.go)** - String manipulation utilities
- **[`string_utils_test.go`](string_utils_test.go)** - Advanced testify testing patterns

### Advanced Examples
- **[`user_service.go`](user_service.go)** - Service with dependencies
- **[`user_service_test.go`](user_service_test.go)** - Mocking example
- **[`api_suite_test.go`](api_suite_test.go)** - Test suite example

### Running the Examples

First, install testify:

```bash
go get github.com/stretchr/testify
```

Then run the tests:

```bash
# Run all tests
go test -v

# Run specific test
go test -run TestAdd -v

# Run with coverage
go test -cover

# Run test suites
go test -run TestMathSuite -v
```

### Learning Points

These examples demonstrate:
- **Basic assertions** with `assert` and `require`
- **Table-driven tests** with testify
- **Mocking external dependencies**
- **Test suites** with setup/teardown
- **Error handling patterns**
- **HTTP testing** with testify
- **JSON validation**

Start with `math_test.go` to see basic testify usage, then explore the mocking and suite examples!