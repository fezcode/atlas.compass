package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthModel struct {
	Input     textinput.Model
	Err       error
	IsLoading bool
}

func NewAuthModel() AuthModel {
	ti := textinput.New()
	ti.Placeholder = "Master Password"
	ti.EchoMode = textinput.EchoPassword
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return AuthModel{
		Input: ti,
	}
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(msg tea.Msg) (AuthModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m AuthModel) View() string {
	title := StyleAuthHeader.Render("ATLAS COMPASS")
	
	input := m.Input.View()
	
	hint := StyleSubtext.Render("Enter Master Password to Unlock")
	
	errView := ""
	if m.Err != nil {
		errView = lipgloss.NewStyle().Foreground(ColorError).MarginTop(1).Render(m.Err.Error())
	}

	content := lipgloss.JoinVertical(lipgloss.Center, title, hint, input, errView)
	
	return StyleAuthBox.Render(content)
}
