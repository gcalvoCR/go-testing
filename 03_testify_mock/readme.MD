## Testify mock

## 1. Why `mock` is needed

When testing, you often don’t want to:

- Hit a real database
- Make actual API calls
- Depend on slow or unreliable external services

Instead, you create a **mock object** that pretends to be the dependency. This way, you can:

- Define what the mock should return
- Verify if certain methods were called (and how many times, with what arguments)

---

## 2. The core of `testify/mock`

The key type is `mock.Mock`.

You **embed it** into your mock struct, which represents the dependency you’re faking.

---

## 3. Creating a mock

Say you have this interface:

```go
type UserRepository interface {
    GetUser(id int) (string, error)
}

```

You want to mock it. With `testify/mock`, you’d do:

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

## 4. Using the mock in a test

Now you can define what should happen when `GetUser` is called:

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

## 5. Key methods in `mock`

- **`On(method string, arguments...)`** → define expected calls
- **`Return(values...)`** → define return values for those calls
- **`Called(arguments...)`** → used inside your mock implementation to track a call
- **`AssertExpectations(t)`** → check if all expected calls happened
- **`AssertCalled(t, method, args...)`** → check if a method was called with args
- **`AssertNotCalled(t, method, args...)`** → check if a method was *not* called
- **`Once()`, `Twice()`, `Times(n)`** → control how many times a method should be called
- **`Run(func(args mock.Arguments))`** → let you add custom behavior (e.g., modify inputs)

---

## 6. Example with `Run`

```go
mockRepo.On("GetUser", mock.Anything).Return("Bob", nil).Run(func(args mock.Arguments) {
    id := args.Int(0)
    fmt.Println("GetUser called with ID:", id)
})

```

This lets you inspect or react dynamically to calls.