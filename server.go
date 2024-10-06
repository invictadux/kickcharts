package main

import (
	"context"
	"encoding/json"
	"html/template"
	"invictadux/code/db"
	"invictadux/code/funcmaps"
	"invictadux/code/scraper"
	"log"
	"net/http"
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

	page := map[string]interface{}{}
	page["ChannelsChart"] = db.GetOverallChannelsGraph(t1, t2)
	page["ViewsChart"] = db.GetOverallViewsGraph(t1, t2)
	page["Categories"], _ = db.GetCategories(0, 14)
	page["Channels"], _ = db.GetChannels(0, 16)
	page["Clips"], _ = db.GetClips(0, 6)

	indexTemplate.Execute(w, page)
}

func ChannelPage(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")

	now := time.Now().UTC()
	t1 := now.Add(-time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	t2 := now.Add(-time.Hour).Format("2006-01-02 15:04:05")
	t1Followers := now.Add(-time.Hour * 24 * 30).Format("2006-01-02 15:04:05")
	t2Followers := now.Add(-time.Hour * 24).Format("2006-01-02 15:04:05")

	page := map[string]interface{}{}
	page["Channel"], _ = db.GetChannel(channel)
	followersGraph := db.GetChannelFollowersGraph(channel, t1Followers, t2Followers)
	followersTable := followersGraph.ToTable()
	page["FollowersChart"] = followersGraph
	page["FollowersTable"] = followersTable
	page["ViewsChart"] = db.GetChannelViewsGraph(channel, t1, t2)
	page["Clips"], _ = db.GetChannelClips(channel, 0, 6)

	channelTemplate.Execute(w, page)
}

func ChannelsPage(w http.ResponseWriter, r *http.Request) {
	page := map[string]interface{}{}
	page["Channels"], _ = db.GetChannels(0, 60)
	channelsTemplate.Execute(w, page)
}

func CategoryPage(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")

	now := time.Now().UTC()
	t1 := now.Add(-time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	t2 := now.Add(-time.Hour).Format("2006-01-02 15:04:05")

	page := map[string]interface{}{}
	page["Category"], _ = db.GetCategory(category)
	page["ChannelsChart"] = db.GetCategoryChannelsGraph(category, t1, t2)
	page["ViewsChart"] = db.GetCategoryViewsGraph(category, t1, t2)
	page["Clips"], _ = db.GetCategoryClips(category, 0, 6)

	categoryTemplate.Execute(w, page)
}

func CategoriesPage(w http.ResponseWriter, r *http.Request) {
	page := map[string]interface{}{}
	page["Categories"], _ = db.GetCategories(0, 60)
	categoriesTemplate.Execute(w, page)
}

func ClipsPage(w http.ResponseWriter, r *http.Request) {
	page := map[string]interface{}{}
	page["Clips"], _ = db.GetClips(0, 28)
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

func ChannelsAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var err error
	page := 1

	if params.Has("page") {
		page, err = strconv.Atoi(params.Get("page"))

		if err != nil {
			ClientError(w, r, 400)
			return
		}
	}

	limit := 30
	offset := (page - 1) * limit

	channels, err := db.GetChannels(offset, limit)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	json.NewEncoder(w).Encode(channels)
}

func CategoriesAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var err error
	page := 1

	if params.Has("page") {
		page, err = strconv.Atoi(params.Get("page"))

		if err != nil {
			ClientError(w, r, 400)
			return
		}
	}

	limit := 30
	offset := (page - 1) * limit

	categories, err := db.GetCategories(offset, limit)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func ClipsAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var err error
	page := 1

	if params.Has("page") {
		page, err = strconv.Atoi(params.Get("page"))

		if err != nil {
			ClientError(w, r, 400)
			return
		}
	}

	limit := 30
	offset := (page - 1) * limit

	clips, err := db.GetClips(offset, limit)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	json.NewEncoder(w).Encode(clips)
}

func main() {
	mux := http.NewServeMux()
	db.Init()
	go scraper.Run()

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
	mux.Handle("GET /api/v1/channels", APIMiddleware(http.HandlerFunc(ChannelsAPI)))
	mux.Handle("GET /api/v1/categories", APIMiddleware(http.HandlerFunc(CategoriesAPI)))
	mux.Handle("GET /api/v1/clips", APIMiddleware(http.HandlerFunc(ClipsAPI)))

	server := &http.Server{
		Handler:      mux,
		Addr:         ":8005",
		ReadTimeout:  7 * time.Second,
		WriteTimeout: 7 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
