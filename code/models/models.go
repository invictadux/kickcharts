package models

import (
	"database/sql"
	"net/http"
	"time"
)

type Channel struct {
	ID             int    `json:"-"`
	Username       string `json:"username"`
	Slug           string `json:"slug"`
	Banner         string `json:"banner"`
	Picture        string `json:"picture"`
	IsBanned       bool   `json:"is_banned"`
	Language       string `json:"language"`
	Live           bool   `json:"live"`
	LiveViewers    int    `json:"live_viewers"`
	FollowersCount int    `json:"followers_count"`
	PeakViewers    int    `json:"peak_viewers"`
	Description    string `json:"description"`
	Discord        string `json:"discord"`
	Facebook       string `json:"facebook"`
	Instagram      string `json:"instagram"`
	Tiktok         string `json:"tiktok"`
	Twitter        string `json:"twitter"`
	Youtube        string `json:"youtube"`
}

func (c *Channel) Scan(rows *sql.Rows) error {
	err := rows.Scan(&c.ID, &c.Username, &c.Slug, &c.Banner,
		&c.Picture, &c.IsBanned, &c.Language, &c.Live, &c.LiveViewers,
		&c.FollowersCount, &c.PeakViewers, &c.Description, &c.Discord,
		&c.Facebook, &c.Instagram, &c.Tiktok, &c.Twitter, &c.Youtube)

	return err
}

func (c *Channel) ScanRow(row *sql.Row) error {
	err := row.Scan(&c.ID, &c.Username, &c.Slug, &c.Banner,
		&c.Picture, &c.IsBanned, &c.Language, &c.Live, &c.LiveViewers,
		&c.FollowersCount, &c.PeakViewers, &c.Description, &c.Discord,
		&c.Facebook, &c.Instagram, &c.Tiktok, &c.Twitter, &c.Youtube)

	return err
}

type Category struct {
	ID           int    `json:"-"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Banner       string `json:"banner"`
	LiveViewers  int    `json:"live_viewers"`
	LiveChannels int    `json:"live_channels"`
	PeakViewers  int    `json:"peak_viewers"`
	PeakChannels int    `json:"peak_channels"`
	Description  string `json:"description"`
}

func (c *Category) Scan(rows *sql.Rows) error {
	err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Banner,
		&c.LiveViewers, &c.LiveChannels, &c.PeakViewers,
		&c.PeakChannels, &c.Description)

	return err
}

func (c *Category) ScanRow(row *sql.Row) error {
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Banner,
		&c.LiveViewers, &c.LiveChannels, &c.PeakViewers,
		&c.PeakChannels, &c.Description)

	return err
}

type Clip struct {
	ID           string    `json:"-"`
	CategoryName string    `json:"category_name"`
	CategorySlug string    `json:"category_slug"`
	ChannelName  string    `json:"channel_name"`
	ChannelSlug  string    `json:"channel_slug"`
	IsMature     bool      `json:"is_mature"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Likes        int       `json:"likes"`
	LivestreamID string    `json:"livestream_id"`
	Thumbnail    string    `json:"thumbnail"`
	Views        int       `json:"views"`
	Duration     int       `json:"duration"`
	CreatedAt    time.Time `json:"created_at"`
}

func (c *Clip) Scan(rows *sql.Rows) error {
	err := rows.Scan(&c.ID, &c.CategoryName, &c.CategorySlug, &c.ChannelName,
		&c.ChannelSlug, &c.IsMature, &c.Title, &c.URL, &c.Likes, &c.LivestreamID,
		&c.Thumbnail, &c.Views, &c.Duration, &c.CreatedAt)

	return err
}

type Graph struct {
	Dates  []string
	Values []int
}

func (g *Graph) Scan(rows *sql.Rows) error {
	for rows.Next() {
		var date string
		var value int
		rows.Scan(&date, &value)
		t, _ := time.Parse("2006-01-02 15:04:05", date)
		g.Dates = append(g.Dates, t.Format("Jan 02, 2006 15:04"))
		g.Values = append(g.Values, value)
	}

	return nil
}

type Pagination struct {
	Page int
	Path string
	Sort string
}

func (p *Pagination) AddPath(r *http.Request) {
	params := r.URL.Query()
	params.Del("page")

	if len(params) == 0 {
		p.Path = r.URL.Path + "?"
	} else {
		p.Path = r.URL.Path + "?" + params.Encode() + "&"
	}

	if params.Has("sort") {
		p.Sort = params.Get("sort")
	} else {
		p.Sort = "lv"
	}
}
