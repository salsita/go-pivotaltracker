// Copyright (c) 2014-2018 Salsita Software
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

// TimeZone is a Pivotal Tracker object for services to return
// as part of the date objects.
type TimeZone struct {
	OlsonName string `json:"olson_name,omitempty"`
	Offset    string `json:"offset,omitempty"`
}
