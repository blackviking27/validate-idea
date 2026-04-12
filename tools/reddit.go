package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RedditPageJson []struct {
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Title    string `json:"title"`
				SelfText string `json:"selftext"`
				Body     string `json:"body"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditSearchResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Permalink string `json:"permalink"`
			}
		} `json:"children"`
	} `json:"data"`
}

func getSearchResultUrls(query string) (urls []string) {
	redditSearchPageUrl := fmt.Sprintf("https://www.reddit.com/search.json?q=%s", query)

	req, err := http.NewRequest("GET", redditSearchPageUrl, nil)
	if err != nil {
		fmt.Errorf("Unable to create reddit search request")
		return urls
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to search reddit for query: ", query)
		return urls
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Errorf("Error while fetching the data, status: %i", res.StatusCode)
		return urls
	}

	var redditSearchResponse RedditSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&redditSearchResponse); err != nil {
		fmt.Errorf("Unable to parse search page respons, err: %v", err)
		return urls
	}

	for _, data := range redditSearchResponse.Data.Children {
		urls = append(urls, fmt.Sprintf("https://www.reddit.com%s", data.Data.Permalink))
	}

	return urls
}

func getPostData(url string) (ParsedSearchResult, error) {

	result := ParsedSearchResult{}

	pageUrl := strings.Split(url, "?")[0]
	pageUrl = strings.TrimSuffix(pageUrl, "/")

	jsonPageUrl := pageUrl + ".json"

	req, _ := http.NewRequest("GET", jsonPageUrl, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return result, fmt.Errorf("Unable to fetch content for %s, err: %v", pageUrl, err)
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var pageResponse RedditPageJson
	if err := json.Unmarshal(bodyBytes, &pageResponse); err != nil {
		return result, fmt.Errorf("Unable to parse the page %s, err: %v", pageUrl, err)
	}

	if len(pageResponse) < 2 {
		return result, fmt.Errorf("Unexpected data received for url: %s", pageUrl)
	}

	// Grab title and content
	postDetails := pageResponse[0].Data.Children
	if len(postDetails) > 0 {
		result.Title = postDetails[0].Data.Title
		result.Content = postDetails[0].Data.SelfText
	}

	// Getting the comment details
	comments := pageResponse[1].Data.Children
	for _, cmt := range comments {
		if cmt.Kind == "t1" && cmt.Data.Body != "" {
			result.Comments = append(result.Comments, cmt.Data.Body)
		}
	}

	return result, nil
}

type Reddit struct {
}

func (this *Reddit) Search(ctx context.Context, query string) ([]ParsedSearchResult, error) {
	fmt.Printf("[1] Searcing the reddit with query: %s\n", query)
	postUrls := getSearchResultUrls(query)

	fmt.Println("[2] Scraping the search post data")
	var parsedSearchResults []ParsedSearchResult
	for _, url := range postUrls {

		post, err := getPostData(url)
		if err != nil {
			fmt.Sprintf("Skipping url %s for err: %v\n", url, err)
			continue
		}

		parsedSearchResults = append(parsedSearchResults, post)

		time.Sleep(2 * time.Second)
	}
	return parsedSearchResults, nil
}

func NewRedditSearch() *Reddit {
	return &Reddit{}
}
