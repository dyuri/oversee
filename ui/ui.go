package ui

import (
	"fmt"
	"time"

	"github.com/dyuri/oversee/proc"

	"github.com/ShinyTrinkets/overseer"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table table.Model
	overseer *overseer.Overseer
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#268bd2"))

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func updateProcesses (ovr *overseer.Overseer, t table.Model) table.Model {
	procs := ovr.ListAll()
	rows := make([]table.Row, len(procs))
	for i, p := range procs {
		pd := ovr.Status(p)
		rows[i] = table.Row{
			fmt.Sprintf("%d.", i + 1),
			pd.ID,
			pd.State,
			pd.Cmd,
			fmt.Sprintf("%d", pd.PID),
		}
	}
	t.SetRows(rows)

	return t
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "s":
			p := m.table.SelectedRow()[1]
			m.overseer.Supervise(p)
			return m, nil
		}
	case TickMsg:
		m.table = updateProcesses(m.overseer, m.table) // TODO update to command
		return m, doTick()
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View())
}

// TODO help

func StartUI() {
	ovr := proc.GetOverseer()

	columns := []table.Column{
		{Title: "Nr.", Width: 3},
		{Title: "Name", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "Command", Width: 20},
		{Title: "PID", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	t = updateProcesses(ovr, t)

	s := table.DefaultStyles()
	t.SetStyles(s)

	p := tea.NewProgram(model{
		overseer: ovr,
		table: t,
	})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
