package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Assuming InputForm type is similar to UpdateInputs
type InputForm struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Link   string   `json:"link"`
	Genre  string   `json:"genre"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
	Type   string   `json:"type"`
}

// ResourceFileArticle structure to hold resources
type ResourceFileUpdateInputs struct {
	Resources []InputForm `json:"resources"`
}

// Function to check if an article exists
func InputExists(articles []InputForm, title string) bool {
	for _, a := range articles {
		if a.Title == title {
			return true
		}
	}
	return false
}

// Function to generate a unique ID based on genre
func generateID(genre string, resourceFile ResourceFileUpdateInputs) string {
	maxID := 0

	// Iterate over existing resources to find the highest ID for the given genre
	for _, article := range resourceFile.Resources {
		if strings.HasPrefix(article.ID, genre) {
			idNum, err := strconv.Atoi(strings.TrimPrefix(article.ID, genre))
			if err == nil && idNum > maxID {
				maxID = idNum
			}
		}
	}

	// Generate new ID by incrementing the highest found ID
	newID := fmt.Sprintf("%s%d", genre, maxID+1)
	return newID
}

// New function to save articles with type
func updateResourcesWithType(sheetID string) error {
	// Construct the URL to fetch CSV data for a specific range (A:F)
	csvURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&range=A:F", sheetID)

	// Fetch the CSV data
	resp, err := http.Get(csvURL)
	if err != nil {
		return fmt.Errorf("error fetching CSV: %v", err)
	}
	defer resp.Body.Close()
	// Parse the CSV
	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true       // Allow lazy quotes
	reader.TrimLeadingSpace = true // Trim leading space

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	var newInputForms []InputForm

	// Process the CSV data
	for i, row := range rows {
		if i == 0 {
			continue // Skip the header row
		}
		if len(row) < 6 {
			continue // Ensure there are enough columns
		}

		article := InputForm{
			ID:     generateID(strings.TrimSpace(row[3]), ResourceFileUpdateInputs{}),
			Title:  strings.TrimSpace(row[0]),
			Author: strings.TrimSpace(row[1]),
			Link:   strings.TrimSpace(row[2]),
			Genre:  strings.TrimSpace(row[3]),
			Status: "unread", // Default status
			Tags:   strings.Split(strings.TrimSpace(row[4]), ","),
			Type:   strings.TrimSpace(row[5]),
		}

		if !InputExists(newInputForms, article.Title) {
			newInputForms = append(newInputForms, article)
		}
	}
	// Save new articles to resources.json
	return saveInputFormsWithType(newInputForms)
}

func saveInputFormsWithType(articles []InputForm) error {
	filePath := "resources.json"

	// Read existing articles from JSON file
	var resourceFile ResourceFileUpdateInputs
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

	// Add the new articles to the existing list
	resourceFile.Resources = append(resourceFile.Resources, articles...)

	// Write the updated list back to the file
	data, err := json.MarshalIndent(resourceFile, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully saved %d new articles to resources.json\n", len(articles))
	return nil
}
