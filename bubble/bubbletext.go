package bubble

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	UserText string
)

func RephraseInput(placeholder string) {
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
}

func initialModel() textModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return textModel{
		textInput: ti,
		err:       nil,
	}
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			UserText = m.textInput.Value()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
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
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n",
		m.textInput.View()) + "\n"
}
