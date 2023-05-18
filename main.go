package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ProcessManager *ProcessManager
}

func NewModel() Model {
	procMan := NewProcessManager()
	procMan.OnStartup()
	return Model{
		ProcessManager: procMan,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			for _, proc := range m.ProcessManager.Processes {
				err := m.ProcessManager.StopProcess(proc)
				if err != nil {
					fmt.Printf("ERROR STOPPING PROCESS: %#v\n", err)
				}
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	var procList []string
	for _, proc := range m.ProcessManager.Processes {
		procList = append(procList, fmt.Sprintf("| %d: %s |", proc.ProcInfo.PID, proc.Service.Name))
	}
	procListStr := strings.Join(procList, "\n")

	return procListStr + "\n\nPress q to quit"
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
