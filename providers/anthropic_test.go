package providers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAnthropicEffortNestsUnderOutputConfig verifies that the flat "effort"
// option is sent as output_config.effort — the shape the Anthropic API
// expects — and never as an unknown top-level parameter (HTTP 400).
func TestAnthropicEffortNestsUnderOutputConfig(t *testing.T) {
	provider, ok := NewAnthropicProvider("fake-api-key", "claude-sonnet-5", nil).(*AnthropicProvider)
	assert.True(t, ok, "Provider should be of type *AnthropicProvider")

	body, err := provider.PrepareRequest("hello", map[string]interface{}{"effort": "low", "max_tokens": 4096})
	assert.NoError(t, err)

	var req map[string]interface{}
	assert.NoError(t, json.Unmarshal(body, &req))

	_, hasFlatEffort := req["effort"]
	assert.False(t, hasFlatEffort, "effort must not be sent top-level (Anthropic rejects unknown parameters)")

	outputConfig, ok := req["output_config"].(map[string]interface{})
	assert.True(t, ok, "output_config must be present")
	assert.Equal(t, "low", outputConfig["effort"], "effort level must be preserved under output_config")
}
