// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

// Helper functions for retrieving domain models and cases from CouchDB

package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetDomain(couchdbURL string, domainsDBName string, domainId string) (*Domain, error) {
	url := couchdbURL + domainsDBName + "/" + domainId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	domain := Domain{}
	err = json.Unmarshal(body, &domain)

	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func GetCase(couchdbURL string, casesDBName string, caseId string) (*Case, error) {
	url := couchdbURL + casesDBName + "/" + caseId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c := Case{}
	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}
	return &c, nil
}
