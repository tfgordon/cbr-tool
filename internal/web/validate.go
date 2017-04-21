// Copyright © 2015 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"errors"
	"github.com/carneades/carneades-4/src/engine/caes"
)

// Validate the content of an argument graph, for the purpose
// of the EAGLE argumentation tool
func validate(ag *caes.ArgGraph) error {
	v := []string{}
	if ag.Metadata["title"].(string) == "" {
		v = append(v, "Empty title")
	}
	if ag.Metadata["year"].(string) == "" {
		v = append(v, "Empty year")
	}
	if ag.Metadata["citation"].(string) == "" {
		v = append(v, "Empty citation")
	}
	if ag.Metadata["description"].(string) == "" {
		v = append(v, "Empty description")
	}
	if ag.Metadata["language"].(string) == "" {
		v = append(v, "Empty language")
	}
	if len(v) == 0 {
		return nil
	} else {
		msg := ""
		for _, e := range v {
			msg = msg + "- " + e + "\n"
		}
		return errors.New(msg)
	}
}
