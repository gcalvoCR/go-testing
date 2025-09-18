package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -1, -1, -2},
		{"mixed numbers", 5, -3, 2},
		{"zero addition", 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "Add(%d, %d) should equal %d", tt.a, tt.b, tt.expected)
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 2},
		{"negative numbers", -1, -1, 0},
		{"mixed numbers", 5, -3, 8},
		{"zero subtraction", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 6},
		{"negative numbers", -2, -3, 6},
		{"mixed numbers", 5, -3, -15},
		{"zero multiplication", 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("successful division", func(t *testing.T) {
		result, err := Divide(10, 2)
		require.NoError(t, err)
		assert.Equal(t, 5, result)
	})

	t.Run("division by zero", func(t *testing.T) {
		result, err := Divide(10, 0)
		assert.Error(t, err)
		assert.Equal(t, "division by zero", err.Error())
		assert.Equal(t, 0, result)
	})

	t.Run("negative division", func(t *testing.T) {
		result, err := Divide(-10, 2)
		require.NoError(t, err)
		assert.Equal(t, -5, result)
	})
}

func TestMathOperations(t *testing.T) {
	// Test multiple operations in one test
	assert.Equal(t, 8, Add(3, 5))
	assert.Equal(t, 15, Multiply(3, 5))
	assert.Equal(t, -2, Subtract(3, 5))

	// Test with require for critical checks
	result, err := Divide(10, 2)
	require.NoError(t, err)
	assert.Equal(t, 5, result)
}
