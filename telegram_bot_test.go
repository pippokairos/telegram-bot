package main

import (
	"slices"
	"testing"
)

// Stub the triggers without reading from file
func setMockTriggers() {
	triggers = []Trigger{
		{Key: "test", Values: "Response for test"},
		{Key: "hello", Values: []interface{}{"Hello!", "Hi there"}},
		{Key: "say hi to", Values: []interface{}{"Hi __input__, how are you?"}},
	}
}

func TestComputeResponse(t *testing.T) {
	setMockTriggers()

	testCases := []struct {
		Name           string
		InputMessage   string
		ExpectedOutput []string
	}{
		{"Exact match", "test", []string{"Response for test"}},
		{"Case insensitive match", "tEsT", []string{"Response for test"}},
		{"Array random match", "Hello", []string{"Hello!", "Hi there"}},
		{"No match", "This doesn't match", []string{""}},
		{"Match with input", "Hey! Say hi to John Doe", []string{"Hi John Doe, how are you?"}},
		{"Match with input trimmed", "Say hi to  John Doe   ", []string{"Hi John Doe, how are you?"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output := computeResponse(testCase.InputMessage)
			if !slices.Contains(testCase.ExpectedOutput, output) {
				t.Errorf("Unexpected response. Got: %s, Expected: %s", output, testCase.ExpectedOutput)
			}
		})
	}
}
