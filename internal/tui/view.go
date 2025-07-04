package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/akrugru/sshi/internal/config"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	config.SSHConfig
}

func (i item) Title() string { return i.Host }
func (i item) FilterValue() string {
	return strings.ToLower(fmt.Sprintf("%s %s %s %s",
		i.User, i.Host, i.HostName, strings.Join(i.Tags, " ")))
}
func (i item) Description() string {
	host := i.HostName
	if host == "" {
		host = i.Host
	}
	if host == "" {
		host = "<unknown>"
	}
	port := i.Port
	if port == "" {
		port = "22"
	}
	user := i.User
	if user == "" {
		user = "<user>"
	}
	return fmt.Sprintf("%s@%s:%s (Tags: %s)", user, host, port, strings.Join(i.Tags, ", "))
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				selected = &i.SSHConfig
				return m, tea.Quit
			}
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetWidth(msg.Width - h)
		m.list.SetHeight(msg.Height - v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.list.FilterInput.View(), docStyle.Render(m.list.View()))
}

var selected *config.SSHConfig

func SelectHost(cfgs []config.SSHConfig) (*config.SSHConfig, error) {
	items := make([]list.Item, len(cfgs))
	for i, c := range cfgs {
		items[i] = item{c}
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Select SSH Connection"
	m.list.SetFilteringEnabled(true)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return nil, err
	}
	return selected, nil
}
