// Copyright (c) 2014-2018 Salsita Software
// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"net/http"
	"net/url"
)

type SearchResponse struct {
	Stories struct {
		Stories []*Story `json:"stories"`
	} `json:"stories"`
}

// SearchService wraps the client context and allows for interaction
// with the Pivotal Tracker Search API.
type SearchService struct {
	client *Client
}

func newSearchService(client *Client) *SearchService {
	return &SearchService{client}
}

// Search searches the project data and returns the stories and/or epics matching the query
func (service *SearchService) Search(projectID int, query string) (SearchResponse, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/search?query=%s", projectID, url.QueryEscape(query))
	req, _ := service.client.NewRequest("GET", u, nil)
	var searchResponse SearchResponse
	resp, err := service.client.Do(req, &searchResponse)

	return searchResponse, resp, err
}
