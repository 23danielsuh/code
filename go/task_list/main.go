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
	completed      map[int]bool
	textInput      textinput.Model
	textInputFocus bool
	editing        bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "New task"
	ti.CharLimit = 50
	ti.Prompt = ""
	ti.Width = 50
	ti.Focus()

	return model{
		tasks:          []string{},
		completed:      make(map[int]bool),
		textInput:      ti,
		textInputFocus: true,
		editing:        false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func mod(x, m int) int {
    if m == 0 {
        return 0
    }
	return (x%m + m) % m
}
func remove(slice []string, i int) []string {
    if len(slice) == 0 || len(slice) == 1 {
        return []string{}
    }
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
    
    /*
    TODO:
    deal with the case where task list is empty and you always want to add a task
    make it so that hitting enter edits the text not deletes all of it
    clean up code
    */

    if len(m.tasks) == 0 {
        m.textInputFocus = true
        m.textInput.Focus()
    }

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

			case " ":
				task_done := m.completed[m.cursor]
				if task_done {
					m.completed[m.cursor] = false
				} else {
					m.completed[m.cursor] = true
				}

			case "d":
				m.tasks = remove(m.tasks, m.cursor)
				m.cursor = mod(m.cursor, len(m.tasks))

			case "enter":
				m.editing = true
				m.textInputFocus = true
				m.textInput.Focus()

			case "o":
				m.textInputFocus = true
				m.textInput.Focus()
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

            case "esc":
				m.textInputFocus = false
                m.editing = false
				m.textInput.Reset()
                return m, nil

			case "enter":
				if m.editing {
					m.tasks[m.cursor] = m.textInput.Value()
					m.editing = false
				} else {
					m.tasks = append(m.tasks, m.textInput.Value())
					m.cursor = len(m.tasks) - 1
				}
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
	allCompleted := true

	s := "Task list:\n\n"

    for i, choice := range m.tasks {
        if m.editing && i == m.cursor {
            s += fmt.Sprintf("  [ ] %s\n", m.textInput.View())
            allCompleted = false
            continue
        }
        cursor := " "
        if m.cursor == i && !m.textInputFocus {
            cursor = ">"
        }

        checked := " "
        if b := m.completed[i]; b {
            checked = "x"
        } else {
            allCompleted = false
        }

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

	if m.textInputFocus && !m.editing {
		s += fmt.Sprintf("  [ ] %s\n", m.textInput.View())
	}
	if allCompleted && (!m.textInputFocus || m.editing) && len(m.tasks) > 0 {
		s += fmt.Sprintf("\nAll tasks completed ✅")
	} else {
		s += fmt.Sprintln()
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
