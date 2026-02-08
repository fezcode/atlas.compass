package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fezcode/atlas.compass/pkg/model"
)

type EditorField int

const (
	FieldTitle EditorField = iota
	FieldUsername
	FieldPassword
	FieldURL
	FieldNotes
	FieldCount
)

type EditorModel struct {
	Inputs  []textinput.Model
	Focused EditorField
	Entry   *model.Entry // nil if creating new
}

func NewEditorModel() EditorModel {
	inputs := make([]textinput.Model, FieldCount)

	inputs[FieldTitle] = textinput.New()
	inputs[FieldTitle].Placeholder = "Title (e.g. GitHub)"
	inputs[FieldTitle].Focus()

	inputs[FieldUsername] = textinput.New()
	inputs[FieldUsername].Placeholder = "Username / Email"

	inputs[FieldPassword] = textinput.New()
	inputs[FieldPassword].Placeholder = "Password"
	// Don't mask in editor so user can see what they type? 
	// Or maybe toggle? For now, show it.

	inputs[FieldURL] = textinput.New()
	inputs[FieldURL].Placeholder = "URL (optional)"

	inputs[FieldNotes] = textinput.New()
	inputs[FieldNotes].Placeholder = "Notes (optional)"

	return EditorModel{
		Inputs:  inputs,
		Focused: FieldTitle,
	}
}

func (m *EditorModel) SetEntry(e model.Entry) {
	m.Entry = &e
	m.Inputs[FieldTitle].SetValue(e.Title)
	m.Inputs[FieldUsername].SetValue(e.Username)
	m.Inputs[FieldPassword].SetValue(e.Password)
	m.Inputs[FieldURL].SetValue(e.URL)
	m.Inputs[FieldNotes].SetValue(e.Notes)
	m.Focused = FieldTitle
	m.Inputs[FieldTitle].Focus()
}

func (m EditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did user press enter on last field?
			if s == "enter" && m.Focused == FieldCount-1 {
				// Handled by parent
				return m, nil 
			}

			// Cycle focus
			if s == "up" || s == "shift+tab" {
				m.Focused--
			} else {
				m.Focused++
			}

			if m.Focused > FieldCount-1 {
				m.Focused = 0
			} else if m.Focused < 0 {
				m.Focused = FieldCount - 1
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= int(FieldCount-1); i++ {
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

func (m *EditorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m EditorModel) View() string {
	var b strings.Builder

	b.WriteString(StyleListHeader.Render("Entry Editor"))
	b.WriteString("\n\n")

	for i, input := range m.Inputs {
		label := ""
		switch EditorField(i) {
		case FieldTitle:
			label = "Title"
		case FieldUsername:
			label = "Username"
		case FieldPassword:
			label = "Password"
		case FieldURL:
			label = "URL"
		case FieldNotes:
			label = "Notes"
		}

		style := StyleEditorLabel
		if EditorField(i) == m.Focused {
			style = style.Foreground(ColorPrimary)
		}
		
		b.WriteString(style.Render(label))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}

	b.WriteString(StyleSubtext.Render("Press Tab/Enter to navigate. Enter on last field to save. Esc to cancel."))

	return b.String()
}

func (m EditorModel) GetEntry() model.Entry {
	e := model.Entry{
		Title:    m.Inputs[FieldTitle].Value(),
		Username: m.Inputs[FieldUsername].Value(),
		Password: m.Inputs[FieldPassword].Value(),
		URL:      m.Inputs[FieldURL].Value(),
		Notes:    m.Inputs[FieldNotes].Value(),
	}
	if m.Entry != nil {
		e.ID = m.Entry.ID
		e.CreatedAt = m.Entry.CreatedAt
	}
	return e
}
