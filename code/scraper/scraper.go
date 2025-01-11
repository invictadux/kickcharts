package scraper

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"invictadux/code/db"
	"invictadux/code/kickapi"
	"invictadux/code/models"
)

var maxChannelViewerCount = 5

func GetAllLivestreams() *chan (string) {
	params := url.Values{
		"page":  []string{"1"},
		"limit": []string{"32"},
		"sort":  []string{"desc"},
	}

	channels := make(chan string, 10000)
	run := true
	page := 1

	for run {
		for {
			data, err := kickapi.GetLivestreams(params)

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
				continue
			}

			fmt.Printf("Page: %v, len: %v\n", params["page"], len(data.Data))

			for _, stream := range data.Data {
				if stream.ViewerCount >= maxChannelViewerCount {
					channels <- stream.Channel.Slug
				} else {
					run = false
					break
				}
			}

			page++
			params["page"] = []string{strconv.Itoa(page)}
			break
		}
	}

	close(channels)
	return &channels
}

func GetAllCategories() *chan (string) {
	params := url.Values{
		"page":  []string{"1"},
		"limit": []string{"32"},
	}

	categories := make(chan string, 10000)
	run := true
	page := 1
	totalViewers := 0

	for run {
		for {
			data, err := kickapi.GetCategories(params)

			if data.CurrentPage != page {
				run = false
				break
			}

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
				continue
			}

			fmt.Printf("Page: %v, len: %v\n", params["page"], len(data.Data))

			for _, category := range data.Data {
				if category.Viewers > 0 {
					categories <- category.Slug

					c := models.Category{
						Name:         category.Name,
						Slug:         category.Slug,
						Banner:       category.Banner.URL,
						LiveViewers:  category.Viewers,
						LiveChannels: 0,
						PeakViewers:  category.Viewers,
						PeakChannels: 0,
						Description:  category.Description,
					}

					dbCategory, err := db.GetCategory(category.Slug)

					if err != nil {
						db.InsertCategory(c)
					} else {
						if dbCategory.PeakViewers > c.PeakViewers {
							c.PeakViewers = dbCategory.PeakViewers
						}

						updateParams := map[string]interface{}{
							"name":          c.Name,
							"banner":        c.Banner,
							"live_viewers":  c.LiveViewers,
							"peak_channels": c.PeakChannels,
							"peak_viewers":  c.PeakViewers,
							"description":   c.Description,
						}

						db.UpdateCategory(category.Slug, updateParams)
					}

					db.InsertCategoryViewsChartPoint(category.Slug, category.Viewers)
					totalViewers += category.Viewers
				} else {
					run = false
					break
				}
			}

			page++
			params["page"] = []string{strconv.Itoa(page)}
			break
		}
	}

	close(categories)
	db.InsertOverallViewersChartPoint(totalViewers)
	return &categories
}

func GetAllCategoriesLivestreams(chSlugs *chan string) []string {
	totalChannels := 0
	liveChannels := []string{}

	for {
		slug, ok := <-*chSlugs

		if !ok {
			break
		}

		params := url.Values{
			"page":        []string{"1"},
			"limit":       []string{"32"},
			"subcategory": []string{slug},
		}

		page := 1
		totalStreams := 0

		//Get all category streams -------------------------------------------------------
		for {
			data, err := kickapi.GetCategoryLivestreams(params)
			totalChannels += len(data.Data)
			totalStreams += len(data.Data)

			if data.CurrentPage != page {
				break
			}

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * 30)
				continue
			}

			for _, livestream := range data.Data {
				if livestream.Viewers > maxChannelViewerCount {
					liveChannels = append(liveChannels, livestream.Channel.Slug)

					channel := models.Channel{
						Username:       livestream.Channel.User.Username,
						Slug:           livestream.Channel.Slug,
						Banner:         "",
						Picture:        livestream.Channel.User.Profilepic,
						IsBanned:       livestream.Channel.IsBanned,
						Language:       livestream.Language,
						Live:           livestream.IsLive,
						LiveViewers:    livestream.Viewers,
						FollowersCount: 0,
						PeakViewers:    livestream.Viewers,
						Description:    livestream.Channel.User.Bio,
						Discord:        livestream.Channel.User.Discord,
						Facebook:       livestream.Channel.User.Facebook,
						Instagram:      livestream.Channel.User.Instagram,
						Tiktok:         livestream.Channel.User.Tiktok,
						Twitter:        livestream.Channel.User.Twitter,
						Youtube:        livestream.Channel.User.Youtube,
					}

					dbChannel, err := db.GetChannel(livestream.Channel.Slug)

					if err != nil {
						db.InsertChannel(channel)
					} else {
						if dbChannel.PeakViewers > channel.PeakViewers {
							channel.PeakViewers = dbChannel.PeakViewers
						}

						updateParams := map[string]interface{}{
							"username":     channel.Username,
							"picture":      channel.Picture,
							"is_banned":    channel.IsBanned,
							"live":         channel.Live,
							"live_viewers": channel.LiveViewers,
							"peak_viewers": channel.PeakViewers,
							"description":  channel.Description,
							"discord":      channel.Discord,
							"facebook":     channel.Facebook,
							"instagram":    channel.Instagram,
							"tiktok":       channel.Tiktok,
							"twitter":      channel.Twitter,
							"youtube":      channel.Youtube,
						}

						db.UpdateChannel(dbChannel.ID, updateParams)
					}

					db.InsertChannelViewersChartPoint(livestream.Channel.Slug, livestream.Viewers)
				}
			}

			page++
			params["page"] = []string{strconv.Itoa(page)}

			if len(data.Data) < 32 {
				peakChannels := totalStreams
				cat, _ := db.GetCategory(slug)

				if cat.PeakChannels > peakChannels {
					peakChannels = cat.PeakChannels
				}

				updateParams := map[string]interface{}{
					"live_channels": totalStreams,
					"peak_channels": peakChannels,
				}

				db.UpdateCategory(slug, updateParams)

				db.InsertCategoryLiveChannelsChartPoint(slug, totalStreams)
				break
			}
		}
		//--------------------------------------------------------------------------

		fmt.Println("[+]", slug)
	}

	db.InsertOverallLiveChannelsChartPoint(totalChannels)
	return liveChannels
}

