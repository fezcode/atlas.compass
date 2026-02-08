package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CPField int

const (
	CPFieldConfirm CPField = iota
	CPFieldCurrent
	CPFieldNew
	CPFieldNewConfirm
	CPFieldCount
)

type ChangePassModel struct {
	Inputs  []textinput.Model
	Focused CPField
}

func NewChangePassModel() ChangePassModel {
	inputs := make([]textinput.Model, CPFieldCount)

	inputs[CPFieldConfirm] = textinput.New()
	inputs[CPFieldConfirm].Placeholder = "Type YES"
	inputs[CPFieldConfirm].Focus()
	inputs[CPFieldConfirm].CharLimit = 3

	inputs[CPFieldCurrent] = textinput.New()
	inputs[CPFieldCurrent].Placeholder = "Current Password"
	inputs[CPFieldCurrent].EchoMode = textinput.EchoPassword

	inputs[CPFieldNew] = textinput.New()
	inputs[CPFieldNew].Placeholder = "New Password"
	inputs[CPFieldNew].EchoMode = textinput.EchoPassword

	inputs[CPFieldNewConfirm] = textinput.New()
	inputs[CPFieldNewConfirm].Placeholder = "Confirm New Password"
	inputs[CPFieldNewConfirm].EchoMode = textinput.EchoPassword

	return ChangePassModel{
		Inputs:  inputs,
		Focused: CPFieldConfirm,
	}
}

func (m ChangePassModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChangePassModel) Update(msg tea.Msg) (ChangePassModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did user press enter on last field?
			if s == "enter" && m.Focused == CPFieldCount-1 {
				return m, nil // Handled by parent
			}

			// Cycle focus
			if s == "up" || s == "shift+tab" {
				m.Focused--
			} else {
				m.Focused++
			}

			if m.Focused > CPFieldCount-1 {
				m.Focused = 0
			} else if m.Focused < 0 {
				m.Focused = CPFieldCount - 1
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= int(CPFieldCount-1); i++ {
				if i == int(m.Focused) {
					cmds[i] = m.Inputs[i].Focus()
					continue
				}
				m.Inputs[i].Blur()
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m *ChangePassModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m ChangePassModel) View() string {
	var b strings.Builder

	b.WriteString(StyleListHeader.Render("Change Master Password"))
	b.WriteString("\n\n")
	
	b.WriteString(StyleSubtext.Render("WARNING: If you proceed, your vault will be re-encrypted."))
	b.WriteString("\n")
	b.WriteString(StyleSubtext.Render("There is NO UNDO."))
	b.WriteString("\n\n")

	for i, input := range m.Inputs {
		label := ""
		switch CPField(i) {
		case CPFieldConfirm:
			label = "Confirm (YES)"
		case CPFieldCurrent:
			label = "Current Pass"
		case CPFieldNew:
			label = "New Pass"
		case CPFieldNewConfirm:
			label = "Retype New"
		}

		style := StyleEditorLabel.Copy().Width(14)
		if CPField(i) == m.Focused {
			style = style.Foreground(ColorPrimary)
		}
		
		b.WriteString(style.Render(label))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}

	b.WriteString(StyleSubtext.Render(" [tab] next • [enter] save • [esc] cancel"))

	return b.String()
}
