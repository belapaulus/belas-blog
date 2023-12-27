package main

import (
	"bufio"
	"html"
	"html/template"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/mattn/go-runewidth"
	"github.com/microcosm-cc/bluemonday"
)

type urlSlice []*url.URL

func (us *urlSlice) String() string {
	var str []string
	for _, u := range *us {
		str = append(str, u.String())
	}
	return strings.Join(str, ", ")
}

func (us *urlSlice) Set(val string) error {
	u, err := url.Parse(val)
	if err != nil {
		return err
	}
	*us = append(*us, u)
	return nil
}

type Article struct {
	Date        time.Time
	Link        string
	SourceLink  string
	SourceTitle string
	Summary     template.HTML
	Title       string
}

const (
	narticles = 3
	perSource = 1
	summaryLen = 256
	urlsFile = "openringURLs.txt"
)

func openring() (output string) {
	input, err := os.ReadFile("templates/openring.html")
	if err != nil {
		panic(err)
	}
	var sources []*url.URL
	file, err := os.Open(urlsFile)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		(*urlSlice)(&sources).Set(sc.Text())
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	file.Close()

	tmpl, err := template.
		New("template").
		Funcs(map[string]interface{}{
			"date": func(t time.Time) string {
				return t.Format("January 2, 2006")
			},
			"datef": func(fmt string, t time.Time) string {
				return t.Format(fmt)
			},
		}).
		Parse(string(input))
	if err != nil {
		panic(err)
	}

	log.Println("Fetching feeds...")
	var feeds []*rss.Feed
	for _, source := range sources {
		feed, err := rss.Fetch(source.String())
		if err != nil {
			log.Printf("Error fetching %s: %s", source.String(), err.Error())
			continue
		}
		if feed.Title == "" {
			log.Printf("Warning: feed from %s has no title", source.Host)
			feed.Title = source.Host
		}
		feeds = append(feeds, feed)
		log.Printf("Fetched %s", feed.Title)
	}
	if len(feeds) == 0 {
		log.Fatal("Expected at least one feed to successfully fetch")
	}

	policy := bluemonday.StrictPolicy()

	var articles []*Article
	for _, feed := range feeds {
		if len(feed.Items) == 0 {
			log.Printf("Warning: feed %s has no items", feed.Title)
			continue
		}
		items := feed.Items
		if len(items) > perSource {
			items = items[:perSource]
		}
		base, err := url.Parse(feed.UpdateURL)
		if err != nil {
			log.Fatalf("%s: failed parsing update URL of the feed",
				feed.UpdateURL)
		}
		feedLink, err := url.Parse(feed.Link)
		if err != nil {
			log.Fatal("failed parsing canonical feed URL of the feed")
		}
		for _, item := range items {
			raw_summary := item.Summary
			if len(raw_summary) == 0 {
				raw_summary = html.UnescapeString(item.Content)
			}
			summary := runewidth.Truncate(
				policy.Sanitize(raw_summary), summaryLen, "…")

			itemLink, err := url.Parse(item.Link)
			if err != nil {
				log.Fatalf("%s: failed parsing article URL of the feed item",
					item.Link)
			}

			articles = append(articles, &Article{
				Date:        item.Date,
				SourceLink:  base.ResolveReference(feedLink).String(),
				SourceTitle: feed.Title,
				Summary:     template.HTML(summary),
				Title:       item.Title,
				Link:        base.ResolveReference(itemLink).String(),
			})
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})
	if len(articles) < narticles {
		articles = articles[:len(articles)]
	} else {
		articles = articles[:narticles]
	}
	var b strings.Builder
	err = tmpl.Execute(&b, struct {
		Articles []*Article
	}{
		Articles: articles,
	})
	if err != nil {
		panic(err)
	}
	return b.String()
}
