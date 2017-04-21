// Copyright © 2015-16 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

// JSON import of domain models

package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func Import(inFile io.Reader) (*Domain, error) {
	data, err := ioutil.ReadAll(inFile)
	if err != nil {
		return nil, err
	}

	m := Domain{}
	err = json.Unmarshal(data, &m)

	if err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}
