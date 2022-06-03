package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

type pathData struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

type unmarshalFunc func(bytes []byte, out interface{}) error

func MakeInputHandler(f unmarshalFunc) func(inputData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return func(d []byte, fb http.Handler) (http.HandlerFunc, error) {
		pathSlice, err := makeUnmarshaler(f)(d)

		if err != nil {
			return nil, err
		}

		pathMap := sliceToMap(pathSlice)

		fmt.Printf("%v\n", pathMap)

		return MapHandler(pathMap, fb), nil
	}
}

func makeUnmarshaler(f unmarshalFunc) func(d []byte) ([]pathData, error) {
	return func(byteData []byte) ([]pathData, error) {
		var t []pathData

		err := f(byteData, &t)

		if err != nil {
			return nil, err
		}

		return t, nil
	}

}

func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := paths[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func sliceToMap(urlSlice []pathData) map[string]string {
	mapData := make(map[string]string)
	for _, urlData := range urlSlice {
		mapData[urlData.Path] = urlData.URL
	}

	return mapData
}

var YAMLHandler = MakeInputHandler(yaml.Unmarshal)
var JSONHandler = MakeInputHandler(json.Unmarshal)
