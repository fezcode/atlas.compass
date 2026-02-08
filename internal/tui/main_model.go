package tui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fezcode/atlas.compass/internal/store"
	"github.com/fezcode/atlas.compass/pkg/model"
)

type State int

const (
	StateAuth State = iota
	StateList
	StateDetail
	StateEditor
)

type MainModel struct {
	State          State
	Auth           AuthModel
	List           ListModel
	Detail         DetailModel
	Editor         EditorModel
	Vault          *model.Vault
	MasterPassword string
	WindowWidth    int
	WindowHeight   int
	StatusMsg      string
}

func NewMainModel() MainModel {
	return MainModel{
		State: StateAuth,
		Auth:  NewAuthModel(),
		List:  NewListModel([]model.Entry{}, 0, 0), // Initialize empty list to prevent crash on resize
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.Auth.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		// Update child models with size
		m.List.List.SetSize(msg.Width, msg.Height-4) // Reserve space for header/status
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.State {
	case StateAuth:
		// Handle Auth Logic
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				// Try to unlock
				pass := m.Auth.Input.Value()
				if pass == "" {
					m.Auth.Err = fmt.Errorf("password cannot be empty")
					return m, nil
				}
				
				vault, err := store.Load(pass)
				if err != nil {
					// Check if it's a decryption error vs file error
					// For now assume decryption error if file exists
					if store.Exists() {
						m.Auth.Err = fmt.Errorf("invalid password")
						m.Auth.Input.SetValue("")
						return m, nil
					}
					// New vault logic could go here, but Load returns empty vault if not exists
					// So actually if err != nil here it's a real error (corruption/permission)
					m.Auth.Err = err
					return m, nil
				}

				m.Vault = vault
				m.MasterPassword = pass
				m.State = StateList
				m.List = NewListModel(vault.Entries, m.WindowWidth, m.WindowHeight-4)
				return m, nil
			}
		}
		
		var authCmd tea.Cmd
		m.Auth, authCmd = m.Auth.Update(msg)
		cmds = append(cmds, authCmd)

	case StateList:
		// Handle List Logic
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// If filtering, let the list handle it
			if m.List.List.FilterState() == list.Filtering {
				break
			}

			switch msg.String() {
			case "?":
				m.List.List.SetShowHelp(!m.List.List.ShowHelp())
				return m, nil
			case "q":
				return m, tea.Quit
			case "a":
				m.State = StateEditor
				m.Editor = NewEditorModel()
				return m, tea.Batch(m.Editor.Init())
			case "enter":
				// View details
				if item, ok := m.List.List.SelectedItem().(item); ok {
					m.State = StateDetail
					m.Detail = DetailModel{Entry: item.entry}
					return m, nil
				}
			case "e":
				// Direct Edit
				if item, ok := m.List.List.SelectedItem().(item); ok {
					m.State = StateEditor
					m.Editor = NewEditorModel()
					m.Editor.SetEntry(item.entry)
					return m, m.Editor.Init()
				}
			case "c":
				// Copy Password
				if item, ok := m.List.List.SelectedItem().(item); ok {
					clipboard.WriteAll(item.entry.Password)
					m.StatusMsg = "Password copied to clipboard!"
					return m, m.clearStatusAfter(2 * time.Second)
				}
			case "u":
				// Copy Username
				if item, ok := m.List.List.SelectedItem().(item); ok {
					clipboard.WriteAll(item.entry.Username)
					m.StatusMsg = "Username copied to clipboard!"
					return m, m.clearStatusAfter(2 * time.Second)
				}
			case "d":
				// Delete
				if index := m.List.List.Index(); index >= 0 && index < len(m.Vault.Entries) {
					// We need to match the item from list to vault
					if item, ok := m.List.List.SelectedItem().(item); ok {
						m.deleteEntry(item.entry.ID)
						m.saveVault()
						m.refreshList()
						m.StatusMsg = "Entry deleted."
						return m, m.clearStatusAfter(2 * time.Second)
					}
				}
			}
		}

		var listCmd tea.Cmd
		m.List, listCmd = m.List.Update(msg)
		cmds = append(cmds, listCmd)

	case StateDetail:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "backspace":
				m.State = StateList
				return m, nil
			case "e":
				m.State = StateEditor
				m.Editor = NewEditorModel()
				m.Editor.SetEntry(m.Detail.Entry)
				return m, m.Editor.Init()
			case "c":
				clipboard.WriteAll(m.Detail.Entry.Password)
				m.StatusMsg = "Password copied!"
				return m, m.clearStatusAfter(2 * time.Second)
			case "u":
				clipboard.WriteAll(m.Detail.Entry.Username)
				m.StatusMsg = "Username copied!"
				return m, m.clearStatusAfter(2 * time.Second)
			}
		}

	case StateEditor:
		// Handle Editor Logic
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEsc {
				m.State = StateList
				return m, nil
			}
			if msg.Type == tea.KeyEnter && m.Editor.Focused == FieldCount-1 {
				// Save
				newEntry := m.Editor.GetEntry()
				if newEntry.ID == "" {
					// Create new
					newEntry.ID = generateID()
					newEntry.CreatedAt = time.Now()
					newEntry.UpdatedAt = time.Now()
					m.Vault.Entries = append(m.Vault.Entries, newEntry)
				} else {
					// Update existing
					for i, e := range m.Vault.Entries {
						if e.ID == newEntry.ID {
							newEntry.UpdatedAt = time.Now()
							m.Vault.Entries[i] = newEntry
							break
						}
					}
				}
				m.saveVault()
				m.refreshList()
				m.State = StateList
				m.StatusMsg = "Entry saved."
				return m, m.clearStatusAfter(2 * time.Second)
			}
		}

		var editorCmd tea.Cmd
		m.Editor, editorCmd = m.Editor.Update(msg)
		cmds = append(cmds, editorCmd)
	}

	// Global status clear
	switch msg.(type) {
	case clearStatusMsg:
		m.StatusMsg = ""
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.State {
	case StateAuth:
		// Center auth box
		return lipgloss.Place(
			m.WindowWidth, m.WindowHeight,
			lipgloss.Center, lipgloss.Center,
			m.Auth.View(),
		)
	case StateList:
		view := m.List.View()
		helpHint := StyleSubtext.Render(" [a] add • [enter] view • [e] edit • [c] copy pass • [u] copy user • [d] delete • [?] help")
		if m.StatusMsg != "" {
			status := StyleStatusBar.Render("❯ " + m.StatusMsg)
			view = lipgloss.JoinVertical(lipgloss.Left, view, status, helpHint)
		} else {
			view = lipgloss.JoinVertical(lipgloss.Left, view, "", helpHint)
		}
		return view
	case StateDetail:
		return lipgloss.Place(
			m.WindowWidth, m.WindowHeight,
			lipgloss.Center, lipgloss.Center,
			m.Detail.View(),
		)
	case StateEditor:
		return lipgloss.Place(
			m.WindowWidth, m.WindowHeight,
			lipgloss.Center, lipgloss.Center,
			m.Editor.View(),
		)
	}
	return ""
}

// Helpers

func (m *MainModel) saveVault() {
	if err := store.Save(m.Vault, m.MasterPassword); err != nil {
		m.StatusMsg = "Error saving vault: " + err.Error()
	}
}

func (m *MainModel) refreshList() {
	// Re-create list model with current entries
	m.List = NewListModel(m.Vault.Entries, m.WindowWidth, m.WindowHeight-4)
}

func (m *MainModel) deleteEntry(id string) {
	newEntries := []model.Entry{}
	for _, e := range m.Vault.Entries {
		if e.ID != id {
			newEntries = append(newEntries, e)
		}
	}
	m.Vault.Entries = newEntries
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Status clearing
type clearStatusMsg struct{}

func (m MainModel) clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}