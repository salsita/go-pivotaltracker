// Copyright (c) 2014 Salsita Software
// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.
package pivotal

import (
	"fmt"
	"net/http"
	"time"
)

type Iteration struct {
	Number          int        `json:"number,omitempty"`
	ProjectId       int        `json:"project_id,omitempty"`
	Length          int        `json:"length,omitempty"`
	TeamStrength    float64    `json:"team_strength,omitempty"`
	StoryIds        []int      `json:"story_ids,omitempty"`
	Stories         []*Story   `json:"stories,omitempty"`
	Start           *time.Time `json:"start,omitempty"`
	Finish          *time.Time `json:"finish,omitempty"`
	Velocity        float64    `json:"velocity,omitempty"`
	Points          int        `json:"points,omitempty"`
	AcceptedPoints  int        `json:"accepted_points,omitempty"`
	EffectivePoints float64    `json:"effective_points,omitempty"`
	Kind            string     `json:"kind,omitempty"`
}

type IterationService struct {
	client *Client
}

func newIterationService(client *Client) *IterationService {
	return &IterationService{client}
}

// Return an iteration from the project.
func (service *IterationService) Get(projectId int, iterationNumber int) (*Iteration, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/iterations/%v", projectId, iterationNumber)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var iteration Iteration
	resp, err := service.client.Do(req, &iteration)
	if err != nil {
		return nil, resp, err
	}

	return &iteration, resp, err
}
