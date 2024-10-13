package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"invictadux/code/models"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() *sql.DB {
	db, _ = sql.Open("sqlite3", "file:database.db?cache=shared")
	return db
}

func InsertChannel(c models.Channel) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO channels (username, slug, banner, picture, is_banned, language,
	live, live_viewers, followers_count, peak_viewers, description, discord, facebook, instagram, tiktok,
	twitter, youtube) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, c.Username, c.Slug, c.Banner, c.Picture, c.IsBanned,
		c.Language, c.Live, c.LiveViewers, c.FollowersCount, c.PeakViewers, c.Description, c.Discord, c.Facebook,
		c.Instagram, c.Tiktok, c.Twitter, c.Youtube)

	if err != nil {
		return err
	}

	return nil
}

func InsertCategory(c models.Category) error {
	_, err := db.Exec(`INSERT INTO categories (name, slug, banner, live_viewers, live_channels,
	peak_viewers, peak_channels, description) VALUES (?,?,?,?,?,?,?,?)`, c.Name, c.Slug, c.Banner,
		c.LiveViewers, c.LiveChannels, c.PeakViewers, c.PeakChannels, c.Description)

	if err != nil {
		return err
	}

	return nil
}

func InsertClip(c models.Clip) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO clips (id, category_name, category_slug, 
	channel_name, channel_slug, is_mature, title, url, likes, livestream_id, thumbnail, 
	views, duration, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, c.ID, c.CategoryName,
		c.CategorySlug, c.ChannelName, c.ChannelSlug, c.IsMature, c.Title, c.URL, c.Likes,
		c.LivestreamID, c.Thumbnail, c.Views, c.Duration, c.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func InsertOverallLiveChannelsChartPoint(n int) error {
	_, err := db.Exec(`INSERT INTO overall_live_channels_chart (ts, n) 
	VALUES (DATETIME('now'),?)`, n)

	if err != nil {
		return err
	}

	return nil
}

func InsertOverallViewersChartPoint(n int) error {
	_, err := db.Exec(`INSERT INTO overall_viewers_chart (ts, n) VALUES (DATETIME('now'),?)`, n)

	if err != nil {
		return err
	}

	return nil
}

func InsertChannelFollowersChartPoint(slug string, n int) error {
	_, err := db.Exec(`INSERT INTO channels_followers_chart (slug, ts, n) VALUES (?,DATETIME('now'),?)`, slug, n)

	if err != nil {
		return err
	}

	return nil
}

func InsertChannelViewersChartPoint(slug string, n int) error {
	_, err := db.Exec(`INSERT INTO channels_views_chart (slug, ts, n) VALUES (?,DATETIME('now'),?)`, slug, n)

	if err != nil {
		return err
	}

	return nil
}

func InsertCategoryLiveChannelsChartPoint(slug string, n int) error {
	_, err := db.Exec(`INSERT INTO categories_live_channels_chart (slug, ts, n) VALUES (?,DATETIME('now'),?)`, slug, n)

	if err != nil {
		return err
	}

	return nil
}

func InsertCategoryViewsChartPoint(slug string, n int) error {
	_, err := db.Exec(`INSERT INTO categories_views_chart (slug, ts, n) VALUES (?,DATETIME('now'),?)`, slug, n)

	if err != nil {
		return err
	}

	return nil
}

func UpdateChannel(id int, params map[string]interface{}) error {
	if len(params) == 0 {
		return fmt.Errorf("no parameters provided to update")
	}

	keys := []string{}
	values := []interface{}{}

	for key, val := range params {
		keys = append(keys, fmt.Sprintf("%s=?", key))
		values = append(values, val)
	}

	setClause := strings.Join(keys, ", ")

	sqlStmt := fmt.Sprintf("UPDATE channels SET %s WHERE id=?", setClause)

	values = append(values, id)

	_, err := db.Exec(sqlStmt, values...)

	if err != nil {
		return err
	}

	return nil
}

func UpdateCategory(slug string, params map[string]interface{}) error {
	if len(params) == 0 {
		return fmt.Errorf("no parameters provided to update")
	}

	keys := []string{}
	values := []interface{}{}

	for key, val := range params {
		keys = append(keys, fmt.Sprintf("%s=?", key))
		values = append(values, val)
	}

	setClause := strings.Join(keys, ", ")

	sqlStmt := fmt.Sprintf("UPDATE categories SET %s WHERE slug=?", setClause)

	values = append(values, slug)

	_, err := db.Exec(sqlStmt, values...)

	if err != nil {
		return err
	}

	return nil
}

func UpdateClip(id string, params map[string]interface{}) error {
	if len(params) == 0 {
		return fmt.Errorf("no parameters provided to update")
	}

	keys := []string{}
	values := []interface{}{}

	for key, val := range params {
		keys = append(keys, fmt.Sprintf("%s=?", key))
		values = append(values, val)
	}

	setClause := strings.Join(keys, ", ")

	sqlStmt := fmt.Sprintf("UPDATE clips SET %s WHERE id=?", setClause)

	values = append(values, id)

	_, err := db.Exec(sqlStmt, values...)

	if err != nil {
		return err
	}

	return nil
}

func GetChannel(slug string) (models.Channel, error) {
	row := db.QueryRow(`SELECT * FROM channels WHERE slug=?`, slug)

	var channel models.Channel
	err := channel.ScanRow(row)
	return channel, err
}

func GetChannels(offset, limit int) ([]models.Channel, error) {
	channels := []models.Channel{}
	rows, err := db.Query(`SELECT * FROM channels ORDER BY peak_viewers DESC limit ?,?`, offset, limit)

	if err != nil {
		return channels, err
	}

	for rows.Next() {
		var channel models.Channel
		channel.Scan(rows)
		channels = append(channels, channel)
	}

	return channels, nil
}

func GetAllChannels() []string {
	channels := []string{}
	rows, err := db.Query(`SELECT slug FROM channels`)

	if err != nil {
		return channels
	}

	for rows.Next() {
		var slug string
		rows.Scan(&slug)
		channels = append(channels, slug)
	}

	return channels
}

func GetCategoriesStats(offset, limit int, sortField string) ([]models.Category, int, int, int, int) {
	categories := []models.Category{}
	mostViews := 0
	mostChannels := 0
	peakViews := 0
	peakChannels := 0

	switch sortField {
	case "lv":
		sortField = "live_viewers"
	case "pv":
		sortField = "peak_viewers"
	case "lc":
		sortField = "live_channels"
	case "pc":
		sortField = "peak_channels"
	default:
		sortField = "live_viewers"
	}

	query := fmt.Sprintf(`SELECT * FROM categories ORDER BY %s DESC limit ?,?`, sortField)

	rows, err := db.Query(query, offset, limit)

	if err != nil {
		return categories, mostViews, mostChannels, peakViews, peakChannels
	}

	for rows.Next() {
		var c models.Category
		c.Scan(rows)
		categories = append(categories, c)

		if c.LiveViewers > mostViews {
			mostViews = c.LiveViewers
		}

		if c.LiveChannels > mostChannels {
			mostChannels = c.LiveChannels
		}

		if c.PeakViewers > peakViews {
			peakViews = c.PeakViewers
		}

		if c.PeakChannels > peakChannels {
			peakChannels = c.PeakChannels
		}
	}

	return categories, mostViews, mostChannels, peakViews, peakChannels
}

func GetChannelsStats(offset, limit int, sortField string) ([]models.Channel, int, int, int) {
	channels := []models.Channel{}
	mostViews := 0
	peakViews := 0
	mostFollowers := 0

	switch sortField {
	case "lv":
		sortField = "live_viewers"
	case "pv":
		sortField = "peak_viewers"
	case "f":
		sortField = "followers_count"
	default:
		sortField = "live_viewers"
	}

	query := fmt.Sprintf(`SELECT * FROM channels ORDER BY %s DESC limit ?,?`, sortField)
	rows, err := db.Query(query, offset, limit)

	if err != nil {
		return channels, mostViews, peakViews, mostFollowers
	}

	for rows.Next() {
		var c models.Channel
		c.Scan(rows)
		channels = append(channels, c)

		if c.LiveViewers > mostViews {
			mostViews = c.LiveViewers
		}

		if c.PeakViewers > peakViews {
			peakViews = c.PeakViewers
		}

		if c.FollowersCount > mostFollowers {
			mostFollowers = c.FollowersCount
		}
	}

	return channels, mostViews, peakViews, mostFollowers
}

func GetChannelStats(slug string) map[string]interface{} {
	row := db.QueryRow(`WITH last_30_days AS (
		SELECT slug, ts, n 
		FROM channels_views_chart
		WHERE slug = ? 
		AND ts >= datetime('now', '-30 days')
	),
	views_stats AS (
		-- Calculate peak viewers and average viewers
		SELECT 
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers,
			COUNT(DISTINCT ts) * 1.0 / 60 AS airtime_hours
		FROM last_30_days
	),
	followers_stats AS (
		SELECT 
			MAX(n) - MIN(n) AS followers_gain
		FROM channels_followers_chart
		WHERE slug = ? 
		AND ts >= datetime('now', '-30 days')
	)
	SELECT 
		v.peak_viewers,
		v.average_viewers,
		v.airtime_hours,
		f.followers_gain,
		(v.average_viewers * v.airtime_hours) AS hours_watched
	FROM views_stats v
	JOIN followers_stats f;`, slug, slug)

	var peakViewers int
	var averageViewers string
	var airtime int
	var followersGain string
	var hoursWatched int

	row.Scan(&peakViewers, &averageViewers, &airtime,
		&followersGain, &hoursWatched)

	data := map[string]interface{}{}
	data["Airtime"] = airtime
	data["PeakViewers"] = peakViewers
	data["AverageViewers"] = averageViewers
	data["FollowersGain"] = followersGain
	data["HoursWatched"] = hoursWatched
	return data
}

func GetClips(params url.Values, offset, limit int) ([]models.Clip, error) {
	clips := []models.Clip{}

	query := `SELECT * FROM clips`
	var args []interface{}

	conditions := []string{}

	if params.Has("category") {
		conditions = append(conditions, "category_slug = ?")
		args = append(args, params.Get("category"))
	}

	if params.Has("channel") {
		conditions = append(conditions, "channel_slug = ?")
		args = append(args, params.Get("channel"))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	switch params.Get("sort") {
	case "v":
		query += " ORDER BY views DESC"
	case "dd":
		query += " ORDER BY duration DESC"
	case "da":
		query += " ORDER BY duration ASC"
	case "cd":
		query += " ORDER BY created_at DESC"
	case "ca":
		query += " ORDER BY created_at ASC"
	default:
		query += " ORDER BY views DESC"
	}

	query += " LIMIT ?, ?"
	args = append(args, offset, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return clips, err
	}
	defer rows.Close()

	for rows.Next() {
		var clip models.Clip
		if err := clip.Scan(rows); err != nil {
			return clips, err
		}
		clips = append(clips, clip)
	}

	return clips, nil
}

func GetCategory(slug string) (models.Category, error) {
	row := db.QueryRow(`SELECT * FROM categories WHERE slug=?`, slug)

	var category models.Category
	err := category.ScanRow(row)
	return category, err
}

func GetCategories(offset, limit int) ([]models.Category, error) {
	categorries := []models.Category{}
	rows, err := db.Query(`SELECT * FROM categories ORDER BY peak_viewers DESC limit ?,?`, offset, limit)

	if err != nil {
		return categorries, err
	}

	for rows.Next() {
		var category models.Category
		category.Scan(rows)
		categorries = append(categorries, category)
	}

	return categorries, nil
}

func GetAllCategories() []string {
	categories := []string{}
	rows, err := db.Query(`SELECT slug FROM categories`)

	if err != nil {
		return categories
	}

	for rows.Next() {
		var slug string
		rows.Scan(&slug)
		categories = append(categories, slug)
	}

	return categories
}

func GetOverallChannelsGraph(t1, t2 string) models.Graph {
	rows, err := db.Query(`	WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN overall_live_channels_chart cl ON cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC`, t1, t2)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetOverallViewsGraph(t1, t2 string) models.Graph {
	rows, err := db.Query(`	WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN overall_viewers_chart cl ON cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC`, t1, t2)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetChannelFollowersGraph(channel, t1, t2 string) models.Graph {
	rows, err := db.Query(`WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 day'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN channels_followers_chart cl ON cl.slug=? AND cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 day') GROUP BY c.dt ORDER BY c.dt ASC;`, t1, t2, channel)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetChannelViewsGraph(channel, t1, t2 string) models.Graph {
	rows, err := db.Query(`WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN channels_views_chart cl ON cl.slug=? AND cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC`, t1, t2, channel)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetCategoryChannelsGraph(channel, t1, t2 string) models.Graph {
	rows, err := db.Query(`WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN categories_live_channels_chart cl ON cl.slug=? AND cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC;`, t1, t2, channel)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetCategoryViewsGraph(channel, t1, t2 string) models.Graph {
	rows, err := db.Query(`WITH RECURSIVE cte AS (SELECT DATETIME(?) as dt, DATETIME(?)
	last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT cl.ts,
	n FROM cte c LEFT JOIN categories_views_chart cl ON cl.slug=? AND cl.ts >= c.dt AND
	cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC;`, t1, t2, channel)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	graph := models.Graph{}
	graph.Scan(rows)

	return graph
}

func GetCategoryGrowthData(category string) map[string]interface{} {
	rows, err := db.Query(`WITH views_data AS (
		SELECT 
			slug,
			DATE(ts) AS date,
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers
		FROM 
			categories_views_chart
		WHERE 
			slug = ?
		GROUP BY 
			slug, date
	),
	channels_data AS (
		SELECT 
			slug,
			DATE(ts) AS date,
			MAX(n) AS peak_channels,
			AVG(n) AS average_channels
		FROM 
			categories_live_channels_chart
		WHERE 
			slug = ?
		GROUP BY 
			slug, date
	)
	SELECT 
		v.date,
		v.peak_viewers,
		v.average_viewers,
		c.peak_channels,
		c.average_channels
	FROM 
		views_data v
	JOIN 
		channels_data c
	ON 
		v.slug = c.slug
	AND 
		v.date = c.date
	ORDER BY 
		v.date`, category, category)

	if err != nil {
		return map[string]interface{}{}
	}

	var dates []string
	var peakViewers []int
	var avgViewers []int
	var peakChannels []int
	var avgChannels []int

	for rows.Next() {
		var d string
		var pv int
		var av float64
		var pc int
		var ac float64

		rows.Scan(&d, &pv, &av, &pc, &ac)
		t, _ := time.Parse("2006-01-02", d)
		dates = append(dates, t.Format("Jan 02, 2006")+" UTC")
		peakViewers = append(peakViewers, pv)
		avgViewers = append(avgViewers, int(av))
		peakChannels = append(peakChannels, pc)
		avgChannels = append(avgChannels, int(ac))
	}

	data := map[string]interface{}{}
	data["Dates"] = dates
	data["PeakViewers"] = peakViewers
	data["AverageViewers"] = avgViewers
	data["PeakChannels"] = peakChannels
	data["AverageChannels"] = avgChannels
	return data
}

func GetChannelGrowthData(channel string) map[string]interface{} {
	rows, err := db.Query(`
	WITH views_data AS (
		SELECT 
			slug,
			DATE(ts) AS date,
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers
		FROM 
			channels_views_chart
		WHERE 
			slug = ?
		GROUP BY 
			slug, date
	)
	SELECT 
		date,
		peak_viewers,
		average_viewers
	FROM 
		views_data
	ORDER BY 
		date`, channel)

	if err != nil {
		return map[string]interface{}{}
	}

	var dates []string
	var peakViewers []int
	var avgViewers []int

	for rows.Next() {
		var d string
		var pv int
		var av float64

		rows.Scan(&d, &pv, &av)
		t, _ := time.Parse("2006-01-02", d)
		dates = append(dates, t.Format("Jan 02, 2006")+" UTC")
		peakViewers = append(peakViewers, pv)
		avgViewers = append(avgViewers, int(av))
	}

	data := map[string]interface{}{}
	data["Dates"] = dates
	data["PeakViewers"] = peakViewers
	data["AverageViewers"] = avgViewers
	return data
}

func GetChannelLast30DGraphs(channel string) map[string]interface{} {
	rows, err := db.Query(`WITH RECURSIVE date_series AS (
		SELECT date('now', '-29 days') AS date
		UNION ALL
		SELECT date(date, '+1 day')
		FROM date_series
		WHERE date < date('now')
	),
	views_data AS (
		SELECT 
			DATE(ts) AS date,
			SUM(n) AS hours_watched,
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers
		FROM 
			channels_views_chart
		WHERE 
			slug = ?
		GROUP BY 
			DATE(ts)
	),
	followers_data AS (
		SELECT 
			DATE(ts) AS date,
			SUM(n) AS followers_gain
		FROM 
			channels_followers_chart
		WHERE 
			slug = ?
		GROUP BY 
			DATE(ts)
	),
	combined_data AS (
		SELECT
			ds.date,
			IFNULL(vd.hours_watched, 0) AS hours_watched,
			IFNULL(vd.peak_viewers, 0) AS peak_viewers,
			IFNULL(vd.average_viewers, 0) AS average_viewers,
			IFNULL(fd.followers_gain, 0) AS followers_gain
		FROM 
			date_series ds
		LEFT JOIN 
			views_data vd ON ds.date = vd.date
		LEFT JOIN 
			followers_data fd ON ds.date = fd.date
	)
	SELECT
		date,
		hours_watched,
		peak_viewers,
		average_viewers,
		followers_gain,
		CASE 
			WHEN hours_watched > 0 THEN 1
			ELSE 0
		END AS airtime
	FROM 
		combined_data
	ORDER BY 
		date`, channel, channel)

	if err != nil {
		return map[string]interface{}{}
	}

	var dates []string
	var hoursWatched []int
	var peakViewers []int
	var avgViewers []int
	var followersGain []int
	var airTime []int

	for rows.Next() {
		var d string
		var hw int
		var pv int
		var av float64
		var fg int
		var at int

		rows.Scan(&d, &hw, &pv, &av, &fg, &at)
		t, _ := time.Parse("2006-01-02", d)
		dates = append(dates, t.Format("02 Jan"))
		hoursWatched = append(hoursWatched, hw)
		peakViewers = append(peakViewers, pv)
		avgViewers = append(avgViewers, int(av))
		followersGain = append(followersGain, fg)
		airTime = append(airTime, at)
	}

	data := map[string]interface{}{}
	data["Dates"] = dates
	data["HoursWatched"] = hoursWatched
	data["PeakViewers"] = peakViewers
	data["AverageViewers"] = avgViewers
	data["FollowersGain"] = followersGain
	data["AirTime"] = airTime
	return data
}

func GetChannelsSlug(offset, limit int) *chan string {
	ch := make(chan string, limit)
	rows, err := db.Query(`SELECT slug FROM channels LIMIT ?,?`, offset, limit)

	if err != nil {
		close(ch)
		return &ch
	}

	for rows.Next() {
		var slug string
		rows.Scan(&slug)
		ch <- slug
	}

	close(ch)
	return &ch
}

func GetOverallViewsStats() (int, string, int, string, int) {
	row := db.QueryRow(`WITH recent_data AS (
	SELECT n AS last_30_days_peak, ts AS peak_date
	FROM overall_viewers_chart
	WHERE DATE(ts) >= DATE('now', '-30 days')
	ORDER BY n DESC
	LIMIT 1
	),
	all_time_data AS (
	SELECT MAX(n) AS all_time_views, ts AS peak_date
	FROM overall_viewers_chart
	),
	last_7_days_avg AS (
	SELECT AVG(n) AS last_7_days_average
	FROM overall_viewers_chart
	WHERE DATE(ts) >= DATE('now', '-7 days')
	)
	SELECT
	recent_data.last_30_days_peak,
	recent_data.peak_date,
	all_time_data.all_time_views,
	all_time_data.peak_date,
	last_7_days_avg.last_7_days_average
	FROM
	recent_data,
	all_time_data,
	last_7_days_avg;`)

	var last30DaysPeak int
	var last30DaysPeakDate string
	var allTimePeak int
	var allTimePeakDate string
	var last7DaysAvg float64
	row.Scan(&last30DaysPeak, &last30DaysPeakDate, &allTimePeak, &allTimePeakDate, &last7DaysAvg)

	return last30DaysPeak, last30DaysPeakDate, allTimePeak, allTimePeakDate, int(last7DaysAvg)
}

func GetOverallChannelsStats() (int, string, int, string, int) {
	row := db.QueryRow(`WITH recent_data AS (
	SELECT n AS last_30_days_peak, ts AS peak_date
	FROM overall_live_channels_chart
	WHERE DATE(ts) >= DATE('now', '-30 days')
	ORDER BY n DESC
	LIMIT 1
	),
	all_time_data AS (
	SELECT MAX(n) AS all_time_views, ts AS peak_date
	FROM overall_live_channels_chart
	),
	last_7_days_avg AS (
	SELECT AVG(n) AS last_7_days_average
	FROM overall_live_channels_chart
	WHERE DATE(ts) >= DATE('now', '-7 days')
	)
	SELECT
	recent_data.last_30_days_peak,
	recent_data.peak_date,
	all_time_data.all_time_views,
	all_time_data.peak_date,
	last_7_days_avg.last_7_days_average
	FROM
	recent_data,
	all_time_data,
	last_7_days_avg;`)

	var last30DaysPeak int
	var last30DaysPeakDate string
	var allTimePeak int
	var allTimePeakDate string
	var last7DaysAvg float64
	row.Scan(&last30DaysPeak, &last30DaysPeakDate, &allTimePeak, &allTimePeakDate, &last7DaysAvg)

	return last30DaysPeak, last30DaysPeakDate, allTimePeak, allTimePeakDate, int(last7DaysAvg)
}

func GetCategoryViewsStats(slug string) map[string]interface{} {
	row := db.QueryRow(`SELECT 
		(SELECT ts 
		FROM categories_views_chart 
		WHERE slug = ? 
		AND ts >= date('now', '-7 days') 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_7_days,
		
		MAX(CASE WHEN ts >= date('now', '-7 days') THEN n END) AS peak_views_7_days,

		(SELECT ts 
		FROM categories_views_chart 
		WHERE slug = ? 
		AND ts >= date('now', '-30 days') 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_30_days,
		
		MAX(CASE WHEN ts >= date('now', '-30 days') THEN n END) AS peak_views_30_days,

		(SELECT ts 
		FROM categories_views_chart 
		WHERE slug = ? 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_all_time,
		
		MAX(n) AS peak_views_all_time
	FROM categories_views_chart
	WHERE slug = ?`, slug, slug, slug, slug)

	var last7DaysPeak int
	var last7DaysPeakDate string
	var last30DaysPeak int
	var last30DaysPeakDate string
	var allTimePeak int
	var allTimePeakDate string

	row.Scan(&last7DaysPeakDate, &last7DaysPeak, &last30DaysPeakDate,
		&last30DaysPeak, &allTimePeakDate, &allTimePeak)

	data := map[string]interface{}{}
	data["Last7DaysPeak"] = last7DaysPeak
	data["Last7DaysPeakDate"] = last7DaysPeakDate
	data["Last30DaysPeak"] = last30DaysPeak
	data["Last30DaysPeakDate"] = last30DaysPeakDate
	data["AllTimePeak"] = allTimePeak
	data["AllTimePeakDate"] = allTimePeakDate
	return data
}

func GetCategoryChannelStats(slug string) map[string]interface{} {
	row := db.QueryRow(`SELECT 
		(SELECT ts 
		FROM categories_live_channels_chart 
		WHERE slug = ? 
		AND ts >= date('now', '-7 days') 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_7_days,
		
		MAX(CASE WHEN ts >= date('now', '-7 days') THEN n END) AS peak_views_7_days,

		(SELECT ts 
		FROM categories_live_channels_chart 
		WHERE slug = ? 
		AND ts >= date('now', '-30 days') 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_30_days,
		
		MAX(CASE WHEN ts >= date('now', '-30 days') THEN n END) AS peak_views_30_days,

		(SELECT ts 
		FROM categories_live_channels_chart 
		WHERE slug = ? 
		ORDER BY n DESC, ts DESC 
		LIMIT 1) AS ts_peak_all_time,
		
		MAX(n) AS peak_views_all_time
	FROM categories_live_channels_chart
	WHERE slug = ?`, slug, slug, slug, slug)

	var last7DaysPeak int
	var last7DaysPeakDate string
	var last30DaysPeak int
	var last30DaysPeakDate string
	var allTimePeak int
	var allTimePeakDate string

	row.Scan(&last7DaysPeakDate, &last7DaysPeak, &last30DaysPeakDate,
		&last30DaysPeak, &allTimePeakDate, &allTimePeak)

	data := map[string]interface{}{}
	data["Last7DaysPeak"] = last7DaysPeak
	data["Last7DaysPeakDate"] = last7DaysPeakDate
	data["Last30DaysPeak"] = last30DaysPeak
	data["Last30DaysPeakDate"] = last30DaysPeakDate
	data["AllTimePeak"] = allTimePeak
	data["AllTimePeakDate"] = allTimePeakDate
	return data
}

func GetViewersChartStats(t string) map[string]interface{} {
	v := ""

	switch t {
	case "w":
		v = "strftime('%Y-%W', ts) AS dt"
	case "m":
		v = "strftime('%Y-%m', ts) AS dt"
	case "q":
		v = "strftime('%Y-Q%q', ts) AS dt"
	}

	rows, _ := db.Query(fmt.Sprintf(`WITH last AS (
		SELECT 
			%s,
			ts,
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers
		FROM 
			overall_viewers_chart
		GROUP BY 
			dt
		ORDER BY 
			dt DESC
		LIMIT 24
	)
	SELECT * FROM last ORDER BY dt ASC`, v))

	dates := []string{}
	peakValues := []int{}
	avgValues := []int{}

	for rows.Next() {
		var (
			d  string
			ts string
			pv int
			av float64
		)
		rows.Scan(&d, &ts, &pv, &av)
		t, _ := time.Parse("2006-01-02 15:04:05", ts)
		dates = append(dates, t.Format("Mon, Jan 02, 2006")+" UTC")
		peakValues = append(peakValues, pv)
		avgValues = append(avgValues, int(av))
	}

	data := map[string]interface{}{
		"dates": dates,
		"peak":  peakValues,
		"avg":   avgValues,
	}

	return data
}

func GetChannelsChartStats(t string) map[string]interface{} {
	v := ""

	switch t {
	case "w":
		v = "strftime('%Y-%W', ts) AS dt"
	case "m":
		v = "strftime('%Y-%m', ts) AS dt"
	case "q":
		v = "strftime('%Y-Q%q', ts) AS dt"
	}

	rows, _ := db.Query(fmt.Sprintf(`WITH last AS (
		SELECT 
			%s,
			ts,
			MAX(n) AS peak_viewers,
			AVG(n) AS average_viewers
		FROM 
			overall_live_channels_chart
		GROUP BY 
			dt
		ORDER BY 
			dt DESC
		LIMIT 24
	)
	SELECT * FROM last ORDER BY dt ASC`, v))

	dates := []string{}
	peakValues := []int{}
	avgValues := []int{}

	for rows.Next() {
		var (
			d  string
			ts string
			pv int
			av float64
		)
		rows.Scan(&d, &ts, &pv, &av)
		t, _ := time.Parse("2006-01-02 15:04:05", ts)
		dates = append(dates, t.Format("Mon, Jan 02, 2006")+" UTC")
		peakValues = append(peakValues, pv)
		avgValues = append(avgValues, int(av))
	}

	data := map[string]interface{}{
		"dates": dates,
		"peak":  peakValues,
		"avg":   avgValues,
	}

	return data
}

func SetChannelsLiveStatus(liveChannels *[]string) error {
	placeholders := make([]string, len(*liveChannels))
	args := make([]interface{}, len(*liveChannels))
	for i, channel := range *liveChannels {
		placeholders[i] = "?"
		args[i] = channel
	}

	query := fmt.Sprintf(`UPDATE channels SET live_viewers = 0, live = 0 
        WHERE slug NOT IN (%s)`, strings.Join(placeholders, ", "))

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

/*
func GetMainSuggestions(query string) []models.Sugg {
	suggestions := []models.Sugg{}

	rows, err := db.Query(`SELECT q, t, bt, p FROM suggestions_fts WHERE q MATCH ? LIMIT 10`, query+"*")

	if err != nil {
		return suggestions
	}

	defer rows.Close()

	for rows.Next() {
		var sugg models.Sugg
		sugg.Scan(rows)
		suggestions = append(suggestions, sugg)
	}
	return suggestions
}*/

/* CREATE TABLE clips (rowid integer primary key, id text not null unique,
broadcaster_id text, broadcaster_name text, creator_id text, creator_name text, video_id text,
game_id text, language text, title text, view_count int, duration REAL, created_at text, thumbnail_url text); */

/* CREATE VIRTUAL TABLE IF NOT EXISTS clips_fts USING FTS5(title, id); */

/* CREATE TABLE subscriptions (id TEXT, user_id TEXT, email TEXT,
item TEXT, amount TEXT, currency_code TEXT, created_at TEXT, ended_at TEXT); */

/* CREATE TABLE comments (id INTEGER PRIMARY KEY AUTOINCREMENT, clip_id TEXT, user_id TEXT,
comment TEXT, created_at TEXT); */

/* CREATE TABLE streamers (id TEXT primary key unique,
login TEXT, display_name TEXT, type TEXT, broadcaster_type TEXT, description TEXT, profile_image_url TEXT,
offline_image_url TEXT, view_count int, email TEXT, created_at TEXT, access_token TEXT, refresh_token TEXT, session TEXT); */

/* CREATE TABLE games (id text primary key unique, name text, box_art_url text); */

/*
WITH RECURSIVE cte AS (SELECT DATETIME('2024-04-15') as dt, DATETIME('now')
last_dt UNION ALL SELECT DATETIME(dt, '+1 hour'), last_dt FROM cte WHERE dt < last_dt) SELECT c.dt,
n FROM cte c LEFT JOIN categories_views_chart cl ON cl.slug='valorant' AND cl.ts >= c.dt AND
cl.ts < DATETIME(c.dt, '+1 hour') GROUP BY c.dt ORDER BY c.dt ASC;
*/

/*
INSERT INTO global_stats (clips, clips_30d_change, views, views_30d_change, streamers, streamers_30d_change, games, games_30d_change)
WITH t1 AS (SELECT COUNT(*) AS c, SUM(view_count) AS v, COUNT(DISTINCT broadcaster_name) AS b, COUNT(DISTINCT game_id) AS g FROM clips WHERE
created_at < strftime('%Y-%m-%dT%H:%M:%SZ', DATETIME('now')) AND created_at > strftime('%Y-%m-%dT%H:%M:%SZ', DATETIME('now', '-2 day'))),
t2 AS (SELECT COUNT(*) AS c, SUM(view_count) AS v, COUNT(DISTINCT broadcaster_name) AS b, COUNT(DISTINCT game_id) AS g FROM clips WHERE
created_at < strftime('%Y-%m-%dT%H:%M:%SZ', DATETIME('now', '-2 day')) AND created_at > strftime('%Y-%m-%dT%H:%M:%SZ', DATETIME('now', '-4 day'))),
t3 AS (SELECT COUNT(*) AS c, SUM(view_count) AS v, COUNT(DISTINCT broadcaster_name) AS b, COUNT(DISTINCT game_id) AS g FROM clips)
SELECT t3.c, CAST((t1.c - t2.c) * 100 AS REAL) / t2.c, t3.v, CAST((t1.v - t2.v) * 100 AS REAL) / t2.v, t3.b, CAST((t1.b - t2.b) * 100 AS REAL) / t2.b,
t3.g, CAST((t1.g - t2.g) * 100 AS REAL) / t2.g FROM t1 JOIN t2 JOIN t3;
*/

/*
WITH last_24_weeks AS (
    SELECT
        strftime('%Y-%W', ts) AS week,
		ts,    -- Extract year-week from timestamp
        MAX(n) AS peak_viewers,            -- Maximum viewers for the week
        AVG(n) AS average_viewers          -- Average viewers for the week
    FROM
        overall_viewers_chart
    GROUP BY
        week                              -- Group by week
    ORDER BY
        week DESC                        -- Order by week descending
    LIMIT 24                             -- Limit to the last 24 weeks
)
SELECT *
FROM last_24_weeks
ORDER BY week ASC;
*/
