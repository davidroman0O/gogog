package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
	"github.com/spf13/cobra"
)

/// GoGog CLI immersive environment which helps to manage the backups from the agent

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gogog",
		Short: "GoG backup manager",
		Long:  `Manager of your GoG games backups`,
		Run: func(cmd *cobra.Command, args []string) {
			if !data.HasState[types.GogAuth]() {
				fmt.Println("please use `gogog login` and get your cookies to use `gogog`")
				return
			}
			var err error
			var state tea.Model
			init := InitialModel()
			p := tea.NewProgram(init)
			if state, err = p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
			var ok bool
			if init, ok = state.(CliGog); !ok {
				panic("failed to analyze state")
			}
			fmt.Println(init)
		},
	}
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type CliGog struct {
	state *types.Gogog
	table table.Model
}

func InitialModel() CliGog {

	return CliGog{
		table: table.New(
			table.WithColumns([]table.Column{
				{Title: "title", Width: 4},
				{Title: "genre", Width: 4},
				{Title: "year", Width: 4},
				{Title: "status", Width: 4},
			}),
			table.WithRows([]table.Row{
				{"game", "action, guns", "2017", "-"},
				{"game", "action, guns", "2017", "-"},
				{"game", "action, guns", "2017", "-"},
				{"game", "action, guns", "2017", "-"},
				{"game", "action, guns", "2017", "-"},
			}),
			table.WithFocused(true),
			table.WithHeight(7),
			table.WithWidth(100),
		),
	}
}

func (m CliGog) Init() tea.Cmd {
	return nil
}

func (m CliGog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
			return m, tea.Batch(
				tea.Printf("focused %v!", m.table.Focused()),
			)
		// case "up":
		// 	m.table.GotoTop()
		// case "down":
		// 	m.table.GotoBottom()
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m CliGog) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
