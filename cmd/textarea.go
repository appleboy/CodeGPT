package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
}

// initialPrompt initializes a textarea model with the given value.
// It sets the character limit to the length of the value plus 100,
// inserts the value into the textarea, and adjusts the width and height
// of the textarea based on the value's content. The textarea is then focused.
//
// Parameters:
//   - value: A string to be inserted into the textarea.
//
// Returns:
//   - model: A struct containing the initialized textarea and any error encountered.
func initialPrompt(value string) model {
	ti := textarea.New()
	// origin defaultCharLimit = 400
	ti.CharLimit = len(value) + 100
	ti.InsertString(value)

	maxWidth := 0
	for _, line := range strings.Split(value, "\n") {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	ti.SetWidth(maxWidth)
	ti.SetHeight(len(strings.Split(value, "\n")))
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"Please confirm the following commit message:\n\n%s\n\n%s",
		m.textarea.View(),
		"(Press Ctrl+C to continue.)",
	) + "\n\n"
}
