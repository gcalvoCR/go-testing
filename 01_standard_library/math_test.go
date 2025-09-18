package main

import (
	"fmt"
	"testing"
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
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
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
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
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
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		expected    int
		expectError bool
	}{
		{"positive numbers", 6, 3, 2, false},
		{"negative numbers", -6, -3, 2, false},
		{"mixed numbers", 10, -2, -5, false},
		{"division by zero", 5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if tt.expectError {
				if err == nil {
					t.Errorf("Divide(%d, %d) expected error but got none", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func FuzzAdd(f *testing.F) {
	// Seed with some values
	f.Add(1, 2)
	f.Add(-3, 7)
	f.Add(0, 0)

	// Define the fuzz target
	f.Fuzz(func(t *testing.T, a, b int) {
		fmt.Println(a, b)
		got := Add(a, b)

		// Basic property: commutativity → a+b == b+a
		if got != Add(b, a) {
			t.Errorf("Add not commutative: Add(%d,%d)=%d, Add(%d,%d)=%d",
				a, b, got, b, a, Add(b, a))
		}

		// Identity property: a+0 == a
		if Add(a, 0) != a {
			t.Errorf("Add(%d,0) != %d", a, a)
		}

		// Overflow check (optional): prevent silent wraparound
		if (a > 0 && b > 0 && got < 0) || (a < 0 && b < 0 && got > 0) {
			t.Errorf("integer overflow: Add(%d,%d)=%d", a, b, got)
		}
	})
}
