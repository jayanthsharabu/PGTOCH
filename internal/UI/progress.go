package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressModel struct {
	progress      progress.Model
	total         int
	current       int
	operationName string
	status        string
	width         int
}

func NewProgressModel(operationName string, total int) *ProgressModel {

	p := progress.New(progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return &ProgressModel{
		progress:      p,
		total:         total,
		current:       0,
		operationName: operationName,
		status:        "Initializing...",
		width:         80,
	}
}

func (m *ProgressModel) incrementProgress(amount int, status string) tea.Cmd {
	m.current += amount
	if m.current > m.total {
		m.current = m.total
	}
	m.status = status
	return nil
}

func (m *ProgressModel) SetProgress(value int, status string) tea.Cmd {
	m.current = value
	if m.current > m.total {
		m.current = m.total
	}
	m.status = status
	return nil
}

func (m *ProgressModel) Init() tea.Cmd {
	return nil
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = m.width - 20
	}

	progressModel, cmd := m.progress.Update(msg)
	m.progress = progressModel.(progress.Model)
	return m, cmd

}

func (m *ProgressModel) View() string {
	percent := float64(m.current) / float64(m.total)
	if m.total == 0 {
		percent = 0
	}

	pad := strings.Repeat(" ", 2)
	progressBar := m.progress.ViewAs(percent)

	title := HighlightStyle.Render(m.operationName)

	stats := InfoStyle.Render(fmt.Sprintf("%d/%d", m.current, m.total))

	status := InfoStyle.Render(m.status)

	return lipgloss.JoinVertical(lipgloss.Left, title, progressBar, pad+stats, pad+status)

}
