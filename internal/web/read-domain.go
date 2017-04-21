// Copyright © 2016 Thomas F. Gordon
// This Source Code Form is subject to the terms of the
// Mozilla Public License, v. 2.0. If a copy of the MPL
// was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.

package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
)

func makeReadDomainHandler(cdb CouchDBConfig) func(http.ResponseWriter, *http.Request) {
	couchdbURL := cdb.URL
	domainsDBName := cdb.Domains

	return func(w http.ResponseWriter, req *http.Request) {
		domainId := path.Base(req.URL.Path)
		url := couchdbURL + domainsDBName + "/" + domainId

		// return a status and JSON encoded message to the client,
		// including the id of the newly created domain, if successful
		respond := func(status int, body string) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(body)
		}

		resp, err := http.Get(url)
		if err != nil {
			respond(http.StatusBadRequest, "Could not get the requested domain model from the Couchdb server")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			respond(http.StatusInternalServerError, "Body of response from the Couchdb server could not be read")
			return
		}

		respond(http.StatusOK, string(body))
	}
}
