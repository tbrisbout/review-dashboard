package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"review-dashboard/internal/gitlab"
)

type viewState int

const (
	stateLoading viewState = iota
	stateLoaded
	stateError
)

type dataMsg struct {
	data *gitlab.DashboardData
}

type errMsg struct {
	err error
}

type userCount struct {
	name  string
	count int
}

type Model struct {
	client        *gitlab.Client
	spinner       spinner.Model
	state         viewState
	data          *gitlab.DashboardData
	err           error
	width         int
	height        int
	scrollOffsets [2]int
	activePanel   int
}

func NewModel(client *gitlab.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	return Model{
		client:  client,
		spinner: s,
		state:   stateLoading,
		width:   120,
		height:  40,
	}
}

func fetchDataCmd(client *gitlab.Client) tea.Cmd {
	return func() tea.Msg {
		data, err := client.FetchDashboardData()
		if err != nil {
			return errMsg{err: err}
		}
		return dataMsg{data: data}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchDataCmd(m.client),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.state = stateLoading
			m.err = nil
			m.scrollOffsets = [2]int{}
			return m, tea.Batch(m.spinner.Tick, fetchDataCmd(m.client))
		case "tab", "shift+tab":
			m.activePanel = (m.activePanel + 1) % 2
			return m, nil
		case "up", "k":
			if m.scrollOffsets[m.activePanel] > 0 {
				m.scrollOffsets[m.activePanel]--
			}
			return m, nil
		case "down", "j":
			m.scrollOffsets[m.activePanel]++
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case dataMsg:
		m.state = stateLoaded
		m.data = msg.data
		m.scrollOffsets = [2]int{}
		return m, nil

	case errMsg:
		m.state = stateError
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.loadingView()
	case stateError:
		return m.errorView()
	case stateLoaded:
		return m.dashboardView()
	}
	return ""
}

func (m Model) loadingView() string {
	content := m.spinner.View() + "  Fetching GitLab data..."
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) errorView() string {
	content := errorStyle.Render("Error: "+m.err.Error()) + "\n\n" +
		hintStyle.Render("Press r to retry  •  q to quit")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) dashboardView() string {
	if m.data == nil {
		return "No data"
	}

	lines := []string{}

	// Header
	title := fmt.Sprintf("  GitLab MR Dashboard — Last %d Days  ", m.data.Days)
	header := headerStyle.Width(m.width).Render(title)
	lines = append(lines, header)

	subtitle := fmt.Sprintf(
		"Fetched %s  •  Since %s",
		m.data.FetchedAt.Format("2006-01-02 15:04:05"),
		m.data.Since.Format("2006-01-02"),
	)
	lines = append(lines, subtitleStyle.Render(subtitle))
	lines = append(lines, "")

	// Stats row
	lines = append(lines, m.renderStatsRow())
	lines = append(lines, "")

	// Panels
	lines = append(lines, m.renderPanels())

	// Footer
	footer := footerStyle.Render("r refresh  •  q quit  •  tab switch panel  •  ↑/↓ scroll")
	lines = append(lines, footer)

	return strings.Join(lines, "\n")
}

func (m Model) renderStatsRow() string {
	stats := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"Total MRs", fmt.Sprintf("%d", m.data.TotalMRs), colorText},
		{"Open", fmt.Sprintf("%d", m.data.OpenMRs), colorSuccess},
		{"Merged", fmt.Sprintf("%d", m.data.MergedMRs), colorAccent},
		{"Closed", fmt.Sprintf("%d", m.data.ClosedMRs), colorWarning},
		{"Reviewers", fmt.Sprintf("%d", len(m.data.ReviewsByUser)), colorHighlight},
		{"Authors", fmt.Sprintf("%d", len(m.data.MRsByAuthor)), colorHighlight},
	}

	var boxes []string
	for _, s := range stats {
		label := lipgloss.NewStyle().Foreground(colorMuted).Render(s.label)
		value := lipgloss.NewStyle().Bold(true).Foreground(s.color).Render(s.value)
		content := lipgloss.JoinVertical(lipgloss.Center, value, label)
		boxes = append(boxes, statBoxStyle.Render(content))
	}

	return "  " + lipgloss.JoinHorizontal(lipgloss.Center, boxes...)
}

