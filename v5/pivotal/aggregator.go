package pivotal

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

const aggregatorURL = "https://www.pivotaltracker.com/services/v5/aggregator"

// AggregatorService is used to wrap the client.
type AggregatorService struct {
	client *Client
}

// Aggregation is the data object for an aggregation.
// It is used for storing the urls of the requests and the response.
type Aggregation struct {
	// Requests are the GET requests that are bundled by the aggregator.
	requests []string

	// Storing the response json along with the urls.
	aggregatedResponse map[string]interface{}

	// The service of the aggregation
	service *AggregatorService

	// The amount of requests that an aggregation bundle will contain.
	requestsPerAggregation int
}

// GetBuilder returns a builder for building an aggregation.
func (a *AggregatorService) GetBuilder() *Aggregation {
	aggregation := Aggregation{
		aggregatedResponse:     make(map[string]interface{}),
		service:                a,
		requestsPerAggregation: 15, // The default value is 15.
	}
	return &aggregation
}

// SetRequestsPerAggregation sets the number of that each aggregation
// bundle will contain.
func (a *Aggregation) SetRequestsPerAggregation(amount int) {
	a.requestsPerAggregation = amount
}

func newAggregatorService(client *Client) *AggregatorService {
	return &AggregatorService{client}
}

// Story adds a story request to the aggregation.
func (a *Aggregation) Story(projectID, storyID int) *Aggregation {
	a.requests = append(a.requests, buildStoryURL(projectID, storyID))
	return a
}

// StoryByID adds a story request using only the story ID.
func (a *Aggregation) StoryByID(storyID int) *Aggregation {
	a.requests = append(a.requests, buildStoryURLByID(storyID))
	return a
}

// Stories adds a list of story requests to an aggregation.
func (a *Aggregation) Stories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.Story(projectID, storyID)
	}
	return a
}

// CommentsOfStory adds a request for getting the comments of a story.
func (a *Aggregation) CommentsOfStory(projectID, storyID int) {
	a.requests = append(a.requests, buildCommentsURL(projectID, storyID))
}

// CommentsOfStories adds multiple requests for getting the comments of a list of stories.
func (a *Aggregation) CommentsOfStories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.CommentsOfStory(projectID, storyID)
	}
	return a
}

// ReviewsOfStory adds multiple requests for getting the reviews of a story.
func (a *Aggregation) ReviewsOfStory(projectID, storyID int) *Aggregation {
	a.requests = append(a.requests, buildReviewsURL(projectID, storyID))
	return a
}

// ReviewsOfStories adds multiple requests for getting the reviews of multiple stories.
func (a *Aggregation) ReviewsOfStories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.ReviewsOfStory(projectID, storyID)
	}
	return a
}

func buildStoryURLByID(storyID int) string {
	return fmt.Sprintf("/services/v5/stories/%d", storyID)
}

func buildStoryURL(projectID, storyID int) string {
	return fmt.Sprintf("/services/v5/projects/%d/stories/%d", projectID, storyID)
}

func buildCommentsURL(projectID, storyID int) string {
	return fmt.Sprintf("/services/v5/projects/%d/stories/%d/comments", projectID, storyID)
}

func buildReviewsURL(projectID, storyID int) string {
	return fmt.Sprintf("/services/v5/projects/%d/stories/%d/reviews?fields=id,story_id,review_type,review_type_id,reviewer_id,status,created_at,updated_at,kind", projectID, storyID)
}

func maxPagesPagination(total int, perPage int) int {
	pagesNumber := math.Ceil(float64(total) / float64(perPage))
	return int(pagesNumber)
}

func getLastIndex(reqs []string, perPage int) int {
	lastIndex := perPage
	if lastIndex > len(reqs) {
		lastIndex = len(reqs)
	}
	return lastIndex
}

// Send sends the next bulk of aggregation to Pivotal Tracker.
// It returns true, if there are more bulks to be sent.
func (a *Aggregation) Send() (*Aggregation, *http.Response, bool, error) {
	lastIndex := getLastIndex(a.requests, a.requestsPerAggregation)
	if lastIndex == 0 {
		return a, nil, false, nil
	}

	// Getting the current bulk to be sent
	currentBulk := a.requests[0:lastIndex]

	// Popping the current bulk from the requests.
	a.requests = a.requests[lastIndex:]

	aggregatedResponse := make(map[string]interface{})

	req, err := a.service.client.NewRequest("POST", aggregatorURL, currentBulk)
	if err != nil {
		return nil, nil, false, err
	}

	response, err := a.service.client.Do(req, &aggregatedResponse)
	if err != nil {
		return nil, nil, false, err
	}

	// Appending the response body into the current aggregation for getters.
	for url, response := range aggregatedResponse {
		a.aggregatedResponse[url] = response
	}

	// Return true if there is more left
	if len(a.requests) > 0 {
		return a, response, true, nil
	}
	return a, response, false, nil
}

// Send sends all the bulks to PivotalTracker.
func (a *Aggregation) SendAll() (*Aggregation, []*http.Response, error) {
	responses := []*http.Response{}
	for {
		a, response, hasMore, err := a.Send()

		if err != nil {
			return a, nil, err
		}

		responses = append(responses, response)

		if !hasMore {
			return a, responses, nil
		}
	}
}

func mapTo(from interface{}, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

// GetStoryByID returns the story using only the story ID.
// Please use GetStory, if you used the project ID when adding the request.
func (a *Aggregation) GetStoryByID(storyID int) (*Story, error) {
	u := buildStoryURLByID(storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist.", storyID)
	}

	var story Story
	err := mapTo(response, &story)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// GetStory returns the story using both the project ID and the story ID.
// Please use GetStoryByID, if you used the story ID only when adding the request.
func (a *Aggregation) GetStory(projectID, storyID int) (*Story, error) {
	u := buildStoryURL(projectID, storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist for project %d.", storyID, projectID)
	}

	var story Story
	err := mapTo(response, &story)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// GetComments returns the comments of a story.
func (a *Aggregation) GetComments(projectID, storyID int) ([]Comment, error) {
	u := buildCommentsURL(projectID, storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d comments don't exist for project %d.", storyID, projectID)
	}
	var comments []Comment
	err := mapTo(response, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetReviews returns the reviews of a story.
func (a *Aggregation) GetReviews(projectID, storyID int) ([]Review, error) {
	u := buildReviewsURL(projectID, storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d reviews don't exist for project %d.", storyID, projectID)
	}
	var reviews []Review
	err := mapTo(response, &reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
