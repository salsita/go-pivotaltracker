// Copyright (c) 2014-2018 Salsita Software
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	// LibraryVersion is used to give the UserAgent some additional context
	LibraryVersion = "2.0.0"

	defaultBaseURL   = "https://www.pivotaltracker.com/services/v5/"
	defaultUserAgent = "go-pivotaltracker/" + LibraryVersion
)

// ErrNoTrailingSlash is returned when URLs are missing the trailing slash
var ErrNoTrailingSlash = errors.New("trailing slash missing")

// Client wraps all Pivotal Tracker services with necessary auth contexts
type Client struct {
	// Pivotal Tracker access token to be used to authenticate API requests.
	token string

	// HTTP client to be used for communication with the Pivotal Tracker API.
	client *http.Client

	// Base URL of the Pivotal Tracker API that is to be used to form API requests.
	baseURL *url.URL

	// User-Agent header to use when connecting to the Pivotal Tracker API.
	userAgent string

	// Accounts service
	Accounts *AccountsService

	// Me service
	Me *MeService

	// Project service
	Projects *ProjectService

	// Story service
	Stories *StoryService

	// Membership service
	Memberships *MembershipService

	// Iteration service
	Iterations *IterationService

	// Activity Service
	Activity *ActivityService

	// Epic Service
	Epic *EpicService

	// Aggregator Service
	Aggregator *AggregatorService
}

// NewClient takes a Pivotal Tracker API Token (created from the project settings) and
// returns a default Client implementation
func NewClient(apiToken string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	client := &Client{
		token:     apiToken,
		client:    http.DefaultClient,
		baseURL:   baseURL,
		userAgent: defaultUserAgent,
	}
	client.Accounts = newAccountsService(client)
	client.Me = newMeService(client)
	client.Projects = newProjectService(client)
	client.Stories = newStoryService(client)
	client.Memberships = newMembershipService(client)
	client.Iterations = newIterationService(client)
	client.Activity = newActivitiesService(client)
	client.Epic = newEpicService(client)
	client.Aggregator = newAggregatorService(client)
	return client
}

// SetBaseURL overrides the defaultBaseURL in the default Client implementation.
func (c *Client) SetBaseURL(baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	if u.Path != "" && u.Path[len(u.Path)-1] != '/' {
		return ErrNoTrailingSlash
	}

	c.baseURL = u
	return nil
}

// SetHTTPClient overrides the default HTTP Client, http.DefaultClient.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// SetUserAgent overrides the defaultUserAgent in the default Client implementation.
func (c *Client) SetUserAgent(agent string) {
	c.userAgent = agent
}

// NewRequest takes an HTTP request definition and wraps it with the Client context.
func (c *Client) NewRequest(method, urlPath string, body interface{}) (*http.Request, error) {
	path, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(path)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("X-TrackerToken", c.token)
	return req, nil
}

// Do takes a request created from NewRequest and executes the HTTP round trip action.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var errObject Error
		if err := json.NewDecoder(resp.Body).Decode(&errObject); err != nil {
			return resp, &ErrAPI{Response: resp}
		}

		return resp, &ErrAPI{
			Response: resp,
			Err:      &errObject,
		}
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
