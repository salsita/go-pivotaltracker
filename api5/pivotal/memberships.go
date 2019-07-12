// Copyright (c) 2014-2018 Salsita Software
// Copyright (c) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"net/http"
	"time"
)

// ProjectMembership is the primary data object for the MembershipService.
type ProjectMembership struct {
	Person         Person
	ID             int        `json:"id,omitempty"`
	Kind           string     `json:"kind,omitempty"`
	AccountID      int        `json:"account_id,omitempty"`
	Owner          bool       `json:"owner,omitempty"`
	Admin          bool       `json:"admin,omitempty"`
	ProjectCreator bool       `json:"project_creator,omitempty"`
	Timekeeper     bool       `json:"timekeeper,omitempty"`
	TimeEnterer    bool       `json:"time_enterer,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// MembershipService wraps the client context for interacting with project members.
type MembershipService struct {
	client *Client
}

func newMembershipService(client *Client) *MembershipService {
	return &MembershipService{client}
}

// List all of the memberships in an account.
func (service *MembershipService) List(projectID int) ([]*ProjectMembership, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/memberships", projectID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var projectMemberships []*ProjectMembership
	resp, err := service.client.Do(req, &projectMemberships)
	if err != nil {
		return nil, resp, err
	}

	return projectMemberships, resp, err
}
