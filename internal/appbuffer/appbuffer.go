// Description: This package provides string buffer management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appbuffer

import "strings"

// Represents a buffer for storing and manipulating a slice of strings.
// It provides methods for appending, joining, clearing, and retrieving strings.
type AppBuffer struct {
	// A slice of strings used to store the buffer.
	buffer []string
}

// Appends the given string to the buffer and returns the buffer.
// It allows method chaining to continue appending strings.
func (buff *AppBuffer) AppendStr(str string) *AppBuffer {
	buff.buffer = append(buff.buffer, str)
	return buff
}

// Returns a copy of the buffer as a slice of strings.
func (buff *AppBuffer) GetBuffer() []string {
	return buff.buffer
}

// Resets the buffer to an empty slice of strings.
func (buff *AppBuffer) Clear() {
	buff.buffer = make([]string, 0)
}

// Concatenates the strings in the buffer with the given separator and returns the concatenated string.
// The separator is placed between each string in the buffer, but not at the beginning or end of the resulting string.
// If the separator is an empty string, the strings in the buffer are concatenated without any separator.
// The resulting string is returned as a new string.
func (buff *AppBuffer) Join(sep string) string {
	return strings.Join(buff.buffer, sep)
}

// Returns the number of strings in the buffer.
func (buff *AppBuffer) Len() int {
	return len(buff.buffer)
}
