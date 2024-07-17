package main

import (
	"encoding/xml"
	"net/http"
)

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchRssFeed(url string) (RssFeed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return RssFeed{}, err
	}

	decoder := xml.NewDecoder(resp.Body)
	defer resp.Body.Close()

	feed := RssFeed{}
	err = decoder.Decode(&feed)
	if err != nil {
		return RssFeed{}, err
	}
	return feed, nil
}
