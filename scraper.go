package main

// THIS IS FOR SCRAPING 400 BOOKS.
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Book represents the structure of each book
type Book struct {
	Title    string   `json:"title"`
	Link     string   `json:"link"`
	Author   string   `json:"author"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// Function to extract the genre from the URL
func extractGenreFromURL(url string) string {
	// Example URL: https://www.shortform.com/best-books/genre/best-history-books-of-all-time
	// We want to extract "history" from this URL

	// Split the URL into parts using "/"
	parts := strings.Split(url, "/")

	// Check if the URL has the expected segments
	if len(parts) >= 4 {
		// Extract the segment containing the genre info
		genrePart := parts[len(parts)-1] // Example: "best-history-books-of-all-time"

		// Remove "best-" and "-books-of-all-time" from the segment
		genrePart = strings.TrimPrefix(genrePart, "best-")
		genrePart = strings.TrimSuffix(genrePart, "-books-of-all-time")

		return genrePart
	}

	return "unknown"
}

// Function to check if the book already exists in the JSON and if it has the same genre
func bookExists(books []Book, title string) (*Book, bool) {
	for _, book := range books {
		if book.Title == title {
			return &book, true
		}
	}
	return nil, false
}

// Function to scrape books from a single URL and return the list
func scrapeBooks(urls []string) ([]Book, error) {
	var allBooks []Book

	for _, url := range urls {
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

		// Find the relevant book elements
		doc.Find("div.card.border").Each(func(i int, s *goquery.Selection) {
			// Get the book title
			title := s.Find("h2.display-4").Text()

			// Get the Amazon buy link
			link, exists := s.Find("a[rel=nofollow]").Attr("href")
			if !exists {
				link = ""
			}

			// Get the author
			author := s.Find("p.byline span").First().Text()

			// Get the genre from the URL
			category := extractGenreFromURL(url)

			// Create a Book object
			book := Book{
				Title:    strings.TrimSpace(title),
				Link:     strings.TrimSpace(link),
				Author:   strings.TrimSpace(author),
				Category: category,
				Tags:     []string{category}, // Start with the genre as a tag
			}

			// Check if the book already exists in the list and update its tags
			existingBook, exists := bookExists(allBooks, book.Title)
			if exists {
				existingBook.Tags = append(existingBook.Tags, category)
			} else {
				allBooks = append(allBooks, book)
			}
		})
	}

	return allBooks, nil
}

// Define a struct to match the resources.json format
type ResourceFile struct {
	Resources []Book `json:"resources"`
}

// Function to save books to resources.json
func saveBooks(books []Book) error {
	filePath := "resources.json"

	// Read existing books from JSON file
	var resourceFile ResourceFile
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

	// Filter out duplicate books and update tags
	for _, book := range books {
		if existingBook, exists := bookExists(resourceFile.Resources, book.Title); exists {
			// Update tags if the book already exists
			existingBook.Tags = append(existingBook.Tags, book.Tags...)
		} else {
			// Add the new book
			resourceFile.Resources = append(resourceFile.Resources, book)
		}
	}

	// Write the updated list back to the file
	data, err := json.MarshalIndent(resourceFile, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully saved %d books to resources.json\n", len(books))
	return nil
}

// BELOW IS THE MAIN CODE I USED TO EXECUTE ALL.
// List of URLs to scrape
// urls := []string{
// 	"https://www.shortform.com/best-books/genre/best-tech-books-of-all-time",
// 	"https://www.shortform.com/best-books/genre/best-finance-books-of-all-time",
// 	"https://www.shortform.com/best-books/genre/best-business-books-of-all-time",
// 	"https://www.shortform.com/best-books/genre/best-self-improvement-books-of-all-time",
// 	"https://www.shortform.com/best-books/genre/best-history-books-of-all-time",
// }

// // Scrape books from the list of URLs
// books, err := scrapeBooks(urls)
// if err != nil {
// 	fmt.Printf("Error while scraping books: %v\n", err)
// 	return
// }

// // Save books to resources.json
// err = saveBooks(books)
// if err != nil {
// 	fmt.Printf("Error while saving books: %v\n", err)
// 	return
// }
