package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	disc "github.com/mrqiwi/samsung-tv-remote/internal/discover"
)

type DeviceListModel struct {
    Devices     []disc.DeviceInfo
    Cursor      int
    Selected    int
    Exit        bool
}

func NewDeviceListModel(devices []disc.DeviceInfo) DeviceListModel {
    return DeviceListModel{
        Devices: devices,
        Cursor:  0,
    }
}

func (m *DeviceListModel) Init() tea.Cmd {
    return nil
}

func (m *DeviceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            m.Cursor--
            if m.Cursor < 0 {
                m.Cursor = len(m.Devices) - 1
            }
        case "down":
            m.Cursor++
            if m.Cursor >= len(m.Devices) {
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

func (m *DeviceListModel) View() string {
    s := "Choose a device:\n\n"

    for i, device := range m.Devices {
        cursor := " "
        if m.Cursor == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %d. %s (%s)\n", cursor, i+1, device.Name, device.IPAddress)
    }

    return s
}

func (m *DeviceListModel) SelectedDevice() disc.DeviceInfo {
    return m.Devices[m.Selected]
}

func (m *DeviceListModel) Esc() bool {
    return m.Exit
}