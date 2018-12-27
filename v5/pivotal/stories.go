// Copyright (c) 2014-2018 Salsita Software
// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// PageLimit is the number of items to fetch at once when getting paginated response.
const PageLimit = 10

const (
	// StoryTypeFeature wraps the string enum in the variable name.
	StoryTypeFeature = "feature"
	// StoryTypeBug wraps the string enum in the variable name.
	StoryTypeBug = "bug"
	// StoryTypeChore wraps the string enum in the variable name.
	StoryTypeChore = "chore"
	// StoryTypeRelease wraps the string enum in the variable name.
	StoryTypeRelease = "release"
)

const (
	// StoryStateUnscheduled wraps the story state enum in the variable name.
	StoryStateUnscheduled = "unscheduled"
	// StoryStatePlanned wraps the story state enum in the variable name.
	StoryStatePlanned = "planned"
	// StoryStateUnstarted wraps the story state enum in the variable name.
	StoryStateUnstarted = "unstarted"
	// StoryStateStarted wraps the story state enum in the variable name.
	StoryStateStarted = "started"
	// StoryStateFinished wraps the story state enum in the variable name.
	StoryStateFinished = "finished"
	// StoryStateDelivered wraps the story state enum in the variable name.
	StoryStateDelivered = "delivered"
	// StoryStateAccepted wraps the story state enum in the variable name.
	StoryStateAccepted = "accepted"
	// StoryStateRejected wraps the story state enum in the variable name.
	StoryStateRejected = "rejected"
)

// Story is the top level data object for a story, it wraps multiple child objects
// but is the primary required for interacting with the StoryService.
type Story struct {
	ID            int        `json:"id,omitempty"`
	ProjectID     int        `json:"project_id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	Type          string     `json:"story_type,omitempty"`
	State         string     `json:"current_state,omitempty"`
	Estimate      *float64   `json:"estimate,omitempty"`
	AcceptedAt    *time.Time `json:"accepted_at,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	RequestedByID int        `json:"requested_by_id,omitempty"`
	OwnerIDs      []int      `json:"owner_ids,omitempty"`
	LabelIDs      []int      `json:"label_ids,omitempty"`
	Labels        []*Label   `json:"labels,omitempty"`
	TaskIDs       []int      `json:"task_ids,omitempty"`
	Tasks         []int      `json:"tasks,omitempty"`
	FollowerIDs   []int      `json:"follower_ids,omitempty"`
	CommentIDs    []int      `json:"comment_ids,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	BeforeID      int        `json:"before_id,omitempty"`
	AfterID       int        `json:"after_id,omitempty"`
	IntegrationID int        `json:"integration_id,omitempty"`
	ExternalID    string     `json:"external_id,omitempty"`
	URL           string     `json:"url,omitempty"`
}

// StoryRequest is a simplified Story object for use in Create/Update/Delete operations.
type StoryRequest struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"story_type,omitempty"`
	State       string    `json:"current_state,omitempty"`
	Estimate    *float64  `json:"estimate,omitempty"`
	OwnerIDs    *[]int    `json:"owner_ids,omitempty"`
	LabelIDs    *[]int    `json:"label_ids,omitempty"`
	Labels      *[]*Label `json:"labels,omitempty"`
	TaskIDs     *[]int    `json:"task_ids,omitempty"`
	Tasks       *[]int    `json:"tasks,omitempty"`
	FollowerIDs *[]int    `json:"follower_ids,omitempty"`
	CommentIDs  *[]int    `json:"comment_ids,omitempty"`
}

