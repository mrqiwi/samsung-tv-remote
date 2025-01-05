package tui

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
)

type CommandListModel struct {
    Commands   []string
    Cursor     int
    Selected   int
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
        }
    }
    return m, nil
}

func (m *CommandListModel) View() string {
    s := "Available commands:\n\n"

    for i, cmd := range m.Commands {
        cursor := " "
        if m.Cursor == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %d. %s\n", cursor, i+1, cmd)
    }

    return s
}

func (m *CommandListModel) SelectedCommand() string {
    return m.Commands[m.Selected]
}