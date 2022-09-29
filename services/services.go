package services

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
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
	RssFeed Channel  `xml:"channel"`
}

type XkcdResponse struct {
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
	url          string
	xkcdresponse XkcdResponse
	xmlresponse  RssFeed
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
	var itemsdata []ItemData = []ItemData{}
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
		itemsdata = append(itemsdata, data...)
	}

	sort.Slice(itemsdata, func(i, j int) bool {
		return itemsdata[i].PublishingDate.After(itemsdata[j].PublishingDate)
	})

	return itemsdata, nil
}

func getXKCD() ([]ItemData, error) {
	var itemsdata []ItemData
	for i := 0; i < 10; i++ {
		if i == 0 {
			url = XKCD + XKCD_CONFIG
		} else {
			url = XKCD + `/` + strconv.Itoa(i) + XKCD_CONFIG
		}
		resp, err := http.Get(url)
		if err != nil {
			return itemsdata, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return itemsdata, err
		}
		err = json.Unmarshal([]byte(string(body)), &xkcdresponse)
		if err != nil {
			return itemsdata, err
		}
		date := xkcdresponse.Year + `-` + xkcdresponse.Month + `-` + xkcdresponse.Day + ` 00:00:00 +0000`
		t, err := timefmt.Parse(date, "%Y-%m-%d %T %z")
		if err != nil {
			return itemsdata, err
		}
		data = ItemData{
			PictureUrl:     xkcdresponse.Img,
			Title:          xkcdresponse.Title,
			Description:    xkcdresponse.News,
			WebUrl:         xkcdresponse.Link,
			PublishingDate: t,
		}
		itemsdata = append(itemsdata, data)
	}
	return itemsdata, nil
}

func getPoorlyDrawnLines() ([]ItemData, error) {
	var itemsdata []ItemData
	resp, err := http.Get(POORLY_DRAWN_LINES)
	if err != nil {
		return itemsdata, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return itemsdata, err
	}
	err = xml.Unmarshal([]byte(string(body)), &xmlresponse)
	if err != nil {
		return itemsdata, err
	}
	for i := 0; i < len(xmlresponse.RssFeed.Item); i++ {
		date := xmlresponse.RssFeed.Item[i].PubDate
		t, err := timefmt.Parse(date, "%a, %d %b %Y %T %z")
		if err != nil {
			return itemsdata, err
		}
		data = ItemData{
			PictureUrl:     xmlresponse.RssFeed.Item[i].Link,
			Title:          xmlresponse.RssFeed.Item[i].Title,
			Description:    xmlresponse.RssFeed.Item[i].Description,
			WebUrl:         POORLY_DRAWN_LINES,
			PublishingDate: t,
		}
		itemsdata = append(itemsdata, data)
	}
	return itemsdata, nil
}
