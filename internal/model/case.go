// Copyright © 2015-16 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package model

type Case struct {
	Id       string `json:"_id"`
	Revision string `json:"_rev"`
	Meta     struct {
		Title       string `json:"title"`
		Citation    string `json:"citation"`
		Court       string
		Year        string // Year Decided
		Description string
		Keywords    string
		Language    string
		Majority    string
		Minority    string
	} `json:"meta"`
	Statements map[string]struct {
		Text  string
		Label string
	} `json:"statements"`
}
