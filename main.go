package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Fact struct {
	UUID string `json:"id"`
	URL string `json:"url"`
	Content string `json:"fact"`
}

const (
	FactsDir string = "facts"
	Port string = ":80"
)

func main() {
	log.Printf("Now Listening on %s", Port)
	
	http.HandleFunc("/", serveFiles)

	http.HandleFunc("/random", randomFactAsJSON)
	http.HandleFunc("/random.json", randomFactAsJSON)
	http.HandleFunc("/random.md", randomFactAsMD)
	http.HandleFunc("/random.txt", randomFactAsTXT)

	log.Fatal(http.ListenAndServe(Port, nil))
}

func randomElement(array []os.FileInfo) os.FileInfo {
	return array[rand.Intn(len(array))]
}

func randomFactAsMD(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(FactsDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/markdown")
	w.WriteHeader(http.StatusCreated)

	fact, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", FactsDir, randomElement(files).Name()))
	if err != nil {
		log.Fatal(err)
	}

	var factStruct Fact
	json.Unmarshal([]byte(string(fact)), &factStruct)

	fmt.Fprintf(w, "> %s\nPermalink: %s", factStruct.Content, factStruct.URL)
}

func randomFactAsTXT(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(FactsDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	fact, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", FactsDir, randomElement(files).Name()))
	if err != nil {
		log.Fatal(err)
	}

	var factStruct Fact
	json.Unmarshal([]byte(string(fact)), &factStruct)

	fmt.Fprintf(w, factStruct.Content)
}

func randomFactAsJSON(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(FactsDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	fact, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", FactsDir, randomElement(files).Name()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(fact))
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	p := "." + r.URL.Path

	if p == "./" {
		http.Redirect(w, r, "https://docs.captainbot.xyz/services/fact-api", 301)
		return	
	} else {
		p = fmt.Sprintf("./%s%s", FactsDir, r.URL.Path)
	}
	
	http.ServeFile(w, r, p)
}