package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid" // for generating unique IDs
	"github.com/olekukonko/tablewriter"
)

// PROCESS: Scraped top websites, populate the resources array, add update button to automatically fetch from the google sheets. add filtering options. add pagination to lists. Add Menu items and improve the CLI UI.
type Resource struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Type   string   `json:"type"`
	Genre  string   `json:"genre"`
	Status string   `json:"status"`
	Link   string   `json:"link"`
	Tags   []string `json:"tags"`
	Author string   `json:"author,omitempty"`
}

type Resources struct {
	List []Resource `json:"resources"`
}

type Playlist struct {
	ID        string     `json:id`
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

type Playlists struct {
	List []Playlist `json:"playlists"`
}

const (
	resourcesFile = "resources.json"
	playlistsFile = "playlists.json"
)

const itemsPerPage = 20

// Colors for genres, statuses, and tags
var genreColors = map[string]color.Attribute{
	"self-improvement": color.FgGreen,
	"AI ML":            color.FgCyan,
	"history":          color.FgYellow,
	"finance":          color.FgBlue,
	"tech":             color.FgRed,
	// TODO: I'll add more tags
}

var statusColors = map[string]color.Attribute{
	"unread":      color.FgWhite,
	"viewed":      color.FgMagenta,
	"in-progress": color.FgYellow,
	"not-started": color.FgRed,
}

var tagColors = map[string]color.Attribute{
	"self-improvement": color.FgGreen,
	"AI ML":            color.FgCyan,
	"history":          color.FgYellow,
	"finance":          color.FgBlue,
	"Tech":             color.FgRed,
	// TODO: add more tags
}

var resourceFields = map[string]bool{
	"ID":     true,
	"Title":  true,
	"Genre":  true,
	"Type":   false,
	"Status": false,
	"Tags":   false,
}

var playlistFields = map[string]bool{
	"Name":      true,
	"Resources": true,
	"ID":        true,
}

func loadResources() (Resources, error) {
	var resources Resources
	file, err := ioutil.ReadFile(resourcesFile)
	if err != nil {
		return resources, err
	}
	err = json.Unmarshal(file, &resources)
	return resources, err
}

func saveResources(resources Resources) error {
	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(resourcesFile, data, 0644)
}

func loadPlaylists() (Playlists, error) {
	var playlists Playlists
	file, err := ioutil.ReadFile(playlistsFile)
	if err != nil {
		return playlists, err
	}
	err = json.Unmarshal(file, &playlists)
	return playlists, err
}

func savePlaylists(playlists Playlists) error {
	data, err := json.MarshalIndent(playlists, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(playlistsFile, data, 0644)
}

func addResource(reader *bufio.Reader) {
	resource := Resource{}

	fmt.Print("Enter ID (e.g., prog001): ")
	resource.ID, _ = reader.ReadString('\n')
	resource.ID = strings.TrimSpace(resource.ID)

	fmt.Print("Enter title: ")
	resource.Title, _ = reader.ReadString('\n')
	resource.Title = strings.TrimSpace(resource.Title)

	fmt.Print("Enter type (book/video/podcast/website/course): ")
	resource.Type, _ = reader.ReadString('\n')
	resource.Type = strings.TrimSpace(resource.Type)

	fmt.Print("Enter genre: ")
	resource.Genre, _ = reader.ReadString('\n')
	resource.Genre = strings.TrimSpace(resource.Genre)

	fmt.Print("Enter status (unread/viewed/in-progress/not-started): ")
	resource.Status, _ = reader.ReadString('\n')
	resource.Status = strings.TrimSpace(resource.Status)

	fmt.Print("Enter link: ")
	resource.Link, _ = reader.ReadString('\n')
	resource.Link = strings.TrimSpace(resource.Link)

	fmt.Print("Enter tags (comma-separated): ")
	tags, _ := reader.ReadString('\n')
	tags = strings.TrimSpace(tags)
	resource.Tags = strings.Split(tags, ",")

	if resource.Type == "book" {
		fmt.Print("Enter author: ")
		resource.Author, _ = reader.ReadString('\n')
		resource.Author = strings.TrimSpace(resource.Author)
	}

	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	resources.List = append(resources.List, resource)
	err = saveResources(resources)
	if err != nil {
		color.Red("Error saving resources: %v", err)
	} else {
		color.Green("Resource added successfully!")
	}
}

func deleteResource(reader *bufio.Reader) {
	fmt.Print("Enter ID of resource to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	for i, resource := range resources.List {
		if strings.EqualFold(resource.ID, id) {
			resources.List = append(resources.List[:i], resources.List[i+1:]...)
			err := saveResources(resources)
			if err != nil {
				color.Red("Error saving resources: %v", err)
			} else {
				color.Green("Deleted resource: %s", id)
			}
			return
		}
	}
	color.Yellow("Resource with ID %s not found.", id)
}

func filterResources(reader *bufio.Reader) {
	fmt.Print("Enter filter criteria (genre/tag/status): ")
	criteria, _ := reader.ReadString('\n')
	criteria = strings.TrimSpace(criteria)

	fmt.Printf("Enter %s to filter by: ", criteria)
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	filtered := Resources{}
	for _, r := range resources.List {
		switch criteria {
		case "genre":
			if strings.EqualFold(r.Genre, value) {
				filtered.List = append(filtered.List, r)
			}
		case "tag":
			for _, tag := range r.Tags {
				if strings.EqualFold(tag, value) {
					filtered.List = append(filtered.List, r)
					break
				}
			}
		case "status":
			if strings.EqualFold(r.Status, value) {
				filtered.List = append(filtered.List, r)
			}
		default:
			color.Red("Unknown filter criteria: %s", criteria)
			return
		}
	}

	if len(filtered.List) == 0 {
		color.Yellow("No resources found for %s: %s", criteria, value)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Type", "Genre", "Status"})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})

	for _, r := range filtered.List {
		table.Append([]string{r.ID, r.Title, r.Type, r.Genre, r.Status})
	}

	table.Render()
}

func markResourceStatus(reader *bufio.Reader) {
	fmt.Print("Enter ID of resource to mark: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Enter new status (unread/viewed/in progress/not started): ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	for i, resource := range resources.List {
		if strings.EqualFold(resource.ID, id) {
			resources.List[i].Status = status
			err := saveResources(resources)
			if err != nil {
				color.Red("Error saving resources: %v", err)
			} else {
				color.Green("Updated status of resource: %s to %s", id, status)
			}
			return
		}
	}
	color.Yellow("Resource with ID %s not found.", id)
}

func createPlaylist(reader *bufio.Reader) {
	fmt.Print("Enter playlist name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Generate a unique ID for the playlist
	playlistID := uuid.New().String()

	playlists, err := loadPlaylists()
	if err != nil {
		color.Red("Error loading playlists: %v", err)
		return
	}

	playlist := Playlist{ID: playlistID, Name: name}
	for {
		fmt.Print("Enter resource ID to add to playlist (or 'done' to finish): ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)

		if strings.EqualFold(id, "done") {
			break
		}

		resources, err := loadResources()
		if err != nil {
			color.Red("Error loading resources: %v", err)
			return
		}

		for _, resource := range resources.List {
			if strings.EqualFold(resource.ID, id) {
				playlist.Resources = append(playlist.Resources, resource)
				break
			}
		}
	}

	playlists.List = append(playlists.List, playlist)
	err = savePlaylists(playlists)
	if err != nil {
		color.Red("Error saving playlists: %v", err)
	} else {
		color.Green("Playlist '%s' created successfully with ID: %s !", name, playlistID)
	}
}

func addResourceToPlaylist(reader *bufio.Reader) {
	fmt.Print("Enter playlist name: ")
	playlistName, _ := reader.ReadString('\n')
	playlistName = strings.TrimSpace(playlistName)

	fmt.Print("Enter resource ID to add to playlist: ")
	resourceID, _ := reader.ReadString('\n')
	resourceID = strings.TrimSpace(resourceID)

	playlists, err := loadPlaylists()
	if err != nil {
		color.Red("Error loading playlists: %v", err)
		return
	}

	for i, playlist := range playlists.List {
		if strings.EqualFold(playlist.Name, playlistName) {
			resources, err := loadResources()
			if err != nil {
				color.Red("Error loading resources: %v", err)
				return
			}

			for _, resource := range resources.List {
				if strings.EqualFold(resource.ID, resourceID) {
					playlists.List[i].Resources = append(playlists.List[i].Resources, resource)
					err := savePlaylists(playlists)
					if err != nil {
						color.Red("Error saving playlists: %v", err)
					} else {
						color.Green("Added resource %s to playlist '%s'", resourceID, playlistName)
					}
					return
				}
			}
			color.Yellow("Resource with ID %s not found.", resourceID)
			return
		}
	}
	color.Yellow("Playlist with name %s not found.", playlistName)
}

func removeResourceFromPlaylist(reader *bufio.Reader) {
	fmt.Print("Enter playlist name: ")
	playlistName, _ := reader.ReadString('\n')
	playlistName = strings.TrimSpace(playlistName)

	fmt.Print("Enter resource ID to remove from playlist: ")
	resourceID, _ := reader.ReadString('\n')
	resourceID = strings.TrimSpace(resourceID)

	playlists, err := loadPlaylists()
	if err != nil {
		color.Red("Error loading playlists: %v", err)
		return
	}

	for i, playlist := range playlists.List {
		if strings.EqualFold(playlist.Name, playlistName) {
			for j, r := range playlist.Resources {
				if strings.EqualFold(r.ID, resourceID) {
					playlists.List[i].Resources = append(playlists.List[i].Resources[:j], playlists.List[i].Resources[j+1:]...)
					err := savePlaylists(playlists)
					if err != nil {
						color.Red("Error saving playlists: %v", err)
					} else {
						color.Green("Removed resource %s from playlist '%s'", resourceID, playlistName)
					}
					return
				}
			}
			color.Yellow("Resource with ID %s not found in playlist %s.", resourceID, playlistName)
			return
		}
	}
	color.Yellow("Playlist with name %s not found.", playlistName)
}

func toggleField(fields map[string]bool, field string) {
	if _, exists := fields[field]; exists {
		fields[field] = !fields[field]
	}
}

func showFieldOptions(fields map[string]bool) {
	fmt.Println("Field Options:")
	for field, visible := range fields {
		var indicator string
		var colorFunc func(a ...interface{}) string
		if visible {
			indicator = "✔️"
			colorFunc = color.New(color.FgGreen).SprintFunc()
		} else {
			indicator = "❌"
			colorFunc = color.New(color.FgRed).SprintFunc()
		}
		fmt.Printf("%s %s\n", colorFunc(indicator), field)
	}
	fmt.Println("Press 'b' to go back.")
}

// Function to list resources with pagination
func listResources() {
	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	if len(resources.List) == 0 {
		color.Yellow("No resources found.")
		return
	}

	totalPages := (len(resources.List) + itemsPerPage - 1) / itemsPerPage
	currentPage := 0

	for {
		start := currentPage * itemsPerPage
		end := start + itemsPerPage
		if end > len(resources.List) {
			end = len(resources.List)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetRowLine(true) // Adds a line between each row

		var headers []string
		var columnColors []tablewriter.Colors

		// Prepare headers and column colors
		for field, visible := range resourceFields {
			if visible {
				headers = append(headers, field)
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor}) // Default color
			}
		}
		table.SetHeader(headers)
		table.SetColumnColor(columnColors...)

		for _, r := range resources.List[start:end] {
			// Ensure headers and columns are aligned
			var row []string
			if resourceFields["ID"] {
				row = append(row, r.ID)
			}
			if resourceFields["Title"] {
				row = append(row, r.Title)
			}
			if resourceFields["Genre"] {
				genreColor := genreColors[r.Genre]
				row = append(row, color.New(genreColor).Sprintf(r.Genre))
			}
			if resourceFields["Type"] {
				row = append(row, r.Type)
			}
			if resourceFields["Status"] {
				statusColor := statusColors[r.Status]
				row = append(row, color.New(statusColor).Sprintf(r.Status))
			}
			if resourceFields["Tags"] {
				var tagStrings []string
				for _, tag := range r.Tags {
					tagColor := tagColors[tag]
					tagStrings = append(tagStrings, color.New(tagColor).Sprintf(tag))
				}
				row = append(row, strings.Join(tagStrings, ", "))
			}

			// Check if row length matches header length
			if len(row) != len(headers) {
				fmt.Println("Row length doesn't match header length!")
				continue
			}

			table.Append(row)
		}

		fmt.Print("\033[38;5;201m") // Set text color to light magenta (pink tone)
		table.Render()
		fmt.Print("\033[0m") // Reset text color here

		fmt.Printf("Page %d of %d\n", currentPage+1, totalPages)
		fmt.Println("Options: [1] Go Left, [2] Go Right, [3] Return to Previous Screen")

		choice := getUserChoice()
		switch choice {
		case "1":
			if currentPage > 0 {
				currentPage--
			}
		case "2":
			if currentPage < totalPages-1 {
				currentPage++
			}
		case "3":
			return
		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}

// Function to list playlists with pagination
func listPlaylists() {
	playlists, err := loadPlaylists()
	if err != nil {
		color.Red("Error loading playlists: %v", err)
		return
	}

	if len(playlists.List) == 0 {
		color.Yellow("No playlists found.")
		return
	}

	totalPages := (len(playlists.List) + itemsPerPage - 1) / itemsPerPage
	currentPage := 0

	for {
		start := currentPage * itemsPerPage
		end := start + itemsPerPage
		if end > len(playlists.List) {
			end = len(playlists.List)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetRowLine(true)

		table.SetHeader([]string{"ID", "Name"})
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.FgWhiteColor},
		)

		for _, p := range playlists.List[start:end] {
			table.Append([]string{p.ID, p.Name})
		}

		fmt.Print("\033[32m") // Set text color to green
		table.Render()
		fmt.Print("\033[0m") // Reset text color

		fmt.Printf("Page %d of %d\n", currentPage+1, totalPages)
		fmt.Println("Options: [1] Go Left, [2] Go Right, [3] Return to Previous Screen")

		choice := getUserChoice()
		switch choice {
		case "1":
			if currentPage > 0 {
				currentPage--
			}
		case "2":
			if currentPage < totalPages-1 {
				currentPage++
			}
		case "3":
			return
		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}

// Function to get a random resource
func getRandomResource() {
	resources, err := loadResources()
	if err != nil {
		color.Red("Error loading resources: %v", err)
		return
	}

	if len(resources.List) == 0 {
		color.Yellow("No resources found.")
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(resources.List))
	randomResource := resources.List[randomIndex]

	color.Green("Random Resource: %s (ID: %s)", randomResource.Title, randomResource.ID)
}

// Function to get user choice
func getUserChoice() string {
	var choice string
	fmt.Scan(&choice)
	return choice
}

func fieldOptions(reader *bufio.Reader, fields map[string]bool) {
	for {
		showFieldOptions(fields)
		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "b" {
			break
		}
		toggleField(fields, choice)
	}
}

func viewPlaylistByID(reader *bufio.Reader) {
	fmt.Print("Enter playlist ID: ")
	playlistID, _ := reader.ReadString('\n')
	playlistID = strings.TrimSpace(playlistID)

	playlists, err := loadPlaylists()
	if err != nil {
		color.Red("Error loading playlists: %v", err)
		return
	}

	for _, playlist := range playlists.List {
		if strings.EqualFold(playlist.ID, playlistID) {
			// Display playlist name and ID
			color.Cyan("Playlist: %s\n", playlist.Name)
			color.Yellow("ID: %s\n", playlist.ID)

			// Render resources in the playlist
			if len(playlist.Resources) == 0 {
				color.Red("No resources in this playlist.")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAutoFormatHeaders(false)
			table.SetRowLine(true)

			// Prepare headers and column colors based on the visibility in resourceFields
			var headers []string
			var columnColors []tablewriter.Colors

			if resourceFields["ID"] {
				headers = append(headers, "Resource ID")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}
			if resourceFields["Title"] {
				headers = append(headers, "Title")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}
			if resourceFields["Genre"] {
				headers = append(headers, "Genre")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}
			if resourceFields["Type"] {
				headers = append(headers, "Type")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}
			if resourceFields["Status"] {
				headers = append(headers, "Status")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}
			if resourceFields["Tags"] {
				headers = append(headers, "Tags")
				columnColors = append(columnColors, tablewriter.Colors{tablewriter.FgWhiteColor})
			}

			// Set the headers and column colors in the table
			table.SetHeader(headers)
			table.SetColumnColor(columnColors...)

			// Loop through resources and color fields like Genre, Status, and Tags based on resourceFields visibility
			for _, resource := range playlist.Resources {
				var row []string

				if resourceFields["ID"] {
					row = append(row, resource.ID)
				}
				if resourceFields["Title"] {
					row = append(row, resource.Title)
				}
				if resourceFields["Genre"] {
					genreColor := genreColors[resource.Genre]
					coloredGenre := color.New(genreColor).Sprintf(resource.Genre)
					row = append(row, coloredGenre)
				}
				if resourceFields["Type"] {
					row = append(row, resource.Type)
				}
				if resourceFields["Status"] {
					statusColor := statusColors[resource.Status]
					coloredStatus := color.New(statusColor).Sprintf(resource.Status)
					row = append(row, coloredStatus)
				}
				if resourceFields["Tags"] {
					var tagStrings []string
					for _, tag := range resource.Tags {
						tagColor := tagColors[tag]
						tagStrings = append(tagStrings, color.New(tagColor).Sprintf(tag))
					}
					coloredTags := strings.Join(tagStrings, ", ")
					row = append(row, coloredTags)
				}

				// Append the row to the table
				table.Append(row)
			}

			// Render the table with colored data
			fmt.Print("\033[38;5;117m") // Set to light blue
			table.Render()
			fmt.Print("\033[0m") // Reset text color here
			return
		}
	}
	color.Red("Playlist with ID %s not found.", playlistID)
}

func printHelp() {
	color.Cyan(`
Available Commands:
- add: Add a new resource
- list: List all resources
- delete: Delete a resource
- fetch-updates: Fetch the newest resources
- filter: Filter resources by genre or tag
- mark: Mark a resource as read/viewed/etc.
- create-playlist: Create a new playlist
- list-playlists: List all playlists
- view-playlist: Inspect a specific playlist from its id. 
- add-to-playlist: Add a resource to a playlist
- remove-from-playlist: Remove a resource from a playlist
- filter-fields: Toggle fields for listing resources
- filter-playlist-fields: Toggle fields for listing playlists
- random-resource: Get a single random resource
- help: Show this help message
- update: Fetch and add new resources from YouTube or similar sources
- exit: Exit the application`)
	color.Green(`Credits:
- Developed by Atilla Colak
- Special thanks to the Go community for inspiration and support.`)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	printHelp() // Show help on startup

	for {
		fmt.Print("\nEnter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "add":
			addResource(reader)
		case "list":
			listResources()
		case "delete":
			deleteResource(reader)
		case "fetch-updates":
			sheetURL := "1wganKHEJps87WhFI2O_xyVw-3vkTshmaf665OKczbwc"
			if err := updateResourcesWithType(sheetURL); err != nil {
				fmt.Printf("Error updating resources from Google Sheets: %v\n", err)
			}
		case "filter":
			filterResources(reader)
		case "mark":
			markResourceStatus(reader)
		case "create-playlist":
			createPlaylist(reader)
		case "list-playlists":
			listPlaylists()
		case "view-playlist":
			viewPlaylistByID(reader)
		case "add-to-playlist":
			addResourceToPlaylist(reader)
		case "remove-from-playlist":
			removeResourceFromPlaylist(reader)
		case "filter-fields":
			fieldOptions(reader, resourceFields)
		case "filter-playlist-fields":
			fieldOptions(reader, playlistFields)
		case "random-resource":
			getRandomResource()
		case "help", "?":
			printHelp()
		case "exit", "quit":
			fmt.Println("Exiting the application. Goodbye!")
			return
		default:
			color.Red("Unknown command: %s", command)
			printHelp()
		}
	}
}