// Label is a child object of a Story. This may need to be broken out into a LabelService
// someday but for now is ok here.
type Label struct {
	ID        int        `json:"id,omitempty"`
	ProjectID int        `json:"project_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Kind      string     `json:"kind,omitempty"`
}

// Task is a child object of a Story.
type Task struct {
	ID          int        `json:"id,omitempty"`
	StoryID     int        `json:"story_id,omitempty"`
	Description string     `json:"description,omitempty"`
	Position    int        `json:"position,omitempty"`
	Complete    bool       `json:"complete,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Person is a child object of Story to give assigned/reporter values.
type Person struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Initials string `json:"initials,omitempty"`
	Username string `json:"username,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

// Comment is used to show all comments associated with a Story.
type Comment struct {
	ID                  int        `json:"id,omitempty"`
	StoryID             int        `json:"story_id,omitempty"`
	EpicID              int        `json:"epic_id,omitempty"`
	PersonID            int        `json:"person_id,omitempty"`
	Text                string     `json:"text,omitempty"`
	FileAttachmentIDs   []int      `json:"file_attachment_ids,omitempty"`
	GoogleAttachmentIDs []int      `json:"google_attachment_ids,omitempty"`
	CommitType          string     `json:"commit_type,omitempty"`
	CommitIdentifier    string     `json:"commit_identifier,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

// Blocker shows the relationship between other Stories and blocking states.
type Blocker struct {
	ID          int        `json:"id,omitempty"`
	StoryID     int        `json:"story_id,omitempty"`
	PersonID    int        `json:"person_id,omitempty"`
	Description string     `json:"description,omitempty"`
	Resolved    bool       `json:"resolved,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// BlockerRequest is used to do Create/Update/Delete on blockers.
type BlockerRequest struct {
	Description string `json:"description,omitempty"`
	Resolved    *bool  `json:"resolved,omitempty"`
}

// StoryService wraps the client context and allows for interaction
// with the Pivotal Tracker Story API.
type StoryService struct {
	client *Client
}

func newStoryService(client *Client) *StoryService {
	return &StoryService{client}
}

// List returns all stories matching the filter in case the filter is specified.
//
// List actually sends 2 HTTP requests - one to get the total number of stories,
// another to retrieve the stories using the right pagination setup. The reason
// for this is that the filter might require to fetch all the stories at once
// to get the right results. Since the response as generated by Pivotal Tracker
// is not always sorted when using a filter, this approach is required to get
// the right data. Not sure whether this is a bug or a feature.
func (service *StoryService) List(projectID int, filter string) ([]*Story, error) {
	reqFunc := newStoriesRequestFunc(service.client, projectID, filter)
	cursor, err := newCursor(service.client, reqFunc, 0)
	if err != nil {
		return nil, err
	}

	var stories []*Story
	if err := cursor.all(&stories); err != nil {
		return nil, err
	}
	return stories, nil
}

func newStoriesRequestFunc(client *Client, projectID int, filter string) func() *http.Request {
	return func() *http.Request {
		u := fmt.Sprintf("projects/%v/stories", projectID)
		if filter != "" {
			u += "?filter=" + url.QueryEscape(filter)
		}
		req, _ := client.NewRequest("GET", u, nil)
		return req
	}
}

// StoryCursor is used to implement the iterator pattern.
type StoryCursor struct {
	*cursor
	buff []*Story
}

// Next returns the next story.
//
// In case there are no more stories, io.EOF is returned as an error.
func (c *StoryCursor) Next() (s *Story, err error) {
	if len(c.buff) == 0 {
		_, err = c.next(&c.buff)
		if err != nil {
			return nil, err
		}
	}

	if len(c.buff) == 0 {
		err = io.EOF
	} else {
		s, c.buff = c.buff[0], c.buff[1:]
	}
	return s, err
}

// Iterate returns a cursor that can be used to iterate over the stories specified
// by the filter. More stories are fetched on demand as needed.
func (service *StoryService) Iterate(projectID int, filter string) (c *StoryCursor, err error) {
	reqFunc := newStoriesRequestFunc(service.client, projectID, filter)
	cursor, err := newCursor(service.client, reqFunc, PageLimit)
	if err != nil {
		return nil, err
	}
	return &StoryCursor{cursor, make([]*Story, 0)}, nil
}

// Create is used to make a new Story.
func (service *StoryService) Create(projectID int, story *StoryRequest) (*Story, *http.Response, error) {
	if projectID == 0 {
		return nil, nil, &ErrFieldNotSet{"project_id"}
	}

	if story.Name == "" {
		return nil, nil, &ErrFieldNotSet{"name"}
	}

	u := fmt.Sprintf("projects/%v/stories", projectID)
	req, err := service.client.NewRequest("POST", u, story)
	if err != nil {
		return nil, nil, err
	}

	var newStory Story

	resp, err := service.client.Do(req, &newStory)
	if err != nil {
		return nil, resp, err
	}

	return &newStory, resp, nil
}

// Get will obtain the details about a single Story by ID.
func (service *StoryService) Get(projectID, storyID int) (*Story, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v", projectID, storyID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var story Story
	resp, err := service.client.Do(req, &story)
	if err != nil {
		return nil, resp, err
	}

	return &story, resp, err
}

// Update will change details of an existing story.
func (service *StoryService) Update(projectID, storyID int, story *StoryRequest) (*Story, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v", projectID, storyID)
	req, err := service.client.NewRequest("PUT", u, story)
	if err != nil {
		return nil, nil, err
	}

	var bodyStory Story
	resp, err := service.client.Do(req, &bodyStory)
	if err != nil {
		return nil, resp, err
	}

	return &bodyStory, resp, err

}

// ListTasks will get the Tasks associated with a Story by ID.
func (service *StoryService) ListTasks(projectID, storyID int) ([]*Task, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v/tasks", projectID, storyID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tasks []*Task
	resp, err := service.client.Do(req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

// AddTask will add a new Task to a Story by ID.
func (service *StoryService) AddTask(projectID, storyID int, task *Task) (*http.Response, error) {
	if task.Description == "" {
		return nil, &ErrFieldNotSet{"description"}
	}

	u := fmt.Sprintf("projects/%v/stories/%v/tasks", projectID, storyID)
	req, err := service.client.NewRequest("POST", u, task)
	if err != nil {
		return nil, err
	}

	return service.client.Do(req, nil)
}

// ListOwners will show who is assigned to a story, returning a Person array.
func (service *StoryService) ListOwners(projectID, storyID int) ([]*Person, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/stories/%d/owners", projectID, storyID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var owners []*Person
	resp, err := service.client.Do(req, &owners)
	if err != nil {
		return nil, resp, err
	}

	return owners, resp, err
}

// AddComment will take a Comment object and attach it to a Story.
func (service *StoryService) AddComment(
	projectID int,
	storyID int,
	comment *Comment,
) (*Comment, *http.Response, error) {

	u := fmt.Sprintf("projects/%v/stories/%v/comments", projectID, storyID)
	req, err := service.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, nil, err
	}

	var newComment Comment
	resp, err := service.client.Do(req, &newComment)
	if err != nil {
		return nil, resp, err
	}

	return &newComment, resp, err
}

// ListComments returns the list of Comments in a Story.
func (service *StoryService) ListComments(
	projectID int,
	storyID int,
) ([]*Comment, *http.Response, error) {

	u := fmt.Sprintf("projects/%v/stories/%v/comments", projectID, storyID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var comments []*Comment
	resp, err := service.client.Do(req, &comments)
	if err != nil {
		return nil, resp, err
	}

	return comments, resp, nil
}

// ListBlockers returns the list of Blockers in a Story.
func (service *StoryService) ListBlockers(
	projectID int,
	storyID int,
) ([]*Blocker, *http.Response, error) {

	u := fmt.Sprintf("projects/%v/stories/%v/blockers", projectID, storyID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var blockers []*Blocker
	resp, err := service.client.Do(req, &blockers)
	if err != nil {
		return nil, resp, err
	}

	return blockers, resp, nil
}

// AddBlocker will add a Blocker to a Story by ID
func (service *StoryService) AddBlocker(projectID int, storyID int, description string) (*Blocker, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v/blockers", projectID, storyID)
	req, err := service.client.NewRequest("POST", u, BlockerRequest{
		Description: description,
	})
	if err != nil {
		return nil, nil, err
	}

	var blocker Blocker
	resp, err := service.client.Do(req, &blocker)
	if err != nil {
		return nil, resp, err
	}

	return &blocker, resp, nil
}

// UpdateBlocker will change an existing Blocker attached to a story by ID.
func (service *StoryService) UpdateBlocker(projectID, storyID, blockerID int, blocker *BlockerRequest) (*Blocker, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v/blockers/%v", projectID, storyID, blockerID)
	req, err := service.client.NewRequest("PUT", u, blocker)
	if err != nil {
		return nil, nil, err
	}

	var blockerResp Blocker
	resp, err := service.client.Do(req, &blockerResp)
	if err != nil {
		return nil, resp, err
	}

	return &blockerResp, resp, nil
}
