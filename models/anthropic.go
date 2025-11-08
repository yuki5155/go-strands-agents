package models

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/yuki5155/go-strands-agents/utils"
)

type AnthropicConfig struct {
	ModelId   string
	MaxTokens int64
	ApiKey    string
}

type Option func(c *AnthropicConfig)

// Functional Options Pattern
func WithModelId(modelId string) Option {
	return func(c *AnthropicConfig) {
		c.ModelId = modelId
	}
}

func WithMaxTokens(maxTokens int64) Option {
	return func(c *AnthropicConfig) {
		c.MaxTokens = maxTokens
	}
}

const DefaultModelId = "claude-sonnet-4-5-20250929"
const DefaultMaxTokens = 1024

// default is nil, will be set from envs
func WithApiKey(apiKey string) Option {
	return func(c *AnthropicConfig) {
		if apiKey != "" {
			c.ApiKey = apiKey
			return
		}
		key, ok := utils.GetApiKeyFromEnv()
		if !ok {
			panic("ANTHROPIC_API_KEY is not set")
		}
		c.ApiKey = key
	}
}

func NewAnthropicConfig(options ...Option) *AnthropicConfig {
	// Set defaults
	config := &AnthropicConfig{
		ModelId:   DefaultModelId,
		MaxTokens: DefaultMaxTokens,
	}
	// Apply options (can override defaults)
	for _, option := range options {
		option(config)
	}
	return config
}

// StreamingResponse represents the complete response from a streaming call
type StreamingResponse struct {
	MessageID                string
	Model                    string
	Role                     string
	Content                  string
	ContentBlockType         string
	ContentBlockIndex        int
	StopReason               string
	StopSequence             string
	InputTokens              int64
	OutputTokens             int64
	CacheCreationInputTokens int64
	CacheReadInputTokens     int64
}

// NewStreamingResponse creates a new StreamingResponse with the given model
func NewStreamingResponse(model string) *StreamingResponse {
	return &StreamingResponse{
		Model: model,
	}
}

// ProcessEvent processes a streaming event and updates the response accordingly
// Returns the text delta for content_block_delta events, empty string otherwise
func (r *StreamingResponse) ProcessEvent(event anthropic.MessageStreamEventUnion) string {
	switch event.Type {
	case "message_start":
		messageStart := event.AsMessageStart()
		r.MessageID = messageStart.Message.ID
		r.Role = string(messageStart.Message.Role)
		r.InputTokens = messageStart.Message.Usage.InputTokens
	case "content_block_start":
		blockStart := event.AsContentBlockStart()
		r.ContentBlockIndex = int(blockStart.Index)
		r.ContentBlockType = string(blockStart.ContentBlock.Type)
		// Check if content block has initial text
		if blockStart.ContentBlock.Text != "" {
			r.Content += blockStart.ContentBlock.Text
			return blockStart.ContentBlock.Text
		}
	case "content_block_delta":
		delta := event.AsContentBlockDelta()
		r.Content += delta.Delta.Text
		return delta.Delta.Text
	case "message_delta":
		messageDelta := event.AsMessageDelta()
		r.StopReason = string(messageDelta.Delta.StopReason)
		if messageDelta.Delta.StopSequence != "" {
			r.StopSequence = messageDelta.Delta.StopSequence
		}
		r.OutputTokens = messageDelta.Usage.OutputTokens
		r.CacheCreationInputTokens = messageDelta.Usage.CacheCreationInputTokens
		r.CacheReadInputTokens = messageDelta.Usage.CacheReadInputTokens
	}
	return ""
}
