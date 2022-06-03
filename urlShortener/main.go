package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rossh87/urlShortener/database"
	"github.com/rossh87/urlShortener/handlers"
)

func main() {
	var yamlFlag string
	// var jsonFlag string

	flag.StringVar(&yamlFlag, "y", "default.yaml", "sets YAML file to use for url shortening router")
	// flag.StringVar(&jsonFlag, "j", "default.json", "sets JSON file to use for url shortening router")

	flag.Parse()

	db, err := database.SetupDB()

	if err != nil {
		err = fmt.Errorf("failed to setup db:\n%v", err)
		panic(err)
	}

	dbPaths := []database.PathData{{Path: "/go", URL: "https://www.google.com"}}

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	database.SetPathData(db, dbPaths)

	dbPathData, err := database.GetPathData(db)

	if err != nil {
		panic(err)
	}

	mapHandler := handlers.MapHandler(pathsToUrls, defaultMux())

	yamlHandler := handlerFromFlag(yamlFlag, handlers.YAMLHandler, mapHandler)

	dbHandler, err := handlers.MakeInputHandler(json.Unmarshal)(dbPathData, yamlHandler)

	if err != nil {
		panic(err)
	}
	// jsonHandler := handlerFromFlag(yamlFlag, handlers.JSONHandler, yamlHandler)

	fmt.Println("starting server on 8080...")

	http.ListenAndServe(":8080", dbHandler)
}

type inputHandler func(inputData []byte, fallback http.Handler) (http.HandlerFunc, error)

func handlerFromFlag(flagValue string, inputHandler inputHandler, fallback http.HandlerFunc) http.HandlerFunc {
	fileData, err := ioutil.ReadFile(flagValue)

	if err != nil {
		fmt.Printf("Unable to read file %s for the following reason:\n%v", flagValue, err)
		panic(err)
	}

	handler, err := inputHandler(fileData, fallback)

	if err != nil {
		fmt.Printf("Failed when parsing data from %s to map, for the following reason:\n%v", flagValue, err)
		panic(err)
	}

	return handler
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", greetResponse)

	return mux
}

func greetResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Greetings!")
}
