package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
	"github.com/spf13/cobra"
)

/// GoGog CLI immersive environment which helps to manage the backups from the agent

/// I don't think i'm going to do the immersive environment, but i'm going to make the api (no gin, no fiber) and webui (htmx) just to rush the project

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gogog",
		Short: "GoG backup manager",
		Long:  `Manager of your GoG games backups`,
		Run: func(cmd *cobra.Command, args []string) {
			if !data.Has[types.GogAuthenticationChrome]() {
				fmt.Println("please use `gogog login` and get your cookies to use `gogog`")
				return
			}
			var err error
			var state tea.Model
			init := InitialModel()
			p := tea.NewProgram(init, tea.WithAltScreen())
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
	keys  keyMap
	help  help.Model
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = keyMap{
	// Up: key.NewBinding(
	// 	key.WithKeys("up", "k"),
	// 	key.WithHelp("↑/k", "move up"),
	// ),
	// Down: key.NewBinding(
	// 	key.WithKeys("down", "j"),
	// 	key.WithHelp("↓/j", "move down"),
	// ),
	// Left: key.NewBinding(
	// 	key.WithKeys("left", "h"),
	// 	key.WithHelp("←/h", "move left"),
	// ),
	// Right: key.NewBinding(
	// 	key.WithKeys("right", "l"),
	// 	key.WithHelp("→/l", "move right"),
	// ),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func InitialModel() CliGog {

	return CliGog{
		keys: keys,
		help: help.New(),
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
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		m.table.SetWidth(msg.Width)
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
	return baseStyle.Render(m.table.View()) + "\n" + strings.Repeat("\n", 4) + m.help.View(m.keys)
}
