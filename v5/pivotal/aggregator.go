package pivotal

import (
	"encoding/json"
	"fmt"
	"strings"
)

const aggregatorURL = "https://www.pivotaltracker.com/services/v5/aggregator"

type AggregatorService struct {
	client             *Client
	requests           []string
	aggregatedResponse map[string]interface{} // Storing the response temporarily.
}

func newAggregatorService(client *Client) *AggregatorService {
	return &AggregatorService{client, []string{}, nil}
}

// For making sure the request is in empty state
func (service *AggregatorService) InitRequest() {
	service.client.Aggregator.requests = []string{}
	service.client.Aggregator.aggregatedResponse = make(map[string]interface{})
}

func (service *AggregatorService) QueueStoryCommentsByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d/comments", projectID, storyID)
	service.requests = append(service.requests, u)
}

func (service *AggregatorService) QueueStoryByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d", projectID, storyID)
	service.requests = append(service.requests, u)
}

func (service *AggregatorService) QueueReviewsByID(projectID, storyID int) {
	u := fmt.Sprintf("/services/v5/projects/%d/stories/%d/reviews?fields=id,story_id,review_type,review_type_id,reviewer_id,status,created_at,updated_at,kind", projectID, storyID)
	service.requests = append(service.requests, u)
}

func (service *AggregatorService) FinishRequest() error {
	req, err := service.client.NewRequest("POST", aggregatorURL, service.requests)
	if err != nil {
		return err
	}

	_, err = service.client.Do(req, &service.aggregatedResponse)
	if err != nil {
		return err
	}
	return nil
}

func (service *AggregatorService) GetStoryReviewsFromRequest() (map[int][]Review, error) {
	// For storing the comments using story IDs directly
	reviewsMap := make(map[int][]Review)

	for url, byteStory := range service.aggregatedResponse {
		// Skip the requests that are not comments
		if !strings.Contains(url, "reviews") {
			continue
		}
		byteData, _ := json.Marshal(byteStory)

		var reviews []Review
		err := json.Unmarshal(byteData, &reviews)
		if err != nil {
			return nil, err
		}
		if len(reviews) == 0 {
			break
		}

		// Updating the commentsMap using the StoryID
		reviewsMap[reviews[0].StoryID] = reviews
	}
	return reviewsMap, nil

}

func (service *AggregatorService) GetStoryCommentsFromRequest() (map[int][]Comment, error) {
	// For storing the comments using story IDs directly
	commentsMap := make(map[int][]Comment)

	for url, byteStory := range service.aggregatedResponse {
		// Skip the requests that are not comments
		if !strings.Contains(url, "comments") {
			continue
		}
		byteData, _ := json.Marshal(byteStory)

		var comments []Comment
		err := json.Unmarshal(byteData, &comments)
		if err != nil {
			return nil, err
		}
		if len(comments) == 0 {
			break
		}

		// Updating the commentsMap using the StoryID
		commentsMap[comments[0].StoryID] = comments
	}
	return commentsMap, nil

}

func (service *AggregatorService) GetStoriesFromRequest() (map[int]Story, error) {
	// For storing the stories using IDs directly
	storyIDMap := make(map[int]Story)

	for url, byteStory := range service.aggregatedResponse {
		// Skip the requests that are not labels.
		if !strings.Contains(url, "stories") {
			continue
		}
		byteData, err := json.Marshal(byteStory)
		if err != nil {
			return nil, err
		}

		var story Story
		err = json.Unmarshal(byteData, &story)
		if err != nil {
			return nil, err
		}

		storyIDMap[story.ID] = story
	}
	return storyIDMap, nil

}
