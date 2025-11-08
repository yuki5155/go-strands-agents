package models

import "github.com/yuki5155/go-strands-agents/utils"

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
