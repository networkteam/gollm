package llm

import "testing"

// A plain prompt must render its input once. NewPrompt keeps the input in both
// Input and Messages, and rendering both doubled every flattened request.
func TestPromptStringDoesNotEchoInput(t *testing.T) {
	p := NewPrompt("THE-INPUT", WithSystemPrompt("THE-SYSTEM", ""))
	got := p.String()

	if n := countOccurrences(got, "THE-INPUT"); n != 1 {
		t.Errorf("input rendered %d times, want 1:\n%s", n, got)
	}
	if countOccurrences(got, "THE-SYSTEM") != 1 {
		t.Errorf("system prompt missing or duplicated:\n%s", got)
	}

	// A real conversation still renders its turns.
	p2 := NewPrompt("FIRST")
	p2.Messages = append(p2.Messages, PromptMessage{Role: "assistant", Content: "REPLY"})
	if got2 := p2.String(); countOccurrences(got2, "REPLY") != 1 {
		t.Errorf("conversation turns must still render:\n%s", got2)
	}
}

func countOccurrences(haystack, needle string) int {
	n, from := 0, 0
	for {
		i := indexFrom(haystack, needle, from)
		if i < 0 {
			return n
		}
		n++
		from = i + len(needle)
	}
}

func indexFrom(s, sub string, from int) int {
	if from >= len(s) {
		return -1
	}
	for i := from; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
