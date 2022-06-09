package textarea

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	textarea := newTextArea()
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

func TestInput(t *testing.T) {
	textarea := newTextArea()

	input := "foo"

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	if !strings.Contains(view, input) {
		t.Log(view)
		t.Error("Textarea did not render the input")
	}

	if textarea.col != len(input) {
		t.Log(view)
		t.Error("Textarea did not move the cursor to the correct position")
	}
}

func TestWrap(t *testing.T) {
	textarea := newTextArea()
	textarea.Width = 5
	textarea.LineLimit = 5
	textarea.Height = 5
	textarea.CharLimit = 60

	textarea, _ = textarea.Update(initialBlinkMsg{})

	input := "foo bar baz"

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	for _, word := range strings.Split(input, " ") {
		if !strings.Contains(view, word) {
			t.Log(view)
			t.Error("Textarea did not render the input")
		}
	}

	// Due to the word wrapping, each word will be on a new line and the
	// textarea will look like this:
	//
	// > foo
	// > bar
	// > baz█
	if textarea.row != 2 && textarea.col != 3 {
		t.Log(view)
		t.Error("Textarea did not move the cursor to the correct position")
	}
}

func TestLineNumbers(t *testing.T) {
	textarea := newTextArea()
	textarea.ShowLineNumbers = true

	lines := 5

	textarea.LineLimit = lines
	textarea.Height = lines

	textarea, _ = textarea.Update(initialBlinkMsg{})

	view := textarea.View()

	for i := 0; i < lines; i++ {
		if !strings.Contains(view, fmt.Sprint(i+1)) {
			t.Log(view)
			t.Error("Textarea did not render the line numbers")
		}
	}
}

func TestCharLimit(t *testing.T) {
	textarea := newTextArea()

	// First input (foo bar) should be accepted as it will fall within the
	// CharLimit. Second input (baz) should not appear in the input.
	input := []string{"foo bar", "baz"}
	textarea.CharLimit = len(input[0])

	for _, k := range []rune(strings.Join(input, " ")) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()
	if strings.Contains(view, input[1]) {
		t.Log(view)
		t.Error("Textarea should not include input past the character limit")
	}
}

func TestVerticalScrolling(t *testing.T) {
	textarea := newTextArea()

	// Since Height is 1 and LineLimit is 5, the text area should vertically
	// scroll when the input is longer than the height.
	textarea.LineLimit = 5
	textarea.Height = 1
	textarea.Width = 20
	textarea.CharLimit = 100

	textarea, _ = textarea.Update(initialBlinkMsg{})

	input := "This is a really long line that should wrap around the text area."

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	// The view should contain the first "line" of the input.
	if !strings.Contains(view, "This is a really") {
		t.Log(view)
		t.Error("Textarea did not render the input")
	}

	// But we should be able to scroll to see the next line.
	// Let's scroll down for each line to view the full input.
	lines := []string{
		"long line that",
		"should wrap around",
		"the text area.",
	}
	for _, line := range lines {
		textarea.viewport.LineDown(1)
		view = textarea.View()
		if !strings.Contains(view, line) {
			t.Log(view)
			t.Error("Textarea did not render the correct scrolled input")
		}
	}
}

func TestScrollBehaviors(t *testing.T) {
	textarea := newTextArea()

	textarea.LineLimit = 20
	textarea.Height = 5
	textarea.Width = 8
	textarea.CharLimit = 200

	textarea.ScrollBehavior = ScrollOverflow

	textarea, _ = textarea.Update(initialBlinkMsg{})

	input := "Line 1 Line 2 Line 3 Line 4 Line 5 Line 6 Line 7 Line 8 Line 9"

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
		textarea.View()
	}

	// Typing these 9 lines should cause the text area to scroll. The text area
	// should have an offset of 4 lines to display the 9th line as the viewport
	// is 5 lines tall.

	if textarea.viewport.YOffset != 4 {
		t.Log(textarea.View())
		t.Log(textarea.row)
		t.Log(textarea.viewport.YOffset)
		t.Error("Textarea did not scroll down one line")
	}

	// Now let's scroll up.
	oldRow := textarea.row
	oldOffset := textarea.viewport.YOffset

	textarea.lineUp(2)
	textarea.View()

	// The cursor should be two lines higher but the viewport should be the same.
	// Since it is still contained in the window
	if textarea.row != oldRow-2 || textarea.viewport.YOffset != oldOffset {
		t.Log(textarea.View())
		t.Log(textarea.row)
		t.Log(textarea.viewport.YOffset)
		t.Error("Textarea did not scroll up two lines or did not maintain the scroll offset")
	}

	// Let's scroll up 3 more lines. This time the cursor should be 3 lines
	// higher and the viewport should be 2 lines higher since it will
	// underflow.
	oldRow = textarea.row
	oldOffset = textarea.viewport.YOffset

	textarea.lineUp(4)
	textarea.View()

	if textarea.row != oldRow-4 || textarea.viewport.YOffset != oldOffset-2 {
		t.Log(textarea.View())
		t.Log("Row: ", oldRow, textarea.row)
		t.Log("YOffset: ", oldOffset, textarea.viewport.YOffset)
		t.Error("Textarea did not scroll up three lines or did not shift the viewport")
	}
}

func newTextArea() Model {
	textarea := New()

	textarea.Prompt = "> "
	textarea.Placeholder = "Hello, World!"

	textarea.Focus()

	textarea, _ = textarea.Update(initialBlinkMsg{})

	return textarea
}

func keyPress(key rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}, Alt: false}
}
