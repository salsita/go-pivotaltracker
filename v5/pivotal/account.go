package pivotal

import (
	"fmt"
	"net/http"
	"time"
)

type Account struct {
	CreatedAt *time.Time `json:"created_at"`
	ID        int        `json:"id"`
	Kind      string     `json:"kind"`
	Name      string     `json:"name"`
	Plan      string     `json:"plan"`
	Status    string     `json:"status"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AccountsService struct {
	client *Client
}

func newAccountsService(client *Client) *AccountsService {
	return &AccountsService{client}
}

func (service *AccountsService) Get(accountID int) (*Account, *http.Response, error) {
	u := fmt.Sprintf("accounts/%d", accountID)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var account Account
	resp, err := service.client.Do(req, &account)
	if err != nil {
		return nil, resp, err
	}

	return &account, resp, nil
}
