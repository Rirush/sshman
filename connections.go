package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

var (
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Dark: "0", Light: "15"}).
			Background(lipgloss.AdaptiveColor{Dark: "15", Light: "0"}).
			PaddingLeft(1)

	GroupStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Bold(true)

	ConnectionStyle = lipgloss.NewStyle().
			MarginLeft(3)

	NoConnectionsStyle = lipgloss.NewStyle().
				MarginLeft(3).
				Foreground(lipgloss.Color("8"))

	SelectedStyle = lipgloss.NewStyle().Underline(true)
)

var (
	SelectBinding = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))
)

type Size struct {
	Width, Height int
}

type Group struct {
	Name        string
	Connections []string
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                                 { return 1 }
func (d ItemDelegate) Spacing() int                                { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, cmd *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	switch item := listItem.(type) {
	case FlatGroup:
		if index == m.Index() {
			fmt.Fprint(w, GroupStyle.Copy().Inherit(SelectedStyle).Render(string(item)))
		} else {
			fmt.Fprint(w, GroupStyle.Render(string(item)))
		}
	case FlatConnection:
		if index == m.Index() {
			fmt.Fprint(w, ConnectionStyle.Copy().Inherit(SelectedStyle).Render(string(item)))
		} else {
			fmt.Fprint(w, ConnectionStyle.Render(string(item)))
		}
	case NoConnection:
		if index == m.Index() {
			fmt.Fprint(w, NoConnectionsStyle.Copy().Inherit(SelectedStyle).Render("(no connections)"))
		} else {
			fmt.Fprint(w, NoConnectionsStyle.Render("(no connections)"))
		}
	}
}

type FlatGroup string

func (f FlatGroup) FilterValue() string { return "" }

type FlatConnection string

func (f FlatConnection) FilterValue() string { return string(f) }

type NoConnection bool

func (f NoConnection) FilterValue() string { return "" }

type ConnectionsModel struct {
	Groups []Group
	Flat   []list.Item
	List   list.Model

	Select key.Binding
}

func NewConnectionsModel() ConnectionsModel {
	m := ConnectionsModel{}

	m.Groups = []Group{
		{"Group 1", []string{"Connection 1"}},
		{"Group 2", []string{"Connection 2", "Connection 3"}},
		{"Group 3", []string{"Connection 4"}},
		{"Group 4", nil},
		{"Group 5", []string{"Connection 5"}},
	}

	for _, group := range m.Groups {
		m.Flat = append(m.Flat, FlatGroup(group.Name))

		if len(group.Connections) == 0 {
			m.Flat = append(m.Flat, NoConnection(true))
		}

		for _, conn := range group.Connections {
			m.Flat = append(m.Flat, FlatConnection(conn))
		}
	}

	m.List = list.New(m.Flat, ItemDelegate{}, 0, 0)
	m.List.Title = "SSH Manager"
	m.List.SetShowStatusBar(false)

	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{SelectBinding}
	}

	m.List.Styles.Title = HeaderStyle
	m.List.Styles.TitleBar = lipgloss.NewStyle().
		MarginBottom(1)

	return m
}

func (c ConnectionsModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (c ConnectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.List.SetSize(msg.Width, msg.Height)
		c.List.Styles.Title.Width(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// do connection or whatever
		case "ctrl+c":
			fallthrough
		case "q":
			return c, tea.Quit
		}
	}

	newListModel, cmd := c.List.Update(msg)
	c.List = newListModel
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c ConnectionsModel) View() string {
	return c.List.View()
}
