package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tasks          []string
	cursor         int
	selected       map[int]bool
	textInput      textinput.Model
	textInputFocus bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "New task"
	ti.CharLimit = 250
	ti.Width = 20

	return model{
		tasks:          []string{"testing", "testing", "testing"},
		selected:       make(map[int]bool),
		textInput:      ti,
		textInputFocus: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func mod(x, m int) int {
	return (x%m + m) % m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.textInputFocus {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				m.cursor--
				m.cursor = mod(m.cursor, len(m.tasks))

			case "down", "j":
				m.cursor++
				m.cursor = mod(m.cursor, len(m.tasks))

			case "enter", " ":
				task_done := m.selected[m.cursor]
				if task_done {
					m.selected[m.cursor] = false
				} else {
					m.selected[m.cursor] = true
				}

			case "o":
				m.textInputFocus = true
			}
		} else {
			switch msg.String() {
			case "enter":
                m.tasks = append(m.tasks, m.textInput.Value())
                m.textInputFocus = false
                m.textInput.Reset()
			default:
				m.textInput, cmd = m.textInput.Update(msg)
				m.textInput.Focus()
				return m, cmd
			}
		}
	}
	return m, cmd
}

func (m model) View() string {
	s := "Task list (type o to add task):\n\n"

	for i, choice := range m.tasks {
		cursor := " "
		if m.cursor == i && !m.textInputFocus {
			cursor = ">"
		}

		checked := " "
		if b := m.selected[i]; b {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	if m.textInputFocus {
		s += fmt.Sprintf("%s\n", m.textInput.View())
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
