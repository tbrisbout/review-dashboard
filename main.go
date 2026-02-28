package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"review-dashboard/internal/gitlab"
	"review-dashboard/internal/ui"
)

func main() {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: GITLAB_TOKEN environment variable is required")
		os.Exit(1)
	}

	projectID := os.Getenv("GITLAB_PROJECT_ID")
	groupID := os.Getenv("GITLAB_GROUP_ID")
	if projectID == "" && groupID == "" {
		fmt.Fprintln(os.Stderr, "error: GITLAB_PROJECT_ID or GITLAB_GROUP_ID environment variable is required")
		os.Exit(1)
	}

	cfg := gitlab.Config{
		BaseURL:   getEnv("GITLAB_URL", "https://gitlab.com"),
		Token:     token,
		ProjectID: projectID,
		GroupID:   groupID,
		Days:      getEnvInt("GITLAB_DAYS", 14),
	}

	client := gitlab.NewClient(cfg)
	model := ui.NewModel(client)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
