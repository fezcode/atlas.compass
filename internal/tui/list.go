package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fezcode/atlas.compass/pkg/model"
)

// item implements list.Item
type item struct {
	entry model.Entry
}

func (i item) Title() string       { return i.entry.Title }
func (i item) Description() string { return i.entry.Username }
func (i item) FilterValue() string { return i.entry.Title + " " + i.entry.Username + " " + i.entry.URL }

type ListModel struct {
	List list.Model
}

func NewListModel(entries []model.Entry, width, height int) ListModel {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = item{entry: e}
	}

	// Use default delegate but customize styles later if needed
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(ColorPrimary).BorderLeftForeground(ColorPrimary)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(ColorPrimary)

				l := list.New(items, d, width, height)

				l.Title = "Compass Vault"

					l.Styles.Title = StyleListHeader

					

					// Force the help styles to be bright white

					l.Styles.HelpStyle = StyleSubtext

					l.Help.Styles.ShortKey = StyleSubtext

					l.Help.Styles.ShortDesc = StyleSubtext

					l.Help.Styles.ShortSeparator = StyleSubtext

					l.Help.Styles.FullKey = StyleSubtext

					l.Help.Styles.FullDesc = StyleSubtext

					

					l.SetShowHelp(false)

					l.KeyMap.Quit.SetKeys("q") // Explicitly remove 'esc' from the list's own quit keys

				

			

		

	 // We show our own hint, but ? toggles the big one
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy password")),
			key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "copy username")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	return ListModel{List: l}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.List.View()
}