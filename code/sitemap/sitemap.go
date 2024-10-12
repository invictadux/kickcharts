package sitemap

import (
	"encoding/xml"
	"fmt"
	"invictadux/code/db"
	"os"
)

type URL struct {
	Loc      string  `xml:"loc"`
	Priority float64 `xml:"priority"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

func Generate() {
	channels := db.GetAllChannels()
	categories := db.GetAllCategories()

	urls := []URL{}
	urls = append(urls, URL{Loc: "https://kickcharts.net", Priority: 1})
	urls = append(urls, URL{Loc: "https://kickcharts.net/channels", Priority: 0.9})
	urls = append(urls, URL{Loc: "https://kickcharts.net/categories", Priority: 0.9})
	urls = append(urls, URL{Loc: "https://kickcharts.net/clips", Priority: 0.9})

	for _, channel := range channels {
		url := URL{}
		url.Loc = "https://kickcharts.net/channel/" + channel
		url.Priority = 0.7
		urls = append(urls, url)
	}

	for _, category := range categories {
		url := URL{}
		url.Loc = "https://kickcharts.net/category/" + category
		url.Priority = 0.7
		urls = append(urls, url)
	}

	sitemap := Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	xmlData, err := xml.MarshalIndent(sitemap, "", "    ")
	if err != nil {
		fmt.Println("Error encoding XML:", err)
		return
	}

	file, err := os.Create("sitemap.xml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(xml.Header))
	if err != nil {
		fmt.Println("Error writing XML header:", err)
		return
	}

	_, err = file.Write(xmlData)
	if err != nil {
		fmt.Println("Error writing XML data:", err)
		return
	}
}