func GetChannelsData(chSlugs *chan string, goroutines int) {
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	count := 0

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			failed := false

			for {
				slug, ok := <-*chSlugs

				if !ok {
					break
				}

				for {
					data, err := kickapi.GetChannel(slug)

					if err != nil {
						if failed {
							failed = false
							break
						}

						failed = true
						fmt.Println("Channel error:", err, slug)
						time.Sleep(time.Second * 30)
						continue
					}

					channel := models.Channel{
						Username:       data.User.Username,
						Slug:           data.Slug,
						Banner:         data.BannerImage.URL,
						Picture:        data.User.ProfilePic,
						Language:       data.Livestream.Language,
						IsBanned:       data.IsBanned,
						FollowersCount: data.FollowersCount,
						PeakViewers:    data.Livestream.ViewerCount,
						Description:    data.User.Bio,
						Discord:        data.User.Discord,
						Facebook:       data.User.Facebook,
						Instagram:      data.User.Instagram,
						Tiktok:         data.User.Tiktok,
						Twitter:        data.User.Twitter,
						Youtube:        data.User.Youtube,
					}

					m.Lock()
					dbChannel, err := db.GetChannel(data.Slug)

					if err != nil {
						db.InsertChannel(channel)
					} else {
						if dbChannel.PeakViewers > channel.PeakViewers {
							channel.PeakViewers = dbChannel.PeakViewers
						}

						updateParams := map[string]interface{}{
							"username":        channel.Username,
							"picture":         channel.Picture,
							"is_banned":       channel.IsBanned,
							"followers_count": channel.FollowersCount,
							"peak_viewers":    channel.PeakViewers,
							"description":     channel.Description,
							"discord":         channel.Discord,
							"facebook":        channel.Facebook,
							"instagram":       channel.Instagram,
							"tiktok":          channel.Tiktok,
							"twitter":         channel.Twitter,
							"youtube":         channel.Youtube,
						}

						db.UpdateChannel(dbChannel.ID, updateParams)
					}

					db.InsertChannelFollowersChartPoint(data.Slug, data.FollowersCount)

					count++
					fmt.Printf("#%v, Username: %v\n", count, data.User.Username)
					m.Unlock()

					break
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func GetAllClips() {
	params := url.Values{
		"sort": []string{"view"},
		"time": []string{"day"},
	}

	run := true

	for run {
		for {
			data, err := kickapi.GetClips(params)

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
				continue
			}

			fmt.Printf("len: %v\n", len(data.Clips))

			for _, clip := range data.Clips {
				if clip.ViewCount >= 100 {
					c := models.Clip{
						ID:           clip.ID,
						CategoryName: clip.Category.Name,
						CategorySlug: clip.Category.Slug,
						ChannelName:  clip.Channel.Username,
						ChannelSlug:  clip.Channel.Slug,
						IsMature:     clip.IsMature,
						Title:        clip.Title,
						URL:          clip.ClipURL,
						Likes:        clip.Likes,
						LivestreamID: clip.LivestreamID,
						Thumbnail:    clip.ThumbnailURL,
						Views:        clip.ViewCount,
						Duration:     clip.Duration,
						CreatedAt:    clip.CreatedAt,
					}

					db.InsertClip(c)
				} else {
					run = false
					break
				}
			}

			params["cursor"] = []string{data.NextCursor}
			break
		}
	}
}

func Run() {
	offset, limit := 0, 5000

	for {
		now := time.Now().UTC()

		if now.Minute() == 0 {
			start := time.Now()
			//chSlug := GetAllLivestreams()
			//GetChannelsData(chSlug, 5)

			categories := GetAllCategories()
			liveChannels := GetAllCategoriesLivestreams(categories)
			db.SetChannelsLiveStatus(&liveChannels)
			GetAllClips()

			if now.Hour() == 0 {
				offset = 0
			}

			chSlug := db.GetChannelsSlug(offset, limit)
			GetChannelsData(chSlug, 5)
			offset += limit
			fmt.Println("Elapsed time:", time.Since(start))
		}

		time.Sleep(time.Minute)
	}
}
