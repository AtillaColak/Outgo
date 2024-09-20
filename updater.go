package main

import (
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

	// Filter out duplicate articles and ensure all fields are filled
	var newInputForms []InputForm
	for _, article := range articles {
		if !InputExists(resourceFile.Resources, article.Title) {
			newInputForms = append(newInputForms, InputForm(article)) // Convert UpdateInputs to InputForm
		}
	}

	// Add the new articles to the existing list
	resourceFile.Resources = append(resourceFile.Resources, newInputForms...)

	// Write the updated list back to the file
	data, err := json.MarshalIndent(resourceFile, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully saved %d new articles to resources.json\n", len(newInputForms))
	return nil
}

// Function to process and save articles from the provided input
func updateResourcesWithType(sheetURL string) error {
	// Fetch the data from the public Google Sheets URL
	resp, err := http.Get(sheetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the TSV format directly (assuming Google Sheets provides it as TSV)
	lines := strings.Split(string(body), "\n")
	var newInputForms []InputForm

	// Read existing articles for ID generation
	var resourceFile ResourceFileUpdateInputs
	if _, err := os.Stat("resources.json"); err == nil {
		file, err := ioutil.ReadFile("resources.json")
		if err == nil {
			json.Unmarshal(file, &resourceFile)
		}
	}

	// Process the TSV data
	for i, line := range lines {
		if i == 0 {
			continue // Skip the header
		}
		columns := strings.Split(line, "\t") // Split by tab character
		if len(columns) < 6 {
			continue // Ensure there are enough columns
		}

		article := InputForm{
			ID:     generateID(strings.TrimSpace(columns[3]), resourceFile), // Generate ID based on genre
			Title:  strings.TrimSpace(columns[0]),
			Author: strings.TrimSpace(columns[1]),
			Link:   strings.TrimSpace(columns[2]),
			Genre:  strings.TrimSpace(columns[3]),
			Status: "unread", // Default status; adjust as needed
			Tags:   strings.Split(strings.TrimSpace(columns[4]), ","),
			Type:   strings.TrimSpace(columns[5]), // New field for type
		}

		if !InputExists(newInputForms, article.Title) {
			newInputForms = append(newInputForms, article)
		}
	}

	// Save new articles to resources.json
	return saveInputFormsWithType(newInputForms)
}
