package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type handle struct {
	stories Stories
}

// func (h handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	path := r.URL.Path

// 	f story, ok := h.stories[path]; ok {

// 	}
// }

func makeStoryTemplate(story Story) {
	const tpl = `
	<!DOCTYPE html>
	
	`
}

func main() {
	storyJSON, err := getJSON()

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	jsonString, err := json.Marshal(storyJSON)

	if err != nil {

		log.Fatalf("error: %v", err)
	}

	serve()
}

// Routes:
// /: serves intro arc.  Read options from json and construct links to "/<arc-name>"

func makeHandler(stories Stories) http.Handler {

}

type StoryOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Story struct {
	Title   string        `json:"title"`
	Story   []string      `json:"story"`
	Options []StoryOption `json:"options"`
}

type Stories map[string]Story

func (stories *Stories) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, v := range m {
		var story Story

		if err := json.Unmarshal([]byte(v), &story); err != nil {
			return err
		}

		(*stories)[k] = story
	}

	return nil
}

func getJSON() (Stories, error) {
	jsonBytes, err := os.ReadFile("story.json")

	if err != nil {
		return nil, err
	}

	jsonStories := make(Stories)

	if err := json.Unmarshal(jsonBytes, &jsonStories); err != nil {
		return nil, err
	}

	return jsonStories, nil
}
