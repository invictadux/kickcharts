package models

import (
	"database/sql"
	"time"
)

type Channel struct {
	ID             int    `json:"-"`
	Username       string `json:"username"`
	Slug           string `json:"slug"`
	Banner         string `json:"banner"`
	Picture        string `json:"picture"`
	Language       string `json:"language"`
	IsBanned       bool   `json:"is_banned"`
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
		&c.Picture, &c.IsBanned, &c.Language, &c.FollowersCount, &c.PeakViewers,
		&c.Description, &c.Discord, &c.Facebook, &c.Instagram,
		&c.Tiktok, &c.Twitter, &c.Youtube)

	return err
}

func (c *Channel) ScanRow(row *sql.Row) error {
	err := row.Scan(&c.ID, &c.Username, &c.Slug, &c.Banner,
		&c.Picture, &c.IsBanned, &c.Language, &c.FollowersCount, &c.PeakViewers,
		&c.Description, &c.Discord, &c.Facebook, &c.Instagram,
		&c.Tiktok, &c.Twitter, &c.Youtube)

	return err
}

type Category struct {
	ID           int    `json:"-"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Banner       string `json:"banner"`
	PeakChannels int    `json:"peak_channels"`
	PeakViewers  int    `json:"peak_viewers"`
	Description  string `json:"description"`
}

func (c *Category) Scan(rows *sql.Rows) error {
	err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Banner,
		&c.PeakChannels, &c.PeakViewers, &c.Description)

	return err
}

func (c *Category) ScanRow(row *sql.Row) error {
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Banner,
		&c.PeakChannels, &c.PeakViewers, &c.Description)

	return err
}

type Clip struct {
	ID           string    `json:"-"`
	Category     string    `json:"category"`
	Channel      string    `json:"channel"`
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
	err := rows.Scan(&c.ID, &c.Category, &c.Channel, &c.IsMature,
		&c.Title, &c.URL, &c.Likes, &c.LivestreamID, &c.Thumbnail,
		&c.Views, &c.Duration, &c.CreatedAt)

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
		g.Dates = append(g.Dates, date)
		g.Values = append(g.Values, value)
	}

	return nil
}

func (g *Graph) ToTable() []Table {
	table := []Table{}
	v := g.Values[len(g.Values)-1]

	for i := len(g.Dates) - 2; i > 0; i-- {
		increment := v - g.Values[i]
		v = g.Values[i]

		t := Table{}
		t.Date = g.Dates[i+1]
		t.V1 = increment
		t.V2 = g.Values[i+1]
		table = append(table, t)

		if len(g.Dates)-i > 14 {
			break
		}
	}

	return table
}

type Table struct {
	Date string
	V1   int
	V2   int
}
