package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"github.com/google/uuid"
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
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
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

func (cfg *apiConfig) UpdateFeedData(amount int) {
	fmt.Println("Starting feed update cycle")
	feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), int32(amount))
	if err != nil {
		fmt.Printf("Error getting feeds from db to update: %s\n", err)
		time.Sleep(30)
		return
	}
	if len(feeds) < 1 {
		fmt.Printf("No feeds to fetch, waiting %d seconds before checking again.\n", cfg.FeedUpdateIntervalSeconds)
		return
	}
	wg := sync.WaitGroup{}
	for _, val := range feeds {
		wg.Add(1)
		go func(oldFeed database.Feed, c apiConfig) {
			defer wg.Done()
			feed, err := FetchRssFeed(oldFeed.Url)
			if err != nil {
				fmt.Printf("There was an error fetching RSS feed at %s: %s\n", oldFeed.Url, err)
				return
			}
			fmt.Printf("Processing feed: %s. Channel Title: %s\n", oldFeed.Name, feed.Channel.Title)
			err = c.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
				ID:        oldFeed.ID,
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				fmt.Printf("There was an error marking feed %s as fetched\n", oldFeed.ID)
				return
			}
			for _, val := range feed.Channel.Item {
				postParams := database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
					Title:       val.Title,
					Url:         val.Link,
					Description: val.Description,
					PublishedAt: val.PubDate,
					FeedID:      oldFeed.ID,
				}
				query, err := cfg.DB.CreatePost(context.Background(), postParams)
				if err != nil {
					//Dont wanna deal with checking for existing post content in a feed, just let postgres return an error
					if strings.ContainsAny(err.Error(), "duplicate key value violates unique constraint") {
						continue
					} else {
						fmt.Printf("Error inserting post into DB: %s\n", err)
					}
				}
				fmt.Printf("Added post to feed %s:\n%v\n", oldFeed.ID, query)
			}
			return
		}(val, *cfg)
	}
	wg.Wait()
	return
}

func RefreshFeeds(c apiConfig) {
	ticker := time.NewTicker(time.Second * time.Duration(c.FeedUpdateIntervalSeconds))
	for {
		<-ticker.C
		c.UpdateFeedData(c.FeedRefreshAmount)
	}
}
