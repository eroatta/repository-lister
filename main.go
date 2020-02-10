package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	repositories := make(chan string)
	go func() {
		for i := 0; i < 50; i++ {
			repositories <- "."
		}
		close(repositories)
	}()

	for n := range repositories {
		log.Println(n)
	}
}

func main2() {
	page := 1
	size := 50
	more := true
	for more == true {
		url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:>=1000+language:go&sort=stars&order=desc&per_page=%d&page=%d", size, page)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			log.Fatal("Unexpected status code", response.StatusCode)
		}

		data := JSONData{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}
		printData(data)

		if data.Count > (page * size) {
			log.Println("Requesting new page...")
		} else {
			more = false
		}
		page++
	}

}

func printData(data JSONData) {
	log.Printf("Repositories found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Repository", "Stars", "Created at", "Description")
	fmt.Fprintf(tw, format, "----------", "-----", "----------", "----------")
	for _, i := range data.Items {
		desc := i.Description
		if len(desc) > 50 {
			desc = string(desc[:50]) + "..."
		}
		t, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tw, format, i.FullName, i.StargazersCount, t.Year(), desc)
	}
	tw.Flush()
}

// Owner is the repository owner
type Owner struct {
	Login string
}

// Item is the single repository data structure
type Item struct {
	ID              int
	Name            string
	FullName        string `json:"full_name"`
	Owner           Owner
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

// JSONData contains the GitHub API response
type JSONData struct {
	Count int `json:"total_count"`
	Items []Item
}
