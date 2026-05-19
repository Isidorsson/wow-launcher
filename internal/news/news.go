// Package news fetches a server's news/changelog feed.
//
// The feed format is a JSON array of NewsItem objects. Server operators
// host it at the news_feed_url in config.toml. Body is plain markdown.
// We do NOT render markdown server-side — the frontend handles it so the
// launcher can apply its own style.
package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	fetchTimeout = 10 * time.Second
	maxBytes     = 1 << 20 // 1 MB cap
)

type Item struct {
	Title    string `json:"title"`
	Date     string `json:"date"`             // ISO 8601 (YYYY-MM-DD); free-form tolerated
	Body     string `json:"body"`             // markdown
	URL      string `json:"url,omitempty"`    // optional "read more" link
	Category string `json:"category,omitempty"` // e.g. "patch", "event", "news"
}

func Fetch(ctx context.Context, url string) ([]Item, error) {
	if url == "" {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build news request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch news: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news http %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	if err != nil {
		return nil, fmt.Errorf("read news: %w", err)
	}
	var items []Item
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse news: %w", err)
	}
	return items, nil
}
