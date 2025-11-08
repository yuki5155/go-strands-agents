package models

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewAnthropicConfig(t *testing.T) {
	testcases := []struct {
		name     string
		options  []Option
		expected *AnthropicConfig
	}{
		{
			name:    "default",
			options: []Option{},
			expected: &AnthropicConfig{
				ModelId:   "claude-sonnet-4-5-20250929",
				MaxTokens: 1024,
				ApiKey:    "",
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			config := NewAnthropicConfig(testcase.options...)
			fmt.Println(config)
			if !reflect.DeepEqual(config, testcase.expected) {
				t.Errorf("expected %+v, got %+v", testcase.expected, config)
			}
		})
	}
}
