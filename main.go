package main

import (
	"./v5/pivotal"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"
)

func doDumpPeople(client *pivotal.Client) error {
	memberships, _, trackerError := client.AccountMemberships.List()
	if trackerError != nil {
		return trackerError
	}

	for _, membership := range memberships {
		fmt.Printf("[%d] %3s %20s %s\n", membership.Person.Id, membership.Person.Initials, membership.Person.Username, membership.Person.Name)
	}

	return nil
}

func main() {
	apiToken := os.Getenv("TRACKER_API_TOKEN")
	if utf8.RuneCountInString(apiToken) == 0 {
		fmt.Println("Please set TRACKER_API_TOKEN")
		return
	}

	accountIdString := os.Getenv("TRACKER_ACCOUNT_ID")
	accountId, err := strconv.Atoi(accountIdString)
	if err != nil {
		fmt.Printf("Could not convert TRACKER_ACCOUNT_ID '%s': %v\n", accountIdString, err)
		return
	}

	client := pivotal.NewClient(apiToken)
	client.SetAccountId(accountId)

	err = doDumpPeople(client)
	if err != nil {
		fmt.Printf("Got Client Error: %v", err)
		return
	}

}
