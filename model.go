package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int

	choices []string
	cursor  int

	selected map[string]*connect_struct

	current *connect_struct
	result  *link_data
}

func initialModel(connections []connect_struct) model {

	connectMap := make(map[string]*connect_struct)
	var names []string

	for i := range connections {
		connectMap[connections[i].name] = &connections[i]
		names = append(names, connections[i].name)
	}

	var current *connect_struct

	if len(connections) > 0 {
		current = &connections[0]
	}

	return model{
		choices:  names,
		selected: connectMap,
		current:  current,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:

		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.current = m.selected[m.choices[m.cursor]]
			}

		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				m.current = m.selected[m.choices[m.cursor]]
			}

		case "enter":
			if len(m.choices) > 0 {
				m.current = m.selected[m.choices[m.cursor]]
			}

		case "ctrl+g":

			if m.current != nil {
				result, found := get_link_data(m.current)

				fmt.Printf("FOUND: %v\n", found)
				fmt.Printf("RESULT: %+v\n", result)

				if found {
					m.result = &result
				}
			}

		case "ctrl+r":
			m.result = nil

		case "ctrl+s":

			if m.result != nil {
				save_link_data(m.result)
			}

		case "ctrl+h":
			// help later

		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		}
	}

	return m, nil
}
func (m model) View() string {

	leftWidth := 25

	rightWidth := m.width - leftWidth - 8
	if rightWidth < 40 {
		rightWidth = 40
	}

	var interfaceList string

	for i, choice := range m.choices {

		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		selected := " "

		if m.current != nil &&
			m.current.name == choice {
			selected = "*"
		}

		interfaceList += fmt.Sprintf(
			"%s [%s] %s\n",
			cursor,
			selected,
			choice,
		)
	}

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(15).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Interfaces\n\n" +
				interfaceList,
		)

	resultText := "Press Ctrl+G to get switch data"

	if m.result != nil {

		resultText = fmt.Sprintf(
			"Switch Name : %s\n\n"+
				"Port ID     : %s\n"+
				"VLAN        : %s",
			m.result.switch_name,
			m.result.port_id,
			m.result.vlan_id,
		)
	}

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(15).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Link Discovery\n\n" +
				resultText,
		)

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)

	connectionInfo := "No adapter selected"

	if m.current != nil {

		status := "Down"

		if m.current.status {
			status = "Up"
		}

		connectionInfo = fmt.Sprintf(
			"Connection : %s\n"+
				"Network    : %s\n"+
				"MAC Addr   : %s\n"+
				"IPv4 Addr  : %s\n"+
				"Status     : %s",
			m.current.name,
			m.current.network_card,
			m.current.mac_addr,
			m.current.ip_addr,
			status,
		)
	}

	connectionPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(8).
		Border(lipgloss.RoundedBorder()).
		Render(connectionInfo)

	resultsInfo := "No switch information available"

	if m.result != nil {

		resultsInfo = fmt.Sprintf(
			"Switch IP  : %s\n"+
				"Model      : %s\n"+
				"Duplex     : %s\n"+
				"VTP Domain : %s\n"+
				"Protocol   : %s",
			m.result.switch_ip,
			m.result.switch_model,
			m.result.duplex_option,
			m.result.vpt_mgmt_domain,
			m.result.protocol,
		)
	}

	resultsPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(8).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Additional Results\n\n" +
				resultsInfo,
		)

	bottom := lipgloss.JoinHorizontal(
		lipgloss.Top,
		connectionPane,
		resultsPane,
	)

	footer := lipgloss.NewStyle().
		Bold(true).
		Render(
			"^Q Quit    ^R Reset    ^G Get Link Data    ^S Save Link Data    ^H Help",
		)

	title := lipgloss.NewStyle().
		Bold(true).
		Render("LDLinux")

	return title +
		"\n\n" +
		top +
		"\n\n" +
		bottom +
		"\n\n" +
		footer +
		"\n"
}
