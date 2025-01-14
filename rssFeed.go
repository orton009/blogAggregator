package main

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

type FeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type RSSFeed struct {
	Channel struct {
		Title         string     `xml:"title"`
		Link          string     `xml:"link"`
		Description   string     `xml:"description"`
		Generator     string     `xml:"generator"`
		Language      string     `xml:"language"`
		LastBuildDate string     `xml:"lastBuildDate"`
		Items         []FeedItem `xml:"item"`
	} `xml:"channel"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return &rssFeed, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &rssFeed, err
	}

	decoder := xml.NewDecoder(resp.Body)
	decoder.Strict = false
	decoder.Entity = xml.HTMLEntity

	err = decoder.Decode(&rssFeed)
	if err != nil {
		return &rssFeed, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Items {

		rssFeed.Channel.Items[i].Title = html.UnescapeString(rssFeed.Channel.Items[i].Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(rssFeed.Channel.Items[i].Description)
	}

	return &rssFeed, nil
}
