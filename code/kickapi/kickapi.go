package kickapi

import (
	"encoding/json"
	"io"
	"net/url"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Livestreams struct {
	CurrentPage  int          `json:"current_page"`
	Data         []Livestream `json:"data"`
	FirstPageURL string       `json:"first_page_url"`
	From         int          `json:"from"`
	NextPageURL  string       `json:"next_page_url"`
	Path         string       `json:"path"`
	PerPage      int          `json:"per_page"`
	PrevPageURL  string       `json:"prev_page_url"`
	To           int          `json:"to"`
}

type Livestream struct {
	ID            int         `json:"id"`
	Slug          string      `json:"slug"`
	ChannelID     int         `json:"channel_id"`
	CreatedAt     string      `json:"created_at"`
	SessionTitle  string      `json:"session_title"`
	IsLive        bool        `json:"is_live"`
	RiskLevelID   interface{} `json:"risk_level_id"`
	StartTime     string      `json:"start_time"`
	Source        interface{} `json:"source"`
	TwitchChannel interface{} `json:"twitch_channel"`
	Duration      int         `json:"duration"`
	Language      string      `json:"language"`
	IsMature      bool        `json:"is_mature"`
	ViewerCount   int         `json:"viewer_count"`
	Order         int         `json:"order"`
	Thumbnail     struct {
		Srcset string `json:"srcset"`
		Src    string `json:"src"`
	} `json:"thumbnail"`
	Viewers int `json:"viewers"`
	Channel struct {
		ID                  int    `json:"id"`
		UserID              int    `json:"user_id"`
		Slug                string `json:"slug"`
		IsBanned            bool   `json:"is_banned"`
		PlaybackURL         string `json:"playback_url"`
		NameUpdatedAt       string `json:"name_updated_at"`
		VodEnabled          bool   `json:"vod_enabled"`
		SubscriptionEnabled bool   `json:"subscription_enabled"`
		CanHost             bool   `json:"can_host"`
		User                struct {
			ID              int       `json:"id"`
			Username        string    `json:"username"`
			AgreedToTerms   bool      `json:"agreed_to_terms"`
			EmailVerifiedAt time.Time `json:"email_verified_at"`
			Bio             string    `json:"bio"`
			Country         string    `json:"country"`
			State           string    `json:"state"`
			City            string    `json:"city"`
			Instagram       string    `json:"instagram"`
			Twitter         string    `json:"twitter"`
			Youtube         string    `json:"youtube"`
			Discord         string    `json:"discord"`
			Tiktok          string    `json:"tiktok"`
			Facebook        string    `json:"facebook"`
			Profilepic      string    `json:"profilepic"`
		} `json:"user"`
	} `json:"channel"`
	Categories []struct {
		ID          int         `json:"id"`
		CategoryID  int         `json:"category_id"`
		Name        string      `json:"name"`
		Slug        string      `json:"slug"`
		Tags        []string    `json:"tags"`
		Description interface{} `json:"description"`
		DeletedAt   interface{} `json:"deleted_at"`
		Viewers     int         `json:"viewers"`
		Category    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
			Icon string `json:"icon"`
		} `json:"category"`
	} `json:"categories"`
}

