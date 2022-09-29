package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/itchyny/timefmt-go"
)

const (
	POORLY_DRAWN_LINES = `http://feeds.feedburner.com/PoorlyDrawnLines`
	XKCD               = `https://xkcd.com`
)

//	type XkcdResponses struct {
//		XkcdResponses []XkcdResponse
//	}
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
	ItemsData    []ItemData
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

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		getXKCD()
	}()
	go func() {
		defer wg.Done()
		getPoorlyDrawnLines()
	}()
	wg.Wait()
	fmt.Println(ItemsData)
	sort.Slice(ItemsData, func(i, j int) bool {
		return ItemsData[i].PublishingDate.Before(ItemsData[j].PublishingDate)
	})
	fmt.Println(len(ItemsData))

}

func getXKCD() {
	for i := 0; i < 10; i++ {
		if i == 0 {
			url = XKCD + "/info.0.json"
		} else {
			url = XKCD + `/` + strconv.Itoa(i) + "/info.0.json"
		}
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		json.Unmarshal([]byte(string(body)), &xkcdresponse)
		date := xkcdresponse.Year + `-` + xkcdresponse.Month + `-` + xkcdresponse.Day + ` 00:00:00 +0000`
		t, err := timefmt.Parse(date, "%Y-%m-%d %T %z")
		if err != nil {
			log.Fatal(err)
		}
		data = ItemData{
			PictureUrl:     xkcdresponse.Img,
			Title:          xkcdresponse.Title,
			Description:    xkcdresponse.News,
			WebUrl:         xkcdresponse.Link,
			PublishingDate: t,
		}
		ItemsData = append(ItemsData, data)
	}
}
func getPoorlyDrawnLines() {
	resp, err := http.Get(POORLY_DRAWN_LINES)
	if err != nil {
		// handle error
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = xml.Unmarshal([]byte(string(body)), &xmlresponse)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(xmlresponse.RssFeed.Item); i++ {
		date := xmlresponse.RssFeed.Item[i].PubDate
		t, err := timefmt.Parse(date, "%a, %d %b %Y %T %z")
		if err != nil {
			log.Fatal(err)
		}
		data = ItemData{
			PictureUrl:     xmlresponse.RssFeed.Item[i].Link,
			Title:          xmlresponse.RssFeed.Item[i].Title,
			Description:    xmlresponse.RssFeed.Item[i].Description,
			WebUrl:         POORLY_DRAWN_LINES,
			PublishingDate: t,
		}
		ItemsData = append(ItemsData, data)
	}
}
