package gitlab

import "time"

type Config struct {
	BaseURL   string
	Token     string
	ProjectID string
	GroupID   string
	Days      int
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type MergeRequest struct {
	ID        int        `json:"id"`
	IID       int        `json:"iid"`
	ProjectID int        `json:"project_id"`
	Title     string     `json:"title"`
	State     string     `json:"state"` // opened, closed, merged
	Author    User       `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`
	WebURL    string     `json:"web_url"`
}

type Event struct {
	ID         int       `json:"id"`
	ActionName string    `json:"action_name"`
	Author     User      `json:"author"`
	CreatedAt  time.Time `json:"created_at"`
	TargetType string    `json:"target_type"`
}

type DashboardData struct {
	ReviewsByUser map[string]int
	MRsByAuthor   map[string]int
	TotalMRs      int
	OpenMRs       int
	MergedMRs     int
	ClosedMRs     int
	FetchedAt     time.Time
	Since         time.Time
	Days          int
}
