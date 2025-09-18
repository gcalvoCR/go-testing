# Go Testing with Standard Library

This guide covers Go testing using only the standard library's `testing` package. It provides comprehensive examples and best practices for writing effective tests in Go.

## Table of Contents

1. [Basic Test Lifecycle Methods](#1--basic-test-lifecycle-methods)
2. [Logging](#2--logging)
3. [Skipping Tests](#3--skipping-tests)
4. [Cleanup Functions](#4--cleanup-functions)
5. [Subtests](#5--subtests)
6. [Timing & Parallelism](#6--timing--parallelism)
7. [Metadata & State](#7--metadata--state)
8. [How to Run Tests](#8--how-to-run-tests)
9. [Best Practices](#9--best-practices)
10. [Putting It All Together](#-putting-it-all-together)
11. [Example Files](#-example-files)

---

## 1. ✅ Basic Test Lifecycle Methods

These are the **core failure and skip controls**.

- **`t.Fail()`**
    
    Marks the test as failed, but continues running.
    
    ```go
    func TestFail(t *testing.T) {
        t.Fail() // test marked as failed, but code continues
        t.Log("This still runs")
    }
    ```
    
- **`t.FailNow()`**
    
    Marks the test as failed and **stops immediately**.
    
    ```go
    func TestFailNow(t *testing.T) {
        t.FailNow()
        t.Log("This never runs")
    }
    ```
    
- **`t.Fatal(args...)`**
    
    Like `t.Log` + `t.FailNow()`. Logs a message and aborts.
    
    ```go
    func TestFatal(t *testing.T) {
        if 2+2 != 4 {
            t.Fatal("Math is broken")
        }
    }
    ```
    
- **`t.FatalF(format, args...)`**
    
    Same as `Fatal`, but with `fmt.Printf`-style formatting.
    
- **`t.Error(args...)`**
    
    Like `t.Log` + `t.Fail()`. Marks failure but continues.
    
    ```go
    func TestError(t *testing.T) {
        if 2+2 != 5 {
            t.Error("Expected 5, got 4")
        }
        t.Log("This still executes")
    }
    ```
    
- **`t.Errorf(format, args...)`**
    
    Same as `Error`, but formatted.
    

👉 Rule of thumb:

- Use `Error`/`Errorf` when you want to **check multiple conditions**.
- Use `Fatal`/`Fatalf` when you want to **abort immediately**.

---

## 2. 📝 Logging

These let you write messages to the test output.

- **`t.Log(args...)`**
    
    Logs values.
    
    ```go
    t.Log("this is a log message")
    ```
    
- **`t.Logf(format, args...)`**
    
    Logs with formatting.
    
    ```go
    t.Logf("Got %d, expected %d", got, want)
    ```
    
- **`t.Helper()`**
    
    Marks the calling function as a helper. This makes error logs show the *caller* of the helper, not the helper itself.
    
    ```go
    func assertEqual(t *testing.T, got, want int) {
        t.Helper() // error points to caller
        if got != want {
            t.Errorf("got %d, want %d", got, want)
        }
    }
    ```
    

---

## 3. ⏩ Skipping Tests

Sometimes you want to skip tests based on environment.

- **`t.Skip(args...)`**
    
    Marks test as skipped and stops.
    
    ```go
    func TestSkip(t *testing.T) {
        if testing.Short() {
            t.Skip("skipping in short mode")
        }
    }
    ```
    
- **`t.Skipf(format, args...)`**
    
    Same, but formatted message.
    
- **`t.SkipNow()`**
    
    Immediately skips without logging a message.
    

---

## 4. 🧹 Cleanup Functions

- **`t.Cleanup(func())`**
    
    Registers a function to run when the test (or subtest) ends.
    
    ```go
    func TestCleanup(t *testing.T) {
        t.Cleanup(func() {
            t.Log("cleanup ran")
        })
        t.Log("test body")
    }
    ```
    

Useful for closing DB connections, removing files, etc.

---

## 5. 🔄 Subtests

- **`t.Run(name string, f func(t *testing.T)) bool`**
    
    Runs a subtest inside your test. Returns `true` if it passed.
    
    ```go
    func TestSubtests(t *testing.T) {
        tests := []struct {
            name string
            in   int
            want int
        }{
            {"double of 2", 2, 4},
            {"double of 3", 3, 6},
        }
    
        for _, tc := range tests {
            t.Run(tc.name, func(t *testing.T) {
                got := tc.in * 2
                if got != tc.want {
                    t.Errorf("got %d, want %d", got, tc.want)
                }
            })
        }
    }
    
    ```
    

This helps parameterize tests.

---

## 6. ⏱️ Timing & Parallelism

- **`t.Parallel()`**
    
    Marks test (or subtest) to run in parallel with others.
    
    ```go
    func TestParallel(t *testing.T) {
        t.Parallel()
        // runs alongside other parallel tests
    }
    
    ```
    

---

## 7. ℹ️ Metadata & State

- **`t.Name()`**
    
    Returns the test’s name.
    
    ```go
    t.Log("Running:", t.Name())
    
    ```
    
- **`t.TempDir()`**
    
    Creates a temporary directory that is automatically cleaned up.
    
    ```go
    func TestTempDir(t *testing.T) {
        dir := t.TempDir()
        t.Log("temp dir:", dir)
    }
    ```
    
- **`t.Setenv(key, value)`**
    
    Sets an environment variable for the test, automatically restored after.
    
    ```go
    func TestSetenv(t *testing.T) {
        t.Setenv("MY_ENV", "123")
        t.Log("env set inside test")
    }
    ```
    

---

## 8. 🏃 How to Run Tests

Go provides several command-line flags to control test execution:

- **`go test`**
    
    Runs all tests in the current package.
    
- **`go test -v`**
    
    Verbose output - shows all test names and results.
    
- **`go test -run TestName`**
    
    Runs only tests matching the pattern (supports regex).
    
    ```bash
    go test -run TestAdd  # runs TestAdd and TestAddNumbers
    ```
    
- **`go test -short`**
    
    Skips tests that take too long (checks `testing.Short()`).
    
- **`go test -cover`**
    
    Shows test coverage percentage.
    
- **`go test -coverprofile=coverage.out`**
    
    Generates a coverage profile file.
    
- **`go test -race`**
    
    Detects race conditions in concurrent code.
    
- **`go test -parallel n`**
    
    Sets the number of parallel tests (default is GOMAXPROCS).
    
- **`go test -bench=.`**
    
    Runs benchmarks (functions starting with Benchmark).
    
- **`go test -timeout 30s`**
    
    Sets a timeout for test execution.

---

## 9. 💡 Best Practices

### Test Organization
- **Test file naming**: Use `_test.go` suffix for test files
- **Test function naming**: Start with `Test` followed by the function name
- **Table-driven tests**: Use structs for multiple test cases
- **Helper functions**: Use `t.Helper()` for reusable test logic

### Writing Good Tests
- **Test one thing**: Each test should verify a single behavior
- **Descriptive names**: Test names should describe what they're testing
- **Arrange-Act-Assert**: Structure tests clearly
- **Avoid dependencies**: Tests should be independent and isolated

### Assertions and Failures
- **Use appropriate failure methods**:
  - `t.Error`/`t.Errorf` for non-blocking failures
  - `t.Fatal`/`t.Fatalf` for blocking failures
- **Provide clear error messages**: Include expected vs actual values
- **Use helper functions**: Keep test code clean and readable

### Test Coverage
- **Aim for meaningful coverage**: Focus on critical paths
- **Don't test implementation details**: Test behavior, not internals
- **Use coverage tools**: `go test -cover` to measure coverage

### Performance and Reliability
- **Use `t.Parallel()`** for independent tests
- **Set timeouts**: Use `go test -timeout` for long-running tests
- **Clean up resources**: Use `t.Cleanup()` for proper teardown
- **Handle race conditions**: Use `go test -race` to detect issues

### Example of a Well-Structured Test

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -1, -2},
        {"mixed numbers", 5, -3, 2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

---

# 🚀 Putting It All Together

Here’s a test that combines multiple features:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    t.Setenv("APP_MODE", "test")
    dir := t.TempDir()

    t.Cleanup(func() {
        t.Log("Cleaning up resources")
    })

    t.Run("subtest1", func(t *testing.T) {
        t.Parallel()
        got := 2 * 2
        want := 4
        if got != want {
            t.Fatalf("got %d, want %d", got, want)
        }
    })

    t.Run("subtest2", func(t *testing.T) {
        got := 2 + 2
        want := 5
        if got != want {
            t.Error("bad math")
        }
    })
}
```
---

## 11. 📁 Example Files

This folder contains practical examples to help you get started with Go testing:

### Core Files
- **[`math.go`](math.go)** - Simple mathematical functions (Add, Subtract, Multiply, Divide)
- **[`math_test.go`](math_test.go)** - Comprehensive tests for the math functions using table-driven tests
- **[`string_utils.go`](string_utils.go)** - String manipulation utilities (Reverse, IsPalindrome, Capitalize, CountWords)
- **[`string_utils_test.go`](string_utils_test.go)** - Tests demonstrating various testing patterns

### How to Run the Examples

Navigate to this directory and run:

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific test file
go test -run TestAdd

# Run with coverage
go test -cover

# Run with race detection
go test -race
```

### Learning Points

These examples demonstrate:
- **Table-driven tests** for multiple test cases
- **Error handling** in tests
- **Subtests** for organizing related test cases
- **Helper functions** and clean test structure
- **Edge cases** and boundary testing
- **Unicode support** in string operations

Start by examining the test files to see how the concepts from this guide are applied in practice!