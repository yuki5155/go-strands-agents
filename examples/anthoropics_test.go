package examples

import (
	"fmt"
	"testing"

	"github.com/yuki5155/go-strands-agents/models"
	"github.com/yuki5155/go-strands-agents/utils"
)

func init() {
	// Load .env file from the root directory
	utils.LoadEnvs()
}

func TestSampleAnthoropics(t *testing.T) {
	sampleAnthoropics()
}

func TestSampleAnthoropics2(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	simpleCall()
}

func TestNewAnthropicConfigWithDefaults(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	config := models.NewAnthropicConfig(models.WithApiKey(""))

	simpleCallWithDefaultsConfig(config)
}

func TestSimpleStreamCallWithDefaultsConfig(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	config := models.NewAnthropicConfig(models.WithApiKey(""))
	simpleStreamCallWithDefaultsConfig(config)
}

func TestSimpleStreamCallWithSchema(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	config := models.NewAnthropicConfig(models.WithApiKey(""))
	response := simpleStreamCallWithSchema(config)

	fmt.Printf("\n--------------------------------\n")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Role: %s\n", response.Role)
	fmt.Printf("Content Block Type: %s\n", response.ContentBlockType)
	fmt.Printf("Content Block Index: %d\n", response.ContentBlockIndex)
	fmt.Printf("Content: %s\n", response.Content)
	fmt.Printf("Stop Reason: %s\n", response.StopReason)
	fmt.Printf("Stop Sequence: %s\n", response.StopSequence)
	fmt.Printf("Input Tokens: %d\n", response.InputTokens)
	fmt.Printf("Output Tokens: %d\n", response.OutputTokens)
	fmt.Printf("Cache Creation Input Tokens: %d\n", response.CacheCreationInputTokens)
	fmt.Printf("Cache Read Input Tokens: %d\n", response.CacheReadInputTokens)
	fmt.Printf("--------------------------------\n")
}

func TestStreamingWithClient(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	response := streamingWithClient()

	fmt.Printf("\n--------------------------------\n")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Role: %s\n", response.Role)
	fmt.Printf("Content Block Type: %s\n", response.ContentBlockType)
	fmt.Printf("Content Block Index: %d\n", response.ContentBlockIndex)
	fmt.Printf("Content: %s\n", response.Content)
	fmt.Printf("Stop Reason: %s\n", response.StopReason)
	fmt.Printf("Stop Sequence: %s\n", response.StopSequence)
	fmt.Printf("Input Tokens: %d\n", response.InputTokens)
	fmt.Printf("Output Tokens: %d\n", response.OutputTokens)
	fmt.Printf("Cache Creation Input Tokens: %d\n", response.CacheCreationInputTokens)
	fmt.Printf("Cache Read Input Tokens: %d\n", response.CacheReadInputTokens)
	fmt.Printf("--------------------------------\n")
}

func TestStreamingWithClientCustom(t *testing.T) {
	_, ok := getApiKeyFromEnv()
	if !ok {
		t.Skip("Skipping test: ANTHROPIC_API_KEY is not set")
	}
	response := streamingWithClientCustom()
	for delta := range response.GetChannel() {
		fmt.Print(delta)
	}

	fmt.Printf("\n--------------------------------\n")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Model: %s\n", response.Model)
	fmt.Printf("Role: %s\n", response.Role)
	fmt.Printf("Content Block Type: %s\n", response.ContentBlockType)
	fmt.Printf("Content Block Index: %d\n", response.ContentBlockIndex)
	fmt.Printf("Content: %s\n", response.Content)
	fmt.Printf("Stop Reason: %s\n", response.StopReason)
	fmt.Printf("Stop Sequence: %s\n", response.StopSequence)
	fmt.Printf("Input Tokens: %d\n", response.InputTokens)
	fmt.Printf("Output Tokens: %d\n", response.OutputTokens)
	fmt.Printf("Cache Creation Input Tokens: %d\n", response.CacheCreationInputTokens)
	fmt.Printf("Cache Read Input Tokens: %d\n", response.CacheReadInputTokens)
	fmt.Printf("--------------------------------\n")
}
