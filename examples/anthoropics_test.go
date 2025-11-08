package examples

import (
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
