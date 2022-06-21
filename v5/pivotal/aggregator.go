package pivotal

import (
	"encoding/json"
	"fmt"
	"math"
)

const aggregatorURL = "https://www.pivotaltracker.com/services/v5/aggregator"

// PT doesn't allow more than around 15 urls in one aggregation request.
// However, this was found by experimenting. There was no official
// rate-limiting data in their docs.
const perPage = 15 // GET requests per aggregation requests.

// AggregatorService is used to wrap the client.
type AggregatorService struct {
	client *Client
}

type AggregationRequest struct {
	url       string
	projectID int
	storyID   int
}

// Aggregation object stores the state of the aggregation.
type Aggregation struct {
	// Requests are the GET requests that are used by the aggregator.
	requests []string

	// Storing the response json along with the urls.
	aggregatedResponse map[string]interface{}

	// The service of the aggregation
	service *AggregatorService
}

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

// Adds the url for the story to the aggregation data.
func (a *Aggregation) Story(projectID, storyID int) *Aggregation {
	a.requests = append(a.requests, BuildStoryURL(projectID, storyID))
	return a
}

func (a *Aggregation) StoryUsingStoryID(storyID int) *Aggregation {
	a.requests = append(a.requests, BuildStoryURLOnlyUsingStoryID(storyID))
	return a
}

// Stories adds the urls for the slice of story IDs.
func (a *Aggregation) Stories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.Story(projectID, storyID)
	}
	return a
}

// Adds the url for the comments to the aggregation data.
func (a *Aggregation) CommentsOfStory(projectID, storyID int) {
	a.requests = append(a.requests, buildCommentsURL(projectID, storyID))
}

// Comments adds the requests for getting the comments of the stories in storiesToGet.
func (a *Aggregation) CommentsOfStories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.CommentsOfStory(projectID, storyID)
	}
	return a
}

// Adds the url for the reviews to the aggregation data.
func (a *Aggregation) ReviewsOfStory(projectID, storyID int) *Aggregation {
	a.requests = append(a.requests, buildReviewsURL(projectID, storyID))
	return a
}

// Comments adds the requests for getting the comments of the stories in storiesToGet.
func (a *Aggregation) ReviewsOfStories(projectID int, storyIDs []int) *Aggregation {
	for _, storyID := range storyIDs {
		a.ReviewsOfStory(projectID, storyID)
	}
	return a
}

func BuildStoryURLOnlyUsingStoryID(storyID int) string {
	return fmt.Sprintf("/services/v5/stories/%d", storyID)
}

func BuildStoryURL(projectID, storyID int) string {
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

// Sends the request for the aggregation.
func (a *Aggregation) Send() (*Aggregation, error) {
	N := len(a.requests)
	max := maxPagesPagination(N, perPage)
	for currentPage := 1; currentPage <= max; currentPage++ {
		currentReqs := paginate(a.requests, currentPage, N, perPage)

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

// Returns the story using StoryID from the aggregation.
func (a *Aggregation) GetStoryOnlyUsingStoryID(storyID int) (*Story, error) {
	u := BuildStoryURLOnlyUsingStoryID(storyID)
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

// Returns the story using StoryID from the aggregation.
func (a *Aggregation) GetStory(projectID, storyID int) (*Story, error) {
	u := BuildStoryURL(projectID, storyID)
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

// Returns the comments using story id from the aggregation.
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

// Returns the reviews using the story id from the aggregation.
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
