// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"net/http"
)

type AccountMembership struct {
	Kind             string `json:"kind"`
	Id               int    `json:"id"`
	Person           Person `json:"person"`
	IsOwner          bool   `json:"owner"`
	IsAdmin          bool   `json:"admin"`
	IsProjectCreator bool   `json:"project_creator"`
	IsTimeKeepr      bool   `json:"timekeeper"`
	IsTimeEnterer    bool   `json:"time_enterer"`
}

type AccountMembershipService struct {
	client    *Client
	accountId int
}

func newAccountMembershipService(client *Client, accountId int) *AccountMembershipService {
	return &AccountMembershipService{client, accountId}
}

func (service *AccountMembershipService) List() ([]*AccountMembership, *http.Response, error) {
	path := fmt.Sprintf("accounts/%d/memberships", service.accountId)
	req, err := service.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var memberships []*AccountMembership
	resp, err := service.client.Do(req, &memberships)
	if err != nil {
		return nil, resp, err
	}

	return memberships, resp, err
}
