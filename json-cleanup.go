package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Generates incremental ID in the format {genre}{number}, e.g., "tech001"
func generateIDs(resources *Resources) {
	// Track the genre counters
	genreCounters := make(map[string]int)

	for i, resource := range resources.List {
		genre := strings.ToLower(resource.Genre)
		genreCounters[genre]++ // Increment the counter for the genre

		// Format the ID as {genre}{001, 002, 003, ...}
		newID := fmt.Sprintf("%s%03d", genre, genreCounters[genre])
		resources.List[i].ID = newID
	}
}

// Read the JSON file
func updateIds() {
	file, err := os.Open("resources.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Read file contents
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Unmarshal JSON data into Resources struct
	var resources Resources
	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Generate IDs for resources
	generateIDs(&resources)

	// Marshal updated resources back to JSON
	updatedData, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	// Write the updated JSON to a file
	err = ioutil.WriteFile("resources_updated.json", updatedData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// Print confirmation message
	fmt.Println("Updated resources.json successfully written to resources_updated.json")
}
