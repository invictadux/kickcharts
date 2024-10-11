package main

import (
	"context"
	"encoding/json"
	"html/template"
	"invictadux/code/db"
	"invictadux/code/funcmaps"
	"invictadux/code/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var indexTemplate *template.Template
var channelTemplate *template.Template
var channelsTemplate *template.Template
var categoryTemplate *template.Template
var categoriesTemplate *template.Template
var clipsTemplate *template.Template
var notFoundTemplate *template.Template

func MainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "invictadux")
		next.ServeHTTP(w, r)
	})
}

func APIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "invictadux")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := map[string]interface{}{}
		page["User"] = "Dan"
		rcopy := r.WithContext(context.WithValue(r.Context(), "result", page))
		next.ServeHTTP(w, rcopy)
	})
}

func GetCTX(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	result, ok := ctx.Value("result").(map[string]interface{})

	if ok {
		return result
	}

	return nil
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	notFoundTemplate.Execute(w, nil)
}

func NotFoundAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(`{"status":"error","message":"path unknown"}`))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if strings.HasPrefix(r.URL.Path, "/api/v1/") {
			NotFoundAPI(w, r)
		} else {
			NotFoundPage(w, r)
		}

		return
	}

	now := time.Now().UTC()
	t1 := now.Add(-time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	t2 := now.Add(-time.Hour).Format("2006-01-02 15:04:05")

	viewsChart := db.GetOverallViewsGraph(t1, t2)
	channelsChart := db.GetOverallChannelsGraph(t1, t2)
	categories, mostViews, mostChannels, _, _ := db.GetCategoriesStats(0, 10, "lv")
	channels, mostChannelViews, _, _ := db.GetChannelsStats(0, 10, "lv")
	clips, _ := db.GetClips(url.Values{}, 0, 10)
	peak30DViews, peak30DViewsDate, allTimePeakViews, allTimePeakViewsDate, avg7DViews := db.GetOverallViewsStats()
	peak30DChannels, peak30DChannelsDate, allTimePeakChannels, allTimePeakChannelsDate, avg7DChannels := db.GetOverallChannelsStats()

	page := map[string]interface{}{}
	page["ChannelsChart"] = channelsChart
	page["ViewsChart"] = viewsChart
	page["Categories"] = categories
	page["Channels"] = channels
	page["Clips"] = clips

	page["LiveViews"] = viewsChart.Values[len(viewsChart.Values)-1]
	page["Peak30DViews"] = peak30DViews
	page["Peak30DViewsDate"] = peak30DViewsDate
	page["AllTimePeakViews"] = allTimePeakViews
	page["AllTimePeakViewsDate"] = allTimePeakViewsDate
	page["Avg7DViews"] = avg7DViews

	page["LiveChannels"] = channelsChart.Values[len(channelsChart.Values)-1]
	page["Peak30DChannels"] = peak30DChannels
	page["Peak30DChannelsDate"] = peak30DChannelsDate
	page["AllTimePeakChannels"] = allTimePeakChannels
	page["AllTimePeakChannelsDate"] = allTimePeakChannelsDate
	page["Avg7DChannels"] = avg7DChannels

	page["MostCategoryViews"] = mostViews
	page["MostCategoryChannels"] = mostChannels
	page["MostChannelViews"] = mostChannelViews

	indexTemplate.Execute(w, page)
}

func ChannelPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("channel")

	now := time.Now().UTC()
	t1 := now.Add(-time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	t2 := now.Add(-time.Hour).Format("2006-01-02 15:04:05")
	t1Followers := now.Add(-time.Hour * 24 * 30).Format("2006-01-02 15:04:05")
	t2Followers := now.Add(-time.Hour * 24).Format("2006-01-02 15:04:05")

	channel, _ := db.GetChannel(slug)

	clips, _ := db.GetClips(url.Values{"channel": []string{slug}}, 0, 10)

	page := map[string]interface{}{}
	page["Channel"] = channel
	page["Stats"] = db.GetChannelStats(slug)
	followersGraph := db.GetChannelFollowersGraph(slug, t1Followers, t2Followers)
	page["FollowersChart"] = followersGraph
	page["ViewsChart"] = db.GetChannelViewsGraph(slug, t1, t2)
	page["Clips"] = clips
	page["Growth"] = db.GetChannelGrowthData(slug)
	page["Charts"] = db.GetChannelLast30DGraphs(slug)

	channelTemplate.Execute(w, page)
}

func ChannelsPage(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pagination := models.Pagination{}
	pagination.AddPath(r)
	pagination.Page, _ = strconv.Atoi(params.Get("page"))

	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	limit := 50
	offset := (pagination.Page - 1) * limit

	channels, mv, pv, mf := db.GetChannelsStats(offset, limit, params.Get("sort"))

	page := map[string]interface{}{}
	page["Channels"] = channels
	page["Offset"] = offset + 1
	page["MostChannelViewers"] = mv
	page["PeakChannelViewers"] = pv
	page["MostChannelFollowers"] = mf
	page["Pagination"] = pagination
	channelsTemplate.Execute(w, page)
}

func CategoryPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("category")

	now := time.Now().UTC()
	t1 := now.Add(-time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	t2 := now.Add(-time.Hour).Format("2006-01-02 15:04:05")

	clips, _ := db.GetClips(url.Values{"category": []string{slug}}, 0, 10)
	mostClipsViews := 0

	if len(clips) > 0 {
		mostClipsViews = clips[0].Views
	}

	category, _ := db.GetCategory(slug)

	page := map[string]interface{}{}
	page["Slug"] = slug
	page["Category"] = category
	page["Views"] = db.GetCategoryViewsStats(slug)
	page["Channels"] = db.GetCategoryChannelStats(slug)
	page["ChannelsChart"] = db.GetCategoryChannelsGraph(slug, t1, t2)
	page["ViewsChart"] = db.GetCategoryViewsGraph(slug, t1, t2)
	page["Growth"] = db.GetCategoryGrowthData(slug)
	page["Clips"] = clips
	page["MostClipViews"] = mostClipsViews

	categoryTemplate.Execute(w, page)
}

func CategoriesPage(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pagination := models.Pagination{}
	pagination.AddPath(r)
	pagination.Page, _ = strconv.Atoi(params.Get("page"))

	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	limit := 50
	offset := (pagination.Page - 1) * limit

	categories, mv, mc, pv, pc := db.GetCategoriesStats(offset, limit, params.Get("sort"))

	page := map[string]interface{}{}
	page["Categories"] = categories
	page["Offset"] = offset + 1
	page["MostCategoryViewers"] = mv
	page["MostCategoryChannels"] = mc
	page["PeakCategoryViewers"] = pv
	page["PeakCategoryChannels"] = pc
	page["Pagination"] = pagination
	categoriesTemplate.Execute(w, page)
}

func ClipsPage(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pagination := models.Pagination{}
	pagination.AddPath(r)
	pagination.Page, _ = strconv.Atoi(params.Get("page"))

	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	limit := 50
	offset := (pagination.Page - 1) * limit

	clips, _ := db.GetClips(params, offset, limit)

	page := map[string]interface{}{}
	page["Clips"] = clips
	page["Offset"] = offset + 1
	page["Pagination"] = pagination
	clipsTemplate.Execute(w, page)
}

//------------------ API ------------------

func IndexAPI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}

func ClientError(w http.ResponseWriter, r *http.Request, stausCode int) {
	w.WriteHeader(stausCode)
	w.Write([]byte(`{"status":"error","message":"You send some data incorrectly!"}`))
}

func ChartStatsAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if !params.Has("s") {
		w.Write([]byte(`{"status":"error","message":"s parameter required"}`))
		return
	} else if !params.Has("t") {
		w.Write([]byte(`{"status":"error","message":"t parameter required"}`))
		return
	}

	s := params.Get("s")
	t := params.Get("t")

	switch t {
	case "w", "m", "q":
	default:
		w.Write([]byte(`{"status":"error","message":"t parameter not valid"}`))
		return
	}

	var data any

	switch s {
	case "v":
		data = db.GetViewersChartStats(t)
	case "c":
		data = db.GetChannelsChartStats(t)
	default:
		w.Write([]byte(`{"status":"error","message":"s parameter required"}`))
		return
	}

	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	db.Init()
	//go scraper.Run()

	indexTemplate = funcmaps.NewTemplate("templates/index.html", "templates/templates.html")
	channelTemplate = funcmaps.NewTemplate("templates/channel.html", "templates/templates.html")
	channelsTemplate = funcmaps.NewTemplate("templates/channels.html", "templates/templates.html")
	categoryTemplate = funcmaps.NewTemplate("templates/category.html", "templates/templates.html")
	categoriesTemplate = funcmaps.NewTemplate("templates/categories.html", "templates/templates.html")
	clipsTemplate = funcmaps.NewTemplate("templates/clips.html", "templates/templates.html")
	notFoundTemplate = funcmaps.NewTemplate("templates/notfound.html", "templates/templates.html")

	//------------------ Static ------------------
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	//------------------ Main ------------------
	mux.Handle("GET /", MainMiddleware(http.HandlerFunc(IndexPage)))
	mux.Handle("GET /channel/{channel}", MainMiddleware(http.HandlerFunc(ChannelPage)))
	mux.Handle("GET /channels", MainMiddleware(http.HandlerFunc(ChannelsPage)))
	mux.Handle("GET /category/{category}", MainMiddleware(http.HandlerFunc(CategoryPage)))
	mux.Handle("GET /categories", MainMiddleware(http.HandlerFunc(CategoriesPage)))
	mux.Handle("GET /clips", MainMiddleware(http.HandlerFunc(ClipsPage)))

	//------------------ API ------------------
	mux.Handle("GET /api/v1/index", APIMiddleware(http.HandlerFunc(IndexAPI)))
	mux.Handle("GET /api/v1/chart/stats", APIMiddleware(http.HandlerFunc(ChartStatsAPI)))

	server := &http.Server{
		Handler:      mux,
		Addr:         ":8011",
		ReadTimeout:  7 * time.Second,
		WriteTimeout: 7 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
