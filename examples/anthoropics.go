package examples

import (
	"context"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go" // imported as anthropic
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yuki5155/go-strands-agents/models"
)

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

func simpleStreamCallWithSchema(config *models.AnthropicConfig) *models.StreamingResponse {
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

	response := models.NewStreamingResponse(config.ModelId)

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

// Example using the new AnthropicClient struct - Simplest usage with defaults
func streamingWithClient() *models.StreamingResponse {
	// Create client with defaults (uses env var for API key, default model and max tokens)
	client := models.NewAnthropicClient()

	// Simple streaming call with console output
	response, err := client.StreamSimpleMessage(
		context.TODO(),
		"What is a quaternion?",
		true, // print to console
	)
	if err != nil {
		panic(err.Error())
	}

	return response
}

// Example with custom options
func streamingWithClientCustom() *models.StreamingResponse {
	// Create client with custom options
	client := models.NewAnthropicClient(
		models.WithMaxTokens(2048),
		models.WithModelId("claude-sonnet-4-5-20250929"),
	)

	// Stream with custom message and custom delta handler
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("Explain quantum entanglement in simple terms.")),
	}

	response, err := client.StreamMessages(context.TODO(), messages, func(delta string) {
		// Custom handling of each text delta
		fmt.Print(delta)
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println()

	return response
}
