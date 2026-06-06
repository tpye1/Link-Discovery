package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)



type model struct {
    choices []string
    options []string
    // text area
    cursor   int
    selected map[string]*connect_struct

}



func initialModel(connections []connect_struct) model {


	var connect_map = make(map[string]*connect_struct)

	var connect_names_arr []string

	for i:=0; i < len(connections); i++ {
		connect_map[connections[i].name] = &(connections[i])
		connect_names_arr = append(connect_names_arr, connections[i].name)
	}


	return model{
		choices: connect_names_arr,
		selected: make(map[string]*connect_struct),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

	case "ctrl+g":
		// Get the link data
		//return
	case "ctrl+s":
		// Save the link data

	case "ctrl+h":
		// Help

	case "ctrl+r":
		// Reset the program

	// These keys should exit the program.
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit

	}

    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}


func (m model) View() string {
	// The header
	s := "Welcome to LdLinux\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[choice]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
