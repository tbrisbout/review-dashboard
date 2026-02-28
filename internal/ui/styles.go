package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary   = lipgloss.Color("#7C3AED")
	colorAccent    = lipgloss.Color("#06B6D4")
	colorSuccess   = lipgloss.Color("#10B981")
	colorWarning   = lipgloss.Color("#F59E0B")
	colorDanger    = lipgloss.Color("#EF4444")
	colorMuted     = lipgloss.Color("#6B7280")
	colorText      = lipgloss.Color("#F9FAFB")
	colorSubtle    = lipgloss.Color("#374151")
	colorHighlight = lipgloss.Color("#A78BFA")

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			Background(colorPrimary).
			Padding(0, 2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(0, 1)

	panelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			MarginBottom(0)

	rankStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Width(4)

	usernameStyle = lipgloss.NewStyle().
			Foreground(colorText)

	countStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHighlight)

	barFillStyle  = lipgloss.NewStyle().Foreground(colorPrimary)
	barEmptyStyle = lipgloss.NewStyle().Foreground(colorSubtle)

	statBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(0, 2).
			MarginRight(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	emptyStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)
)
