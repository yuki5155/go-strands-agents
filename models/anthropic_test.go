package models

import (
	"fmt"
	"reflect"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
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

func TestNewStreamingResponse(t *testing.T) {
	testcases := []struct {
		name     string
		model    string
		expected *StreamingResponse
	}{
		{
			name:  "creates response with model",
			model: "claude-sonnet-4-5-20250929",
			expected: &StreamingResponse{
				Model:                    "claude-sonnet-4-5-20250929",
				MessageID:                "",
				Role:                     "",
				Content:                  "",
				ContentBlockType:         "",
				ContentBlockIndex:        0,
				StopReason:               "",
				StopSequence:             "",
				InputTokens:              0,
				OutputTokens:             0,
				CacheCreationInputTokens: 0,
				CacheReadInputTokens:     0,
			},
		},
		{
			name:  "creates response with different model",
			model: "claude-opus-4",
			expected: &StreamingResponse{
				Model:                    "claude-opus-4",
				MessageID:                "",
				Role:                     "",
				Content:                  "",
				ContentBlockType:         "",
				ContentBlockIndex:        0,
				StopReason:               "",
				StopSequence:             "",
				InputTokens:              0,
				OutputTokens:             0,
				CacheCreationInputTokens: 0,
				CacheReadInputTokens:     0,
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			response := NewStreamingResponse(testcase.model)
			if !reflect.DeepEqual(response, testcase.expected) {
				t.Errorf("expected %+v, got %+v", testcase.expected, response)
			}
		})
	}
}

func TestStreamingResponse_StructFields(t *testing.T) {
	testcases := []struct {
		name     string
		response *StreamingResponse
	}{
		{
			name: "all fields set",
			response: &StreamingResponse{
				MessageID:                "msg_123",
				Model:                    "claude-sonnet-4-5",
				Role:                     "assistant",
				Content:                  "Hello World",
				ContentBlockType:         "text",
				ContentBlockIndex:        0,
				StopReason:               "end_turn",
				StopSequence:             "",
				InputTokens:              100,
				OutputTokens:             200,
				CacheCreationInputTokens: 50,
				CacheReadInputTokens:     30,
			},
		},
		{
			name: "minimal fields",
			response: &StreamingResponse{
				Model:   "test-model",
				Content: "test content",
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			r := testcase.response
			// Verify fields can be read (basic sanity check)
			_ = r.MessageID
			_ = r.Model
			_ = r.Role
			_ = r.Content
			_ = r.ContentBlockType
			_ = r.ContentBlockIndex
			_ = r.StopReason
			_ = r.StopSequence
			_ = r.InputTokens
			_ = r.OutputTokens
			_ = r.CacheCreationInputTokens
			_ = r.CacheReadInputTokens
		})
	}
}

func TestStreamingResponse_ProcessEvent(t *testing.T) {
	testcases := []struct {
		name          string
		eventType     string
		expectedDelta string
	}{
		{
			name:          "unknown event type",
			eventType:     "unknown_event_type",
			expectedDelta: "",
		},
		{
			name:          "message_start returns empty delta",
			eventType:     "message_start",
			expectedDelta: "",
		},
		{
			name:          "message_stop returns empty delta",
			eventType:     "message_stop",
			expectedDelta: "",
		},
		{
			name:          "message_delta returns empty delta",
			eventType:     "message_delta",
			expectedDelta: "",
		},
		{
			name:          "content_block_stop returns empty delta",
			eventType:     "content_block_stop",
			expectedDelta: "",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			response := NewStreamingResponse("test-model")
			event := anthropic.MessageStreamEventUnion{
				Type: testcase.eventType,
			}
			delta := response.ProcessEvent(event)
			if delta != testcase.expectedDelta {
				t.Errorf("expected delta '%s', got '%s'", testcase.expectedDelta, delta)
			}
		})
	}
}

func TestStreamingResponse_ContentAccumulation(t *testing.T) {
	testcases := []struct {
		name            string
		contentParts    []string
		expectedContent string
	}{
		{
			name:            "accumulates multiple parts",
			contentParts:    []string{"Hello", " ", "World", "!"},
			expectedContent: "Hello World!",
		},
		{
			name:            "single part",
			contentParts:    []string{"Hello"},
			expectedContent: "Hello",
		},
		{
			name:            "empty parts",
			contentParts:    []string{"", "", ""},
			expectedContent: "",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			response := NewStreamingResponse("test-model")
			for _, content := range testcase.contentParts {
				response.Content += content
			}
			if response.Content != testcase.expectedContent {
				t.Errorf("expected Content '%s', got '%s'", testcase.expectedContent, response.Content)
			}
		})
	}
}
