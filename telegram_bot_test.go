package main

import (
	"slices"
	"testing"
)

func TestComputeResponse(t *testing.T) {
	triggers = []Trigger{
		{Key: "test", Values: "Response for test"},
		{Key: "hello", Values: []interface{}{"Hello!", "Hi there"}},
	}

	testCases := []struct {
		Name           string
		InputMessage   string
		ExpectedOutput []string
	}{
		{"Exact match", "test", []string{"Response for test"}},
		{"Case insensitive match", "tEsT", []string{"Response for test"}},
		{"Array random match", "Hello", []string{"Hello!", "Hi there"}},
		{"No match", "This doesn't match", []string{""}},
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