type Categories struct {
	CurrentPage int `json:"current_page"`
	Data        []struct {
		ID          int      `json:"id"`
		CategoryID  int      `json:"category_id"`
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Tags        []string `json:"tags"`
		Description string   `json:"description"`
		DeletedAt   any      `json:"deleted_at"`
		Viewers     int      `json:"viewers"`
		Banner      struct {
			Responsive string `json:"responsive"`
			URL        string `json:"url"`
		} `json:"banner"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
			Icon string `json:"icon"`
		} `json:"category"`
	} `json:"data"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	Links        []struct {
		URL    any    `json:"url"`
		Label  string `json:"label"`
		Active bool   `json:"active"`
	} `json:"links"`
	NextPageURL string `json:"next_page_url"`
	Path        string `json:"path"`
	PerPage     string `json:"per_page"`
	PrevPageURL any    `json:"prev_page_url"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

type Subcategory struct {
	ID             int      `json:"id"`
	CategoryID     int      `json:"category_id"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description"`
	DeletedAt      any      `json:"deleted_at"`
	Viewers        int      `json:"viewers"`
	FollowersCount int      `json:"followers_count"`
	Followed       bool     `json:"followed"`
	Banner         struct {
		Srcset string `json:"srcset"`
		Src    string `json:"src"`
	} `json:"banner"`
}

type Channel struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"user_id"`
	Slug                string `json:"slug"`
	IsBanned            bool   `json:"is_banned"`
	PlaybackURL         string `json:"playback_url"`
	VodEnabled          bool   `json:"vod_enabled"`
	SubscriptionEnabled bool   `json:"subscription_enabled"`
	FollowersCount      int    `json:"followers_count"`
	Following           bool   `json:"following"`
	Subscription        any    `json:"subscription"`
	SubscriberBadges    []struct {
		ID         int `json:"id"`
		ChannelID  int `json:"channel_id"`
		Months     int `json:"months"`
		BadgeImage struct {
			Srcset string `json:"srcset"`
			Src    string `json:"src"`
		} `json:"badge_image"`
	} `json:"subscriber_badges"`
	BannerImage struct {
		URL string `json:"url"`
	} `json:"banner_image"`
	Livestream struct {
		ID            int    `json:"id"`
		Slug          string `json:"slug"`
		ChannelID     int    `json:"channel_id"`
		CreatedAt     string `json:"created_at"`
		SessionTitle  string `json:"session_title"`
		IsLive        bool   `json:"is_live"`
		RiskLevelID   any    `json:"risk_level_id"`
		StartTime     string `json:"start_time"`
		Source        any    `json:"source"`
		TwitchChannel any    `json:"twitch_channel"`
		Duration      int    `json:"duration"`
		Language      string `json:"language"`
		IsMature      bool   `json:"is_mature"`
		ViewerCount   int    `json:"viewer_count"`
		Thumbnail     struct {
			URL string `json:"url"`
		} `json:"thumbnail"`
		Categories []struct {
			ID          int      `json:"id"`
			CategoryID  int      `json:"category_id"`
			Name        string   `json:"name"`
			Slug        string   `json:"slug"`
			Tags        []string `json:"tags"`
			Description any      `json:"description"`
			DeletedAt   any      `json:"deleted_at"`
			Viewers     int      `json:"viewers"`
			Category    struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
				Icon string `json:"icon"`
			} `json:"category"`
		} `json:"categories"`
		Tags []any `json:"tags"`
	} `json:"livestream"`
	Role               any   `json:"role"`
	Muted              bool  `json:"muted"`
	FollowerBadges     []any `json:"follower_badges"`
	OfflineBannerImage struct {
		Src    string `json:"src"`
		Srcset string `json:"srcset"`
	} `json:"offline_banner_image"`
	Verified         bool `json:"verified"`
	RecentCategories []struct {
		ID          int      `json:"id"`
		CategoryID  int      `json:"category_id"`
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Tags        []string `json:"tags"`
		Description any      `json:"description"`
		DeletedAt   any      `json:"deleted_at"`
		Viewers     int      `json:"viewers"`
		Banner      struct {
			Responsive string `json:"responsive"`
			URL        string `json:"url"`
		} `json:"banner"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
			Icon string `json:"icon"`
		} `json:"category"`
	} `json:"recent_categories"`
	CanHost bool `json:"can_host"`
	User    struct {
		ID              int       `json:"id"`
		Username        string    `json:"username"`
		AgreedToTerms   bool      `json:"agreed_to_terms"`
		EmailVerifiedAt time.Time `json:"email_verified_at"`
		Bio             string    `json:"bio"`
		Country         string    `json:"country"`
		State           string    `json:"state"`
		City            string    `json:"city"`
		Instagram       string    `json:"instagram"`
		Twitter         string    `json:"twitter"`
		Youtube         string    `json:"youtube"`
		Discord         string    `json:"discord"`
		Tiktok          string    `json:"tiktok"`
		Facebook        string    `json:"facebook"`
		ProfilePic      string    `json:"profile_pic"`
	} `json:"user"`
	Chatroom struct {
		ID                   int       `json:"id"`
		ChatableType         string    `json:"chatable_type"`
		ChannelID            int       `json:"channel_id"`
		CreatedAt            time.Time `json:"created_at"`
		UpdatedAt            time.Time `json:"updated_at"`
		ChatModeOld          string    `json:"chat_mode_old"`
		ChatMode             string    `json:"chat_mode"`
		SlowMode             bool      `json:"slow_mode"`
		ChatableID           int       `json:"chatable_id"`
		FollowersMode        bool      `json:"followers_mode"`
		SubscribersMode      bool      `json:"subscribers_mode"`
		EmotesMode           bool      `json:"emotes_mode"`
		MessageInterval      int       `json:"message_interval"`
		FollowingMinDuration int       `json:"following_min_duration"`
	} `json:"chatroom"`
}

