package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

const size = 50

func main() {
	token := flag.String("token", "", "GitHub Access Token")
	flag.Parse()

	if *token == "" {
		log.Fatal("invalid token")
	}

	responses := make(chan GitHubResponse)
	go func() {
		page := 1
		more := true
		for more {
			url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:>=1000+language:go&sort=stars&order=desc&per_page=%d&page=%d", size, page)
			request, _ := http.NewRequest("GET", url, nil)
			request.Header.Add("Authorization", fmt.Sprintf("token %s", *token))

			client := http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			// check for errors on the response
			if response.StatusCode != http.StatusOK {
				ghError := GitHubErrorResponse{}
				_ = json.Unmarshal(body, &ghError)
				log.Fatal(fmt.Sprintf("Unexpected status code: %v", ghError), response.StatusCode)
			}

			// extract the results and send them to the chan
			var ghResponse GitHubResponse
			err = json.Unmarshal(body, &ghResponse)
			if err != nil {
				log.Fatal(err)
			}
			responses <- ghResponse

			// check if we should keep looking for
			if ghResponse.Count > (page*size) && (page+1)*size <= 1000 {
				page++
			} else {
				more = false
			}
		}
		close(responses)
	}()

	// print out the incoming results
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Repository", "Stars", "Created at", "Description")
	fmt.Fprintf(tw, format, "----------", "-----", "----------", "-----------")
	for response := range responses {
		for _, i := range response.Items {
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
	}
	tw.Flush()
}

// GitHubResponse contains the GitHub API response.
type GitHubResponse struct {
	Count int `json:"total_count"`
	Items []Item
}

// Item is the single repository data structure
type Item struct {
	ID              int
	Name            string
	FullName        string `json:"full_name"`
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

// GitHubErrorResponse contains the GitHub API error response.
type GitHubErrorResponse struct {
	Message string `json:"message"`
	Errors  []ErrorItem
}

// ErrorItem represents the an error detail on the GitHub API error response.
type ErrorItem struct {
	Resource string
	Field    string
	Code     string
}
