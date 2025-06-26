package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
.______     _______                 ___        ______  __    __  
|   _  \   /  _____|    ______ _____\  \      /      ||  |  |  | 
|  |_)  | |  |  __     |______|______\  \    |  ,----'|  |__|  | 
|   ___/  |  | |_ |     ______ ______ >  >   |  |     |   __   | 
|  |      |  |__| |    |______|______/  /    |   ----.|  |  |  | 
| _|       \______|                 /__/      \______||__|  |__| 
                                                                 
`

type AppModel struct {
	Width       int
	Height      int
	Spinner     spinner.Model
	StatusMsg   string
	IsLoading   bool
	CurrentView string
	Error       error
}

func NewAppModel() *AppModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	return &AppModel{
		Spinner:     s,
		StatusMsg:   "Ready..",
		IsLoading:   false,
		CurrentView: "main",
	}
}

func (m *AppModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m *AppModel) Update(msg tea.Msg) (*AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m *AppModel) View() string {
	if m.Width == 0 {
		return "Loading..."
	}
	var view string
	switch m.CurrentView {
	case "main":
		view = m.renderMainView()
	default:
		view = m.renderMainView()
	}

	footer := FooterStyle.Render("Press q to quit")

	return lipgloss.JoinVertical(lipgloss.Center, view, footer)
}

func (m AppModel) renderMainView() string {
	logo := LogoStyle.Render(logo)

	title := TitleStyle.Render("Chug: Blazing-Fast ETL Pipeline")
	subtitle := SubtitleStyle.Render("PostgreSQL to ClickHouse data transfer")

	var status string
	if m.IsLoading {
		status = fmt.Sprintf("%s %s", m.Spinner.View(), m.StatusMsg)
	} else if m.Error != nil {
		status = ErrorStyle.Render(fmt.Sprintf("Error: %s", m.Error.Error()))
	} else {
		status = InfoStyle.Render(m.StatusMsg)
	}

	helpText := BoxStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		HighlightStyle.Render("Available Commands:"),
		InfoStyle.Render("connect   - Test database connections"),
		InfoStyle.Render("ingest    - Transfer data from PostgreSQL to ClickHouse"),
		InfoStyle.Render("export    - Export data from ClickHouse to CSV"),
		InfoStyle.Render("sample-config - Generate a sample configuration file"),
	))

	currentTime := time.Now().Format("15:04:05")
	stats := BoxStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		HighlightStyle.Render("Stats:"),
		InfoStyle.Render(fmt.Sprintf("Current time: %s", currentTime)),
		InfoStyle.Render("Terminal size: "+fmt.Sprintf("%d×%d", m.Width, m.Height)),
	))

	boxes := lipgloss.JoinHorizontal(lipgloss.Top, helpText, stats)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		title,
		subtitle,
		"",
		boxes,
		"",
		status,
	)
}
