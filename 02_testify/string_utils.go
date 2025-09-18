package main

import (
	"strings"
)

// Reverse returns the reverse of a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome (case-insensitive)
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// Capitalize capitalizes the first letter of each word
func Capitalize(s string) string {
	return strings.Title(s)
}

// CountWords counts the number of words in a string
func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}
