package examples

import (
	"context"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go" // imported as anthropic
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yuki5155/go-strands-agents/models"
)

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

func sampleAnthoropics() {
	fmt.Println("anthoropic")
}

func getApiKeyFromEnv() (string, bool) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", false
	}
	return apiKey, true
}
func simpleCall() {
	apiKey, ok := getApiKeyFromEnv()
	if !ok {
		panic("ANTHROPIC_API_KEY is not set")
	}
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.ModelClaudeSonnet4_5_20250929,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}

func simpleCallWithDefaultsConfig(config *models.AnthropicConfig) {
	client := anthropic.NewClient(
		option.WithAPIKey(config.ApiKey),
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: config.MaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.Model(config.ModelId),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}

func simpleStreamCallWithDefaultsConfig(config *models.AnthropicConfig) {
	client := anthropic.NewClient(
		option.WithAPIKey(config.ApiKey),
	)
	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: config.MaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.Model(config.ModelId),
	})
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		fmt.Printf("%+v\n", event)
	}
	if err := stream.Err(); err != nil {
		panic(err.Error())
	}
}

func simpleStreamCallWithSchema(config *models.AnthropicConfig) *StreamingResponse {
	client := anthropic.NewClient(
		option.WithAPIKey(config.ApiKey),
	)
	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: config.MaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.Model(config.ModelId),
	})
	defer stream.Close()

	response := NewStreamingResponse(config.ModelId)

	for stream.Next() {
		event := stream.Current()
		delta := response.ProcessEvent(event)
		if delta != "" {
			fmt.Print(delta)
		}
	}
	if err := stream.Err(); err != nil {
		panic(err.Error())
	}
	fmt.Println() // newline at the end

	return response
}
