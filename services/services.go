package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/itchyny/timefmt-go"
)

const (
	POORLY_DRAWN_LINES = `http://feeds.feedburner.com/PoorlyDrawnLines`
	XKCD               = `https://xkcd.com`
	XKCD_CONFIG        = "/info.0.json"
)

type Item struct {
	XMLName        xml.Name `xml:"item"`
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Creator        string   `xml:"dc:creator"`
	PubDate        string   `xml:"pubDate"`
	Category       string   `xml:"category"`
	Guid           string   `xml:"guid"`
	Description    string   `xml:"description"`
	ContentEncoded string   `xml:"content:encoded"`
}
type Channel struct {
	XMLName           xml.Name `xml:"channel"`
	Title             string   `xml:"title"`
	Atom              string   `xml:"atom"`
	Link              string   `xml:"link"`
	Description       string   `xml:"description"`
	LastBuildDate     string   `xml:"lastBuildDate"`
	Language          string   `xml:"language"`
	SyUpdatePeriod    string   `xml:"sy:updatePeriod"`
	SyUpdateFrequency string   `xml:"sy:updateFrequency"`
	Generator         string   `xml:"generator"`
	Item              []Item   `xml:"item"`
}
type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	RSSFeed Channel  `xml:"channel"`
}

type XKCDResponse struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

var (
	URL          string
	xkcdResponse XKCDResponse
	xmlResponse  RssFeed
	data         ItemData
)

type ItemData struct {
	PictureUrl     string    `json:"picurl,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description"`
	WebUrl         string    `json:"weburl"`
	PublishingDate time.Time `json:"publishing_date"`
}

func GetData() ([]ItemData, error) {
	var errCh chan error = make(chan error, 2)
	var dataCh chan []ItemData = make(chan []ItemData, 2)
	var itemsData []ItemData = []ItemData{}
	defer close(errCh)
	defer close(dataCh)

	go func(data chan []ItemData, ch chan error) {
		items, err := getXKCD()
		data <- items
		ch <- err
	}(dataCh, errCh)

	go func(data chan []ItemData, ch chan error) {
		items, err := getPoorlyDrawnLines()
		data <- items
		ch <- err
	}(dataCh, errCh)

	for i := 0; i < cap(errCh); i += 1 {
		err := <-errCh
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < cap(dataCh); i += 1 {
		data := <-dataCh
		itemsData = append(itemsData, data...)
	}

	sort.Slice(itemsData, func(i, j int) bool {
		return itemsData[i].PublishingDate.After(itemsData[j].PublishingDate)
	})

	return itemsData, nil
}

func getXKCD() ([]ItemData, error) {
	var itemsData []ItemData
	for i := 0; i < 10; i++ {
		if i == 0 {
			URL = fmt.Sprintf(`%s%s`, XKCD, XKCD_CONFIG)
		} else {
			URL = fmt.Sprintf(`%s/%d%s`, XKCD, i, XKCD_CONFIG)
		}
		resp, err := http.Get(URL)
		if err != nil {
			return itemsData, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return itemsData, err
		}
		err = json.Unmarshal([]byte(string(body)), &xkcdResponse)
		if err != nil {
			return itemsData, err
		}
		date := fmt.Sprintf(`%s-%s-%s %s`, xkcdResponse.Year, xkcdResponse.Month, xkcdResponse.Day, ` 00:00:00 +0000`)
		publishingDate, err := timefmt.Parse(date, "%Y-%m-%d %T %z")
		if err != nil {
			return itemsData, err
		}
		itemsData = append(itemsData, ItemData{
			PictureUrl:     xkcdResponse.Img,
			Title:          xkcdResponse.Title,
			Description:    xkcdResponse.News,
			WebUrl:         xkcdResponse.Link,
			PublishingDate: publishingDate,
		})
	}
	return itemsData, nil
}

func getPoorlyDrawnLines() ([]ItemData, error) {
	var itemsData []ItemData
	resp, err := http.Get(POORLY_DRAWN_LINES)
	if err != nil {
		return itemsData, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return itemsData, err
	}
	err = xml.Unmarshal([]byte(string(body)), &xmlResponse)
	if err != nil {
		return itemsData, err
	}
	for i := 0; i < len(xmlResponse.RSSFeed.Item); i++ {
		date := xmlResponse.RSSFeed.Item[i].PubDate
		t, err := timefmt.Parse(date, "%a, %d %b %Y %T %z")
		if err != nil {
			return itemsData, err
		}
		data = ItemData{
			PictureUrl:     xmlResponse.RSSFeed.Item[i].Link,
			Title:          xmlResponse.RSSFeed.Item[i].Title,
			Description:    xmlResponse.RSSFeed.Item[i].Description,
			WebUrl:         POORLY_DRAWN_LINES,
			PublishingDate: t,
		}
		itemsData = append(itemsData, data)
	}
	return itemsData, nil
}
