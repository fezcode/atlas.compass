package tui

import (
	"strings"

	"github.com/fezcode/atlas.compass/pkg/model"
)

type DetailModel struct {
	Entry model.Entry
}

func (m DetailModel) View() string {
	var b strings.Builder

	b.WriteString(StyleListHeader.Render("Entry Details"))
	b.WriteString("\n\n")

	renderField := func(label, value string) {
		b.WriteString(StyleEditorLabel.Render(label + ":"))
		b.WriteString(" ")
		if value == "" {
			b.WriteString(StyleSubtext.Render("none"))
		} else {
			b.WriteString(StyleBase.Render(value))
		}
		b.WriteString("\n")
	}

	renderField("Title", m.Entry.Title)
	renderField("User", m.Entry.Username)
	renderField("Pass", m.Entry.Password)
	renderField("URL", m.Entry.URL)
	renderField("Notes", m.Entry.Notes)

	b.WriteString("\n")
	b.WriteString(StyleSubtext.Render(" [e] Edit • [c] Copy Pass • [u] Copy User • [esc] Back"))

	return b.String()
}