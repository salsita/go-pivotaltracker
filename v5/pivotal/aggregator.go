package pivotal

import (
	"encoding/json"
	"fmt"
	"strings"
)

const aggregatorURL = "https://www.pivotaltracker.com/services/v5/aggregator"

// AggregatorService is used to wrap the client.
type AggregatorService struct {
	client *Client
}

// Aggregation object stores the state of the aggregation.
type Aggregation struct {
	// Requests are the GET requests that are used by the aggregator.
	requests []string

	// Storing the response json along with the urls.
	aggregatedResponse map[string]interface{}

	// Map for quickly accessing the stories using their IDs
	aggregatedData map[int]*AggregatedStory
}

// AggregatedStory contains the story, comment and review information.
type AggregatedStory struct {
	Story    *Story
	Comments *[]Comment
	Reviews  *[]Review
}

func newAggregatorService(client *Client) *AggregatorService {
	return &AggregatorService{client}
}

func (a *Aggregation) maybeCreateAggregatedStory(StoryID int) *AggregatedStory {
	story, ok := a.aggregatedData[StoryID]
	if !ok {
		a.aggregatedData[StoryID] = &AggregatedStory{}
		return a.aggregatedData[StoryID]
	}
	return story
}

// Fills the urls for the aggregator and sends the request to PT.
func (service *AggregatorService) fillAggregationData(ProjectID int, StoryIDs []int) (*Aggregation, error) {
	aggregation := &Aggregation{
		aggregatedResponse: make(map[string]interface{}),
		aggregatedData:     make(map[int]*AggregatedStory),
	}
	for _, ID := range StoryIDs {
		aggregation.queueStoryByID(ProjectID, ID)
		aggregation.queueStoryCommentsByID(ProjectID, ID)
		aggregation.queueReviewsByID(ProjectID, ID)
	}
	err := service.sendAggregationRequest(aggregation)
	if err != nil {
		return nil, err
	}
	return aggregation, nil
}

// Converts the JSON response from PT into the related structs(Story, Comment, Review).
func (service *AggregatorService) processAggregationMaps(aggregation *Aggregation) (*Aggregation, error) {
	// Creating the map with the story information.
	for url, byteStory := range aggregation.aggregatedResponse {
		byteData, _ := json.Marshal(byteStory)

		if strings.Contains(url, "reviews") {
			var reviews []Review
			err := json.Unmarshal(byteData, &reviews)
			if err != nil {
				return nil, err
			}
			if len(reviews) == 0 {
				continue
			}
			// Accessing the reviews using the story ID
			storyData := aggregation.maybeCreateAggregatedStory(reviews[0].StoryID)
			storyData.Reviews = &reviews
			continue
		}

		if strings.Contains(url, "comments") {
			var comments []Comment
			err := json.Unmarshal(byteData, &comments)
			if err != nil {
				return nil, err
			}
			if len(comments) == 0 {
				continue
			}
			// Accessing the reviews using the story ID
			storyData := aggregation.maybeCreateAggregatedStory(comments[0].StoryID)
			storyData.Comments = &comments
			continue
		}

		// Handling get story requests if it isn't comments/reviews.
		var story Story
		err := json.Unmarshal(byteData, &story)
		if err != nil {
			return nil, err
		}
		if story.ID == 0 {
			continue
		}
		// Accessing the story with the story ID
		storyData := aggregation.maybeCreateAggregatedStory(story.ID)
		storyData.Story = &story

	}

	return aggregation, nil
}

// BuildStoriesCommentsAndReviewsFor returns the reviews, comments and the story information
// for the given project and story IDs.
func (service *AggregatorService) BuildStoriesCommentsAndReviewsFor(ProjectID int, StoryIDs []int) (*Aggregation, error) {
	if len(StoryIDs) > 5 {
		// Aggregator doesn't seem to be returning more than 15 requests in total.
		return nil, fmt.Errorf("There can be a maximum number of 5 Story IDs.")
	}

	aggregation, err := service.fillAggregationData(ProjectID, StoryIDs)
	if err != nil {
		return nil, err
	}

	aggregation, err = service.processAggregationMaps(aggregation)
	if err != nil {
		return nil, err
	}

	return aggregation, nil
}

// Adds the url for the comments to the aggregation data.
func (a *Aggregation) queueStoryCommentsByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d/comments", projectID, storyID)
	a.requests = append(a.requests, u)
}

// Adds the url for the story to the aggregation data.
func (a *Aggregation) queueStoryByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d", projectID, storyID)
	a.requests = append(a.requests, u)
}

// Adds the url for the reviews to the aggregation data.
func (a *Aggregation) queueReviewsByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d/reviews?fields=id,story_id,review_type,review_type_id,reviewer_id,status,created_at,updated_at,kind", projectID, storyID)
	a.requests = append(a.requests, u)
}

// Sends the request for the aggregation.
func (service *AggregatorService) sendAggregationRequest(a *Aggregation) error {
	req, err := service.client.NewRequest("POST", aggregatorURL, a.requests)
	if err != nil {
		return err
	}

	_, err = service.client.Do(req, &a.aggregatedResponse)
	if err != nil {
		return err
	}
	return nil
}

// Returns the story using StoryID from the aggregation.
func (a *Aggregation) GetStory(storyID int) (*Story, error) {
	story, ok := a.aggregatedData[storyID]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist.", storyID)
	}
	return story.Story, nil
}

// Returns the comments using story id from the aggregation.
func (a *Aggregation) GetComments(storyID int) ([]Comment, error) {
	story, ok := a.aggregatedData[storyID]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist.", storyID)
	}
	return *story.Comments, nil
}

// Returns the reviews using the story id from the aggregation.
func (a *Aggregation) GetReviews(storyID int) ([]Review, error) {
	story, ok := a.aggregatedData[storyID]
	if !ok {
		return nil, fmt.Errorf("Story %d doesn't exist.", storyID)
	}
	return *story.Reviews, nil
}
