package bubble

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	UserTextInputs UserInput
)

type UserInput struct {
	UserText  string
	PrevQuery string
}

func RephraseInput() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type textModel struct {
	textInput textinput.Model
	err       error
	quitting  *bool
}

func initialModel() textModel {
	ti := textinput.New()
	ti.Placeholder = "Rephrase your query"
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 200

	return textModel{
		textInput: ti,
		err:       nil,
		quitting:  &Quitting,
	}
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if *m.quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			UserTextInputs.UserText = m.textInput.Value()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			*m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textModel) View() string {
	if *m.quitting {
		return ""
	}
	return fmt.Sprintf(
		"%s\n\n",
		m.textInput.View()) + "\n"
}
