package llm

import "testing"

func TestUsageFromResponseMap(t *testing.T) {
	// No usage object → nil.
	if usageFromResponseMap(map[string]interface{}{}) != nil {
		t.Error("missing usage should yield nil")
	}

	// Anthropic-style: input/output plus the prompt-cache breakdown. JSON
	// numbers decode to float64, which the parser must coerce to int.
	u := usageFromResponseMap(map[string]interface{}{
		"usage": map[string]interface{}{
			"input_tokens":                float64(6341),
			"output_tokens":               float64(8),
			"cache_read_input_tokens":     float64(6163),
			"cache_creation_input_tokens": float64(0),
		},
	})
	if u == nil {
		t.Fatal("expected usage")
	}
	if u.InputTokens != 6341 || u.OutputTokens != 8 || u.CacheReadInputTokens != 6163 {
		t.Errorf("anthropic parse wrong: %+v", u)
	}

	// OpenAI-style: prompt/completion/total.
	u2 := usageFromResponseMap(map[string]interface{}{
		"usage": map[string]interface{}{
			"prompt_tokens":     float64(100),
			"completion_tokens": float64(20),
			"total_tokens":      float64(120),
		},
	})
	if u2 == nil {
		t.Fatal("expected usage")
	}
	if u2.PromptTokens != 100 || u2.CompletionTokens != 20 || u2.TotalTokens != 120 {
		t.Errorf("openai parse wrong: %+v", u2)
	}
}
