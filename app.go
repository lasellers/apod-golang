// fetch NASA APOD data and display as a webpage on localhost:8080
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type apod struct {
	Title          string `json:"title"`
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	URL            string `json:"url"`
	HdURL          string `json:"hdurl"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainPage(w)
	})
	var PORT string = ":8080"
	fmt.Printf("Starting server at port %s\n", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func mainPage(w http.ResponseWriter) {
	const AppAPIKey string = "HAzNETw6VUoq4BIZ0KikSObCBP0xe0gMSbbYaubP"
	const APIUrl string = "https://api.nasa.gov/planetary/apod?api_key=" + AppAPIKey
	const TimeoutMax time.Duration = 3

	// init http client
	var apodClient = &http.Client{
		Timeout: time.Second * TimeoutMax,
	}

	// setup get request
	req, err := http.NewRequest(http.MethodGet, APIUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Intrafoundation Software")

	// do actual REST API GET fetch
	res, getErr := apodClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// if request succeeds we have to explictly close it
	if res.Body != nil {
		defer res.Body.Close()
	}

	// get data from API call
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// get body containing json string and convert to a go struct type
	datastore := apod{}
	jsonErr := json.Unmarshal(body, &datastore)
	if jsonErr != nil {
		log.Fatal(readErr)
	}

	// print html
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html lang=\"en\">")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<h1>Astronomy Photo of the Day</h1>")

	fmt.Fprintf(w, "<h2>%s</h2>\n", datastore.Title)
	fmt.Fprintf(w, "<p><b>Copyright:</b> %s</p>\n", datastore.Copyright)
	fmt.Fprintf(w, "<p><b>Date:</b> %s</p>\n", datastore.Date)
	fmt.Fprintf(w, "<p><b>Explanation:</b> %s</p>\n", datastore.Explanation)
	fmt.Fprintf(w, "<p><b>Media Type:</b> %s</p>\n", datastore.MediaType)
	fmt.Fprintf(w, "<p><b>Service Version:</b> %s</p>\n", datastore.ServiceVersion)
	fmt.Fprintf(w, "<p><b>Url:</b> %s</p>\n", datastore.URL)
	fmt.Fprintf(w, "<p><b>HD Url:</b> %s</p>\n", datastore.HdURL)
	fmt.Fprintf(w, "<img src=\"%s\">\n", datastore.URL)

	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}