type Clips struct {
	Clips []struct {
		ID           string    `json:"id"`
		LivestreamID string    `json:"livestream_id"`
		CategoryID   string    `json:"category_id"`
		ChannelID    int       `json:"channel_id"`
		UserID       int       `json:"user_id"`
		Title        string    `json:"title"`
		ClipURL      string    `json:"clip_url"`
		ThumbnailURL string    `json:"thumbnail_url"`
		Privacy      string    `json:"privacy"`
		Likes        int       `json:"likes"`
		Liked        bool      `json:"liked"`
		Views        int       `json:"views"`
		Duration     int       `json:"duration"`
		StartedAt    time.Time `json:"started_at"`
		CreatedAt    time.Time `json:"created_at"`
		IsMature     bool      `json:"is_mature"`
		VideoURL     string    `json:"video_url"`
		ViewCount    int       `json:"view_count"`
		LikesCount   int       `json:"likes_count"`
		Category     struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Slug           string `json:"slug"`
			Responsive     string `json:"responsive"`
			Banner         string `json:"banner"`
			ParentCategory string `json:"parent_category"`
		} `json:"category"`
		Creator struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Slug           string `json:"slug"`
			ProfilePicture any    `json:"profile_picture"`
		} `json:"creator"`
		Channel struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Slug           string `json:"slug"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"channel"`
	} `json:"clips"`
	NextCursor string `json:"nextCursor"`
}

func NewRequest(urlPath string, params url.Values) (*[]byte, error) {
	var bytes []byte

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Firefox_123),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return &bytes, err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)

	if err != nil {
		return &bytes, err
	}

	req.URL.RawQuery = params.Encode()

	req.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8"},
		"Accept-Encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.5"},
		"Cache-Control":             {"no-cache"},
		"Connection":                {"keep-alive"},
		"Host":                      {"kick.com"},
		"Pragma":                    {"no-cache"},
		"Priority":                  {"u=0, i"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0"},
		http.HeaderOrderKey: {
			"accept",
			"Accept-Encoding",
			"accept-language",
			"Cache-Control",
			"Connection",
			"Host",
			"Pragma",
			"Priority",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"Sec-Fetch-User",
			"Upgrade-Insecure-Requests",
			"user-agent",
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		return &bytes, err
	}

	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return &bytes, err
	}

	return &bytes, nil
}

func GetLivestreams(params url.Values) (Livestreams, error) {
	endpoints := url.Values{
		"page":  []string{"1"},
		"limit": []string{"24"},
		"sort":  []string{"desc"},
	}

	for key, param := range params {
		endpoints[key] = param
	}

	var livestreams Livestreams
	bytes, err := NewRequest("https://kick.com/stream/livestreams/en", endpoints)

	if err != nil {
		return livestreams, err
	}

	err = json.Unmarshal(*bytes, &livestreams)
	return livestreams, err
}

func GetCategories(params url.Values) (Categories, error) {
	endpoints := url.Values{
		"page":  []string{"1"},
		"limit": []string{"24"},
	}

	for key, param := range params {
		endpoints[key] = param
	}

	var categories Categories
	bytes, err := NewRequest("https://kick.com/api/v1/subcategories", endpoints)

	if err != nil {
		return categories, err
	}

	err = json.Unmarshal(*bytes, &categories)
	return categories, err
}

func GetCategoryLivestreams(params url.Values) (Livestreams, error) {
	endpoints := url.Values{
		"page":  []string{"1"},
		"limit": []string{"32"},
		"sort":  []string{"desc"},
	}

	for key, param := range params {
		endpoints[key] = param
	}

	var categoriesLivestreams Livestreams
	bytes, err := NewRequest("https://kick.com/stream/livestreams/en", endpoints)

	if err != nil {
		return categoriesLivestreams, err
	}

	err = json.Unmarshal(*bytes, &categoriesLivestreams)
	return categoriesLivestreams, err
}

func GetSubcategories(subcategory string) (Subcategory, error) {
	var subcategoryData Subcategory
	bytes, err := NewRequest("https://kick.com/api/v1/subcategories/"+subcategory, url.Values{})

	if err != nil {
		return subcategoryData, err
	}

	err = json.Unmarshal(*bytes, &subcategoryData)
	return subcategoryData, err
}

func GetChannel(channel string) (Channel, error) {
	var channelData Channel
	bytes, err := NewRequest("https://kick.com/api/v2/channels/"+channel, url.Values{})

	if err != nil {
		return channelData, err
	}

	err = json.Unmarshal(*bytes, &channelData)
	return channelData, err
}

func GetClips(params url.Values) (Clips, error) {
	endpoints := url.Values{
		"sort": []string{"view"},
		"time": []string{"week"},
	}

	for key, param := range params {
		endpoints[key] = param
	}

	var clips Clips
	bytes, err := NewRequest("https://kick.com/api/v2/clips", endpoints)

	if err != nil {
		return clips, err
	}

	err = json.Unmarshal(*bytes, &clips)
	return clips, err
}
