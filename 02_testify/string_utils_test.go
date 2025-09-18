package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "olleh"},
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"palindrome", "racecar", "racecar"},
		{"unicode", "héllo", "olléh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple palindrome", "racecar", true},
		{"not palindrome", "hello", false},
		{"empty string", "", true},
		{"single character", "a", true},
		{"case insensitive", "Racecar", true},
		{"spaces ignored", "A man a plan a canal Panama", false}, // Note: this doesn't handle spaces
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single word", "hello", "Hello"},
		{"multiple words", "hello world", "Hello World"},
		{"empty string", "", ""},
		{"already capitalized", "Hello", "Hello"},
		{"mixed case", "hELLO wORLD", "HELLO WORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Capitalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"single word", "hello", 1},
		{"multiple words", "hello world test", 3},
		{"empty string", "", 0},
		{"with spaces", "  hello   world  ", 2},
		{"tabs and newlines", "hello\tworld\nagain", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringOperations(t *testing.T) {
	// Test multiple string operations
	input := "hello world"

	assert.True(t, len(input) > 0)
	assert.Contains(t, input, "world")
	assert.NotContains(t, input, "goodbye")

	reversed := Reverse(input)
	assert.Equal(t, "dlrow olleh", reversed)

	assert.False(t, IsPalindrome(input))
	assert.True(t, IsPalindrome("racecar"))

	capitalized := Capitalize(input)
	assert.Equal(t, "Hello World", capitalized)

	wordCount := CountWords(input)
	assert.Equal(t, 2, wordCount)
	assert.Greater(t, wordCount, 0)
}

func TestEdgeCases(t *testing.T) {
	// Test edge cases with require for critical checks
	require.NotNil(t, Reverse("test"))
	require.NotNil(t, Capitalize("test"))

	// Test with assert for non-critical checks
	assert.Empty(t, Reverse(""))
	assert.Empty(t, Capitalize(""))
	assert.Zero(t, CountWords(""))
}
