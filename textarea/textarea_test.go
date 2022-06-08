package textarea

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextarea(t *testing.T) {
	textarea := New()

	textarea.Focus()

	textarea.Prompt = "> "
	textarea.Placeholder = "Hello, World!"

	textarea, _ = textarea.Update(initialBlinkMsg{})
	view := textarea.View()

	if !strings.Contains(view, ">") {
		t.Log(view)
		t.Error("Textarea did not render the prompt")
	}

	if !strings.Contains(view, "World!") {
		t.Log(view)
		t.Error("Textarea did not render the placeholder")
	}
}

func TestTextareaInput(t *testing.T) {
	textarea := New()

	textarea.Prompt = "> "
	textarea.Placeholder = "Hello, World!"
	textarea.CharLimit = 100
	textarea.Height = 5
	textarea.LineLimit = 5
	textarea.Width = 10

	textarea.Focus()

	textarea, _ = textarea.Update(initialBlinkMsg{})

	for _, k := range []rune("foo") {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	input := "foo"

	if !strings.Contains(view, input) {
		t.Log(view)
		t.Error("Textarea did not render the input")
	}

	if textarea.col != len(input) {
		t.Log(view)
		t.Error("Textarea did not move the cursor to the correct position")
	}
}

func TestTextareaWrap(t *testing.T) {
	 
}

func keyPress(key rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}, Alt: false}
}
