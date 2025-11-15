package models

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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
	Channel                  chan string
}

// NewStreamingResponse creates a new StreamingResponse with the given model
func NewStreamingResponse(model string) *StreamingResponse {
	return &StreamingResponse{
		Model:   model,
		Channel: make(chan string),
	}
}

// GetChannel returns the channel for streaming text deltas
func (r *StreamingResponse) GetChannel() chan string {
	return r.Channel
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
			if r.Channel != nil {
				r.Channel <- blockStart.ContentBlock.Text
			}
			return blockStart.ContentBlock.Text
		}
	case "content_block_delta":
		delta := event.AsContentBlockDelta()
		r.Content += delta.Delta.Text
		if r.Channel != nil {
			r.Channel <- delta.Delta.Text
		}
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

// create client struct for anthropic
type AnthropicClient struct {
	Client anthropic.Client
	Config *AnthropicConfig
}

func NewAnthropicClient(options ...Option) *AnthropicClient {
	// Create config with provided options
	config := NewAnthropicConfig(options...)

	// If ApiKey is still empty, try to get it from env
	if config.ApiKey == "" {
		key, ok := utils.GetApiKeyFromEnv()
		if !ok {
			panic("ANTHROPIC_API_KEY is not set")
		}
		config.ApiKey = key
	}

	return &AnthropicClient{
		Client: anthropic.NewClient(
			option.WithAPIKey(config.ApiKey),
		),
		Config: config,
	}
}

// StreamMessages sends messages and streams the response with optional callback for each text delta
func (c *AnthropicClient) StreamMessages(ctx context.Context, messages []anthropic.MessageParam, onDelta func(string)) (*StreamingResponse, error) {
	response := NewStreamingResponse(c.Config.ModelId)

	go func() {
		defer close(response.Channel)

		stream := c.Client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			MaxTokens: c.Config.MaxTokens,
			Messages:  messages,
			Model:     anthropic.Model(c.Config.ModelId),
		})
		defer stream.Close()

		for stream.Next() {
			event := stream.Current()
			delta := response.ProcessEvent(event)
			if delta != "" && onDelta != nil {
				onDelta(delta)
			}
		}
	}()

	return response, nil
}

// StreamSimpleMessage is a convenience method for sending a single text message
func (c *AnthropicClient) StreamSimpleMessage(ctx context.Context, text string, printToConsole bool) (*StreamingResponse, error) {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(text)),
	}

	var onDelta func(string)
	if printToConsole {
		onDelta = func(delta string) {
			//fmt.Print(delta)
		}
	}

	response, err := c.StreamMessages(ctx, messages, onDelta)
	if err != nil {
		return nil, err
	}

	if printToConsole {
		// fmt.Println() // newline at the end
	}

	return response, nil
}
