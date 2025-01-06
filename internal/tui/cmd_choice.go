package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommandListModel struct {
	Commands []string
	Cursor   int
	Selected int
	Exit     bool
}

func NewCommandListModel(commands []string) CommandListModel {
	return CommandListModel{
		Commands: commands,
		Cursor:   0,
	}
}

func (m *CommandListModel) Init() tea.Cmd {
	return nil
}

func (m *CommandListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Commands) - 1
			}
		case "down":
			m.Cursor++
			if m.Cursor >= len(m.Commands) {
				m.Cursor = 0
			}
		case "enter":
			m.Selected = m.Cursor
			return m, tea.Quit
		case "ctrl+c", "q":
			m.Exit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *CommandListModel) View() string {
	if len(m.Commands) == 0 {
		return "No commands available.\nPress q to quit."
	}

	titleStyle := lipgloss.NewStyle()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	commandStyle := lipgloss.NewStyle()

	s := titleStyle.Render("Choose a command:\n\n")
	for i, cmd := range m.Commands {
		cursor := "  "
		if m.Cursor == i {
			cursor = cursorStyle.Render(">")
		}
		s += fmt.Sprintf("%s %d. %s\n", cursor, i+1, commandStyle.Render(cmd))
	}

	s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Use ↑/↓ to navigate, Enter to select, q to quit.")
	return s
}

func (m *CommandListModel) SelectedCommand() string {
	return m.Commands[m.Selected]
}

func (m *CommandListModel) Esc() bool {
	return m.Exit
}
