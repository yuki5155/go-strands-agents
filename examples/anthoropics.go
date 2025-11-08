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
