# review-dashboard

A terminal dashboard for GitLab merge request activity, built with [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss).

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-blue)

## Features

- **Reviewers panel** — ranks contributors by number of approvals over the configured period
- **Authors panel** — ranks contributors by number of merge requests opened
- **Summary stats** — total, open, merged, and closed MR counts at a glance
- **Proportional bars** — visual `█░` bars for quick comparison
- Supports both **project** and **group** scopes
- Scrollable panels, live refresh, responsive layout

## Requirements

- Go 1.24+
- A GitLab personal access token with `read_api` scope

## Build

```bash
git clone https://github.com/your-username/review-dashboard.git
cd review-dashboard
go build -o review-dashboard .
```

Or run directly without building:

```bash
go run .
```

## Configuration

All configuration is done via environment variables.

| Variable | Required | Default | Description |
|---|---|---|---|
| `GITLAB_TOKEN` | Yes | — | Personal access token with `read_api` scope |
| `GITLAB_PROJECT_ID` | One of these | — | Numeric project ID or URL-encoded path (e.g. `mygroup%2Fmyproject`) |
| `GITLAB_GROUP_ID` | One of these | — | Numeric group ID or URL-encoded path (e.g. `mygroup`) |
| `GITLAB_URL` | No | `https://gitlab.com` | Base URL for self-hosted GitLab instances |
| `GITLAB_DAYS` | No | `14` | Number of days to look back |

## Run

```bash
export GITLAB_TOKEN=glpat-xxxxxxxxxxxx
export GITLAB_PROJECT_ID=12345

./review-dashboard
```

Self-hosted GitLab with a group scope over the last 30 days:

```bash
export GITLAB_TOKEN=glpat-xxxxxxxxxxxx
export GITLAB_URL=https://gitlab.mycompany.com
export GITLAB_GROUP_ID=my-team
export GITLAB_DAYS=30

./review-dashboard
```

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `r` | Refresh data |
| `Tab` / `Shift+Tab` | Switch active panel |
| `↑` / `k` | Scroll up in active panel |
| `↓` / `j` | Scroll down in active panel |
| `q` / `Ctrl+C` | Quit |

## How it works

- **Reviewers** are counted from GitLab's [events API](https://docs.gitlab.com/ee/api/events.html) filtered by `action=approved`, making it a single paginated request rather than one call per MR.
- **Authors** are counted from the [merge requests API](https://docs.gitlab.com/ee/api/merge_requests.html) filtered by `created_after`.
- Both requests are made concurrently and all pages are fetched automatically.

## License

MIT — see [LICENSE](LICENSE) for details.
