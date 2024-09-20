package main

// THIS IS FOR SCRAPING AI/ML ARTICLES:
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Article represents the structure of each ML/AI article
type Article struct {
	Title    string   `json:"title"`
	Link     string   `json:"link"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Author   string   `json:"author"`
	Type     string   `json:"type"`
}

// Define a struct to match the resources.json format
type ResourceFileArticle struct {
	Resources []Article `json:"resources"`
}

// Function to check if the article already exists in the JSON
func articleExists(articles []Article, title string) bool {
	for _, article := range articles {
		if article.Title == title {
			return true
		}
	}
	return false
}

// Function to scrape ML/AI articles and return the list
func scrapeArticles() ([]Article, error) {
	// URL to scrape
	url := "https://github.com/dair-ai/ML-Papers-of-the-Week"
	// Fetch page content
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var articles []Article

	// Find the relevant article elements
	doc.Find("article markdown-accessiblity-table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 2 {
				return
			}

			// Extract title and link
			title := cells.First().Text()
			title = strings.Split(title, " - ")[0]
			link, exists := cells.Last().Find("a[href]").Attr("href")
			if !exists {
				return
			}

			// Create an Article object and append it to the list
			article := Article{
				Title:    strings.TrimSpace(title),
				Link:     strings.TrimSpace(link),
				Category: "AI ML",
				Tags:     []string{"AI", "ML"},
				Author:   "...",
			}
			articles = append(articles, article)
		})
	})

	return articles, nil
}

// Function to save articles to resources.json
func saveArticles(articles []Article) error {
	filePath := "resources.json"

	// Read existing articles from JSON file
	var resourceFile ResourceFileArticle
	if _, err := os.Stat(filePath); err == nil {
		// File exists, so load the existing data
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(file, &resourceFile); err != nil {
			return err
		}
	}

	// Filter out duplicate articles
	var newArticles []Article
	for _, article := range articles {
		if !articleExists(resourceFile.Resources, article.Title) {
			newArticles = append(newArticles, article)
		}
	}

	// Add the new articles to the existing list
	resourceFile.Resources = append(resourceFile.Resources, newArticles...)

	// Write the updated list back to the file
	data, err := json.MarshalIndent(resourceFile, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully saved %d new articles to resources.json\n", len(newArticles))
	return nil
}

// // Scrape articles
// articles, err := scrapeArticles()
// if err != nil {
// 	fmt.Printf("Error while scraping articles: %v\n", err)
// 	return
// }

// // Save articles to resources.json
// if err := saveArticles(articles); err != nil {
// 	fmt.Printf("Error while saving articles: %v\n", err)
// 	return
// }
