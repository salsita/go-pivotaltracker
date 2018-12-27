// Copyright (c) 2014-2018 Salsita Software
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"time"
)

// Date wraps a compatible time.Time object for marshal/unmarshal operations
type Date time.Time

// UnmarshalJSON implements the json.Unmarshaler() interface for the Date object
func (date *Date) UnmarshalJSON(content []byte) error {
	s := string(content)

	parsingError := func() error {
		return fmt.Errorf(
			"pivotal.Date.UnmarshalJSON: invalid date string: %s", content)
	}

	// Check whether the leading and trailing " is there.
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return parsingError()
	}

	// Strip the leading and trailing "
	s = s[:len(s)-1][1:]

	// Parse the rest.
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return parsingError()
	}

	*date = Date(t)
	return nil
}

// MarshalJSON implements the json.Marshaler() interface for the Date object
func (date Date) MarshalJSON() ([]byte, error) {
	return []byte((time.Time)(date).Format("2006-01-02")), nil
}
