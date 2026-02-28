package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doGet(path string, params url.Values) ([]byte, string, error) {
	endpoint := fmt.Sprintf("%s/api/v4%s", c.cfg.BaseURL, path)
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("PRIVATE-TOKEN", c.cfg.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	nextPage := resp.Header.Get("X-Next-Page")
	return body, nextPage, nil
}

func (c *Client) paginatedGet(path string, params url.Values, collect func([]byte) (int, error)) error {
	p := make(url.Values)
	for k, v := range params {
		p[k] = v
	}
	p.Set("per_page", "100")

	for page := 1; ; page++ {
		p.Set("page", fmt.Sprintf("%d", page))

		body, nextPage, err := c.doGet(path, p)
		if err != nil {
			return err
		}

		n, err := collect(body)
		if err != nil {
			return err
		}

		if nextPage == "" || n == 0 {
			break
		}
	}
	return nil
}

func (c *Client) FetchMergeRequests(since time.Time) ([]MergeRequest, error) {
	var path string
	if c.cfg.ProjectID != "" {
		path = "/projects/" + url.PathEscape(c.cfg.ProjectID) + "/merge_requests"
	} else {
		path = "/groups/" + url.PathEscape(c.cfg.GroupID) + "/merge_requests"
	}

	params := url.Values{
		"state":         {"all"},
		"created_after": {since.Format(time.RFC3339)},
		"scope":         {"all"},
	}

	var mrs []MergeRequest
	err := c.paginatedGet(path, params, func(body []byte) (int, error) {
		var page []MergeRequest
		if err := json.Unmarshal(body, &page); err != nil {
			return 0, err
		}
		mrs = append(mrs, page...)
		return len(page), nil
	})
	return mrs, err
}

func (c *Client) FetchApprovalEvents(since time.Time) ([]Event, error) {
	var path string
	if c.cfg.ProjectID != "" {
		path = "/projects/" + url.PathEscape(c.cfg.ProjectID) + "/events"
	} else {
		path = "/groups/" + url.PathEscape(c.cfg.GroupID) + "/events"
	}

	params := url.Values{
		"action":      {"approved"},
		"target_type": {"MergeRequest"},
		"after":       {since.Format("2006-01-02")},
	}

	var events []Event
	err := c.paginatedGet(path, params, func(body []byte) (int, error) {
		var page []Event
		if err := json.Unmarshal(body, &page); err != nil {
			return 0, err
		}
		events = append(events, page...)
		return len(page), nil
	})
	return events, err
}

func (c *Client) FetchDashboardData() (*DashboardData, error) {
	since := time.Now().AddDate(0, 0, -c.cfg.Days)

	var (
		mrs       []MergeRequest
		events    []Event
		mrsErr    error
		eventsErr error
		wg        sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		mrs, mrsErr = c.FetchMergeRequests(since)
	}()
	go func() {
		defer wg.Done()
		events, eventsErr = c.FetchApprovalEvents(since)
	}()
	wg.Wait()

	if mrsErr != nil {
		return nil, fmt.Errorf("fetching merge requests: %w", mrsErr)
	}
	if eventsErr != nil {
		return nil, fmt.Errorf("fetching approval events: %w", eventsErr)
	}

	data := &DashboardData{
		ReviewsByUser: make(map[string]int),
		MRsByAuthor:   make(map[string]int),
		FetchedAt:     time.Now(),
		Since:         since,
		Days:          c.cfg.Days,
	}

	for _, mr := range mrs {
		data.TotalMRs++
		data.MRsByAuthor[mr.Author.Username]++
		switch mr.State {
		case "opened":
			data.OpenMRs++
		case "merged":
			data.MergedMRs++
		case "closed":
			data.ClosedMRs++
		}
	}

	for _, event := range events {
		if event.Author.Username != "" {
			data.ReviewsByUser[event.Author.Username]++
		}
	}

	return data, nil
}
