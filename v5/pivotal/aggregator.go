package pivotal

import (
	"encoding/json"
	"fmt"
	"math"
)

const aggregatorURL = "https://www.pivotaltracker.com/services/v5/aggregator"

const requestsPerAggregation = 15 // Number of requests per aggregation.

// AggregatorService is used to wrap the client.
type AggregatorService struct {
	client *Client
}

// Aggregation is the data object for an aggregation.
// It is used for storing the urls of the requests and the response.
type Aggregation struct {
	// Requests are the GET requests that are used by the aggregator.
	requests []string

	// Storing the response json along with the urls.
	aggregatedResponse map[string]interface{}

	// The service of the aggregation
	service *AggregatorService
}

// GetBuilder returns a builder for building an aggregation.
func (a *AggregatorService) GetBuilder() *Aggregation {
	aggregation := Aggregation{
		aggregatedResponse: make(map[string]interface{}),
		service:            a,
	}
	return &aggregation
}

func newAggregatorService(client *Client) *AggregatorService {
	return &AggregatorService{client}
}

// Story adds a story request to the aggregation.
func (a *Aggregation) Story(projectID, storyID int) *Aggregation {
	a.requests = append(a.requests, buildStoryURL(projectID, storyID))
	return a
}

// StoryUsingStoryID adds a story request using only using the story ID.
func (a *Aggregation) StoryUsingStoryID(storyID int) *Aggregation {
	a.requests = append(a.requests, buildStoryURLOnlyUsingStoryID(storyID))
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

func buildStoryURLOnlyUsingStoryID(storyID int) string {
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

func paginate(reqs []string, currentPage, total int, perPage int) []string {
	firstEntry := (currentPage - 1) * perPage
	lastEntry := firstEntry + perPage

	if lastEntry > total {
		lastEntry = total
	}

	return reqs[firstEntry:lastEntry]
}

// Send completes the aggregation and sends it to Pivotal Tracker.
func (a *Aggregation) Send() (*Aggregation, error) {
	N := len(a.requests)
	max := maxPagesPagination(N, requestsPerAggregation)
	for currentPage := 1; currentPage <= max; currentPage++ {
		currentReqs := paginate(a.requests, currentPage, N, requestsPerAggregation)

		aggregatedResponse := make(map[string]interface{})

		req, err := a.service.client.NewRequest("POST", aggregatorURL, currentReqs)
		if err != nil {
			return nil, err
		}

		_, err = a.service.client.Do(req, &aggregatedResponse)
		if err != nil {
			return nil, err
		}

		// Appending the response body into the current aggregation for getters.
		for url, response := range aggregatedResponse {
			a.aggregatedResponse[url] = response
		}
	}

	return a, nil
}

// GetStoryOnlyUsingStoryID returns the story using only the story ID.
func (a *Aggregation) GetStoryOnlyUsingStoryID(storyID int) (*Story, error) {
	u := buildStoryURLOnlyUsingStoryID(storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist.", storyID)
	}
	byteData, _ := json.Marshal(response)

	// Handling get story requests if it isn't comments/reviews.
	var story Story
	err := json.Unmarshal(byteData, &story)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// GetStory returns the story using both the project ID and the story ID.
func (a *Aggregation) GetStory(projectID, storyID int) (*Story, error) {
	u := buildStoryURL(projectID, storyID)
	response, ok := a.aggregatedResponse[u]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist for project %d.", storyID, projectID)
	}
	byteData, _ := json.Marshal(response)

	// Handling get story requests if it isn't comments/reviews.
	var story Story
	err := json.Unmarshal(byteData, &story)
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
	byteData, _ := json.Marshal(response)

	var comments []Comment
	err := json.Unmarshal(byteData, &comments)
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
	byteData, _ := json.Marshal(response)
	var reviews []Review
	err := json.Unmarshal(byteData, &reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
