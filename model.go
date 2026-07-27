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

	errMessage string
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

				if found {
					m.result = &result
				} else {
					m.result = &link_data{
						switch_name: "NOT FOUND",
					}
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

	rightWidth := m.width - leftWidth - 6
	if rightWidth < 40 {
		rightWidth = 40
	}

	var interfaceList string

	for i, choice := range m.choices {

		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		interfaceList += fmt.Sprintf(
			"%s %s\n",
			cursor,
			choice,
		)
	}

	interfacesPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Interfaces\n\n" +
				interfaceList,
		)

	connectionText := "No adapter selected"

	if m.current != nil {

		status := "Down"

		if m.current.status {
			status = "Up"
		}

		connectionText = fmt.Sprintf(
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
		Width(rightWidth).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Connection\n\n" +
				connectionText,
		)

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		interfacesPane,
		connectionPane,
	)

	linkData := "Press Ctrl+G to collect LLDP/CDP data"

	if m.result != nil {

		linkData = fmt.Sprintf(
			"Switch Name : %s\n"+
				"Port ID     : %s\n"+
				"VLAN        : %s\n"+
				"Switch IP   : %s\n"+
				"Model       : %s\n"+
				"Duplex      : %s\n"+
				"VTP Domain  : %s\n"+
				"Protocol    : %s",
			m.result.switch_name,
			m.result.port_id,
			m.result.vlan_id,
			m.result.switch_ip,
			m.result.switch_model,
			m.result.duplex_option,
			m.result.vpt_mgmt_domain,
			m.result.protocol,
		)
	}

	linkPane := lipgloss.NewStyle().
		Width(m.width).
		Height(12).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Link Data\n\n" +
				linkData,
		)
	errorPane := ""
	if m.errMessage != "" {
		errorPane = lipgloss.NewStyle().
        	Foreground(lipgloss.Color("9")).
        	Border(lipgloss.RoundedBorder()).
        	Render("Error\n\n" + m.errMessage)
	}

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
		linkPane +
		"\n\n" +
		errorPane +
		"\n\n" +
		footer +
		"\n"
}