func (m Model) renderPanels() string {
	// Reserve space for header (~5 lines), stats (~3 lines), footer (~1), spacing
	usedLines := 10
	available := m.height - usedLines
	if available < 5 {
		available = 5
	}

	// Panel inner content width
	totalWidth := m.width - 2 // left margin
	panelWidth := (totalWidth / 2) - 2
	if panelWidth < 32 {
		panelWidth = 32
	}
	innerWidth := panelWidth - 4 // account for border + padding

	left := m.renderReviewersPanel(innerWidth, available)
	right := m.renderAuthorsPanel(innerWidth, available)

	leftBox := panelStyle
	rightBox := panelStyle
	if m.activePanel == 0 {
		leftBox = panelActiveStyle
	} else {
		rightBox = panelActiveStyle
	}

	leftRendered := leftBox.Width(panelWidth).Render(left)
	rightRendered := rightBox.Width(panelWidth).Render(right)

	return "  " + lipgloss.JoinHorizontal(lipgloss.Top, leftRendered, "  ", rightRendered)
}

func (m Model) renderReviewersPanel(innerWidth, maxRows int) string {
	entries := sortedEntries(m.data.ReviewsByUser)
	title := panelTitleStyle.Render("Reviewers  (approvals)")
	divider := lipgloss.NewStyle().Foreground(colorSubtle).Render(strings.Repeat("─", innerWidth))

	if len(entries) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			divider,
			emptyStyle.Render("No approval events found"),
		)
	}

	rows := m.renderEntries(entries, innerWidth, maxRows-3, m.scrollOffsets[0])
	hint := m.scrollHint(len(entries), maxRows-3, m.scrollOffsets[0])

	parts := []string{title, divider}
	parts = append(parts, rows...)
	if hint != "" {
		parts = append(parts, hintStyle.Render(hint))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) renderAuthorsPanel(innerWidth, maxRows int) string {
	entries := sortedEntries(m.data.MRsByAuthor)
	title := panelTitleStyle.Render("Authors  (MRs created)")
	divider := lipgloss.NewStyle().Foreground(colorSubtle).Render(strings.Repeat("─", innerWidth))

	if len(entries) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			divider,
			emptyStyle.Render("No merge requests found"),
		)
	}

	rows := m.renderEntries(entries, innerWidth, maxRows-3, m.scrollOffsets[1])
	hint := m.scrollHint(len(entries), maxRows-3, m.scrollOffsets[1])

	parts := []string{title, divider}
	parts = append(parts, rows...)
	if hint != "" {
		parts = append(parts, hintStyle.Render(hint))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) renderEntries(entries []userCount, innerWidth, maxVisible, offset int) []string {
	if maxVisible < 1 {
		maxVisible = 1
	}

	// Clamp offset
	if offset > len(entries)-maxVisible {
		offset = len(entries) - maxVisible
	}
	if offset < 0 {
		offset = 0
	}

	end := offset + maxVisible
	if end > len(entries) {
		end = len(entries)
	}

	maxCount := entries[0].count
	// bar width: total - rank(4) - space(1) - username(variable) - space(1) - count(4) - space(1)
	// keep username max at 20, count right-aligned in 4 chars
	const (
		rankW  = 4
		countW = 4
		barW   = 10
		spaces = 4 // spaces between columns
	)
	nameW := innerWidth - rankW - countW - barW - spaces
	if nameW < 8 {
		nameW = 8
	}

	var rows []string
	for i, entry := range entries[offset:end] {
		rank := rankStyle.Render(fmt.Sprintf("%d.", offset+i+1))

		name := entry.name
		if len(name) > nameW {
			name = name[:nameW-1] + "…"
		}
		nameCol := usernameStyle.Width(nameW).Render(name)

		countStr := fmt.Sprintf("%d", entry.count)
		countCol := countStyle.Width(countW).Render(countStr)

		bar := renderBar(entry.count, maxCount, barW)

		row := rank + " " + nameCol + " " + countCol + " " + bar
		rows = append(rows, row)
	}
	return rows
}

func (m Model) scrollHint(total, visible, offset int) string {
	if total <= visible {
		return ""
	}
	end := offset + visible
	if end > total {
		end = total
	}
	return fmt.Sprintf("%d–%d of %d", offset+1, end, total)
}

func renderBar(value, max, width int) string {
	if max == 0 || width == 0 {
		return strings.Repeat(" ", width)
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	return barFillStyle.Render(strings.Repeat("█", filled)) +
		barEmptyStyle.Render(strings.Repeat("░", empty))
}

func sortedEntries(m map[string]int) []userCount {
	entries := make([]userCount, 0, len(m))
	for name, count := range m {
		entries = append(entries, userCount{name, count})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].name < entries[j].name
	})
	return entries
}
